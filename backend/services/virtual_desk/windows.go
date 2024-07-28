//go:build windows

package virtual_desk

import (
	"bufio"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/duke-git/lancet/v2/slice"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"gobot/backend/constants"
	"gobot/backend/dao"
	"gobot/backend/forms"
	"gobot/backend/global"
	"gobot/backend/log"
	"gobot/backend/models"
	"gobot/backend/services/env"
	"gobot/backend/services/sys_exec"
	"gobot/backend/utils"
	"golang.org/x/sys/windows/registry"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

var visualChan chan string
var visualProcess *exec.Cmd
var conn *websocket.Conn

func StartInVirtualDesk(ctx context.Context, id string) error {
	if visualProcess == nil {
		err := enableAutoRun()
		if err != nil {
			return err
		}
		virtualExe := filepath.Join(env.GetExecutablePath(), "VirtualDesk", "VirtualDesk.exe")
		visualProcess = exec.Command(virtualExe, "1024*768")

		stdout, err := visualProcess.StdoutPipe()
		if err != nil {
			return err
		}
		cancelContext, cancel := context.WithCancel(context.Background())
		visualChan = make(chan string, 2)
		go monitorChildSession(ctx, cancelContext, stdout)
		visualChan <- id
		err = visualProcess.Run()
		visualProcess = nil
		cancel()
		close(visualChan)
		if conn != nil {
			_ = conn.Close()
			conn = nil
		}
		if err != nil {
			return err
		}
	} else {
		visualChan <- id
	}

	return nil
}

func TerminateVirtualDeskFlow(ctx context.Context, id string) error {
	connectVirtualServer(ctx)
	req := &forms.Transfer{
		Action: constants.StopAction,
	}

	data := forms.StopFlow{
		InstanceId: id,
	}
	reqData, err := json.Marshal(data)
	if err != nil {
		log.Logger.Logger.Error().Err(err)
		return err
	}
	req.Data = string(reqData)
	reqMessage, err := json.Marshal(req)
	if err != nil {
		log.Logger.Logger.Error().Err(err)
		return err
	}
	if err := conn.WriteMessage(websocket.TextMessage, reqMessage); err != nil {
		_ = conn.Close()
		conn = nil
		return err
	}
	return nil
}

func monitorChildSession(ctx, cancelContext context.Context, outPipe io.ReadCloser) {
	reader := bufio.NewReader(outPipe)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Logger.Logger.Error().Msg(err.Error())
			return
		}
		if line == "LoginComplete\r\n" {
			break
		}
	}
	connectVirtualServer(ctx)
	for {
		select {
		case <-cancelContext.Done():
			log.Logger.Logger.Info().Msg("虚拟桌面关闭了")
			return
		case id := <-visualChan:
			if conn == nil {
				log.Logger.Logger.Error().Msg("连接虚拟桌面失败")
				continue
			}
			log.Logger.Logger.Info().Msg(id)
			project, err := queryProjectById(id)
			if err != nil {
				log.Logger.Logger.Error().Msg(err.Error())
				continue
			}
			instanceId := uuid.New().String()
			global.ProcessQueue = append(global.ProcessQueue, forms.RunningInstance{
				Id:          instanceId,
				ProjectId:   id,
				ProjectName: project.Name,
				StartTs:     time.Now(),
				TriggerType: "手动触发",
				IsVirtual:   true,
				Process:     nil,
			})
			executeParam, err := generateExecuteParam(project, instanceId, "")
			if err != nil {
				log.Logger.Logger.Error().Msg(err.Error())
				continue
			}
			req := &forms.Transfer{
				Action: constants.StartAction,
			}
			module := "robot_core.robot_interpreter"
			if project.IsFlow {
				module = "robot_core.flow_executor"
			}
			data := forms.StartFlow{
				InstanceId: instanceId,
				ExePath:    utils.GetVenvPython(project.Path),
				Module:     module,
				Param:      executeParam,
			}
			reqData, err := json.Marshal(data)
			if err != nil {
				log.Logger.Logger.Error().Err(err)
				continue
			}
			req.Data = string(reqData)
			reqMessage, err := json.Marshal(req)
			if err != nil {
				log.Logger.Logger.Error().Err(err)
				continue
			}
			if err := conn.WriteMessage(websocket.TextMessage, reqMessage); err != nil {
				log.Logger.Logger.Error().Msg(err.Error())
				_ = conn.Close()
				conn = nil
			}
		}
	}
}

func queryProjectById(id string) (result *models.Project, err error) {
	if err := dao.ReadTx(dao.MainDB, func(tx dao.Tx) error {
		result, err = tx.SelectProject(id)
		return err
	}); err != nil {
		return nil, err
	}
	return result, nil
}

func enableAutoRun() error {
	key, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Run`, registry.SET_VALUE)
	if err != nil {
		log.Logger.Logger.Error().AnErr("无法打开注册表键", err)
		return err
	}
	defer func(key registry.Key) {
		err := key.Close()
		if err != nil {
			log.Logger.Logger.Error().AnErr("关闭注册表键失败", err)
		}
	}(key)
	err = key.SetStringValue("VirtualHost", filepath.Join(env.GetExecutablePath(), "VirtualHost.exe"))
	if err != nil {
		log.Logger.Logger.Error().AnErr("设置注册表值失败", err)
		return err
	}

	log.Logger.Logger.Info().Msg("注册表设置成功，程序已设置为开机自启动。")
	return nil
}

func disableAutoRun() error {
	key, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Run`, registry.SET_VALUE)
	if err != nil {
		log.Logger.Logger.Error().AnErr("无法打开注册表键", err)
		return err
	}
	defer func(key registry.Key) {
		err := key.Close()
		if err != nil {
			log.Logger.Logger.Error().AnErr("关闭注册表键失败", err)
		}
	}(key)

	err = key.DeleteValue("VirtualHost")
	if err != nil {
		log.Logger.Logger.Error().AnErr("删除注册表值失败", err)
		return err
	}

	log.Logger.Logger.Info().Msg("注册表删除成功，程序已取消开机自启动。")
	return nil
}

func connectVirtualServer(ctx context.Context) {
	for visualProcess != nil && conn == nil {
		data, err := os.ReadFile(filepath.Join(env.GetExecutablePath(), "virtual_host_port"))
		if err != nil {
			continue
		}
		log.Logger.Logger.Info().Msgf("nativehost_port: %s", string(data))
		conn, _, err = websocket.DefaultDialer.Dial(fmt.Sprintf("ws://127.0.0.1:%s/ws", string(data)), nil)
		if err != nil {
			log.Logger.Logger.Error().Err(err)
			conn = nil
			time.Sleep(time.Second)
			continue
		}
		log.Logger.Logger.Info().Msg("连接虚拟桌面成功")
		_ = disableAutoRun()
		go listenResponse(ctx, conn)
		break
	}
}

func listenResponse(ctx context.Context, conn *websocket.Conn) {
	for conn != nil {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Logger.Logger.Error().Err(err)
			_ = conn.Close()
			conn = nil
			return
		}
		log.Logger.Logger.Info().Msgf("收到虚拟桌面消息：%s", message)
		resp := forms.Transfer{}
		if err = json.Unmarshal(message, &resp); err != nil {
			log.Logger.Logger.Error().Err(err)
			continue
		}
		if resp.Action == constants.ExecuteResultAction || resp.Action == constants.StopResultAction {
			data := resp.Data
			result := forms.Result{}
			if err = json.Unmarshal([]byte(data), &result); err != nil {
				log.Logger.Logger.Error().Err(err)
				continue
			}
			if resp.Action == constants.StopResultAction && !result.Success {
				log.Logger.Logger.Error().Msgf("停止进程失败%s", result.InstanceId)
				continue
			}
			proc, ok := slice.FindBy(global.ProcessQueue, func(index int, item forms.RunningInstance) bool {
				return item.Id == result.InstanceId
			})
			if !ok {
				log.Logger.Logger.Warn().Msg("找不到对应的进程")
				continue
			}
			global.ProcessQueue = slice.Filter(global.ProcessQueue, func(index int, item forms.RunningInstance) bool {
				return item.Id != result.InstanceId
			})
			execution := models.Execution{Id: proc.Id, ProjectId: proc.ProjectId, TriggerType: proc.TriggerType, StartTs: proc.StartTs, EndTs: time.Now()}
			if resp.Action == constants.StopResultAction {
				execution.ExecuteResult = 3
			} else if result.Success {
				execution.ExecuteResult = 1
			} else {
				execution.ExecuteResult = 0
				execution.ErrorMsg = result.Message
			}
			runtime.EventsEmit(ctx, "execute_event")
			_ = dao.WriteTx(dao.MainDB, func(tx dao.Tx) error {
				if err := tx.InsertExecution(&execution); err != nil {
					return err
				}
				return nil
			})
			project, err := queryProjectById(execution.ProjectId)
			if err != nil {
				log.Logger.Logger.Error().Err(err)
				continue
			}
			logPath := filepath.Join(project.Path, constants.LogDir, result.InstanceId+".log")
			logs, err := os.ReadFile(logPath)
			if err == nil {
				_ = dao.WriteTx(dao.LogDB, func(tx dao.Tx) error {
					if err := tx.InsertExecutionLog(string(logs), execution.Id); err != nil {
						return err
					}
					return nil
				})
			}
			_ = os.Remove(logPath)
		}
	}

}

func generateExecuteParam(project *models.Project, instanceId, subId string) (string, error) {
	projectPath := filepath.Join(project.Path, constants.BaseDir)
	logPath := filepath.Join(project.Path, constants.LogDir, instanceId+".log")
	params := make(map[string]interface{})
	logDir := filepath.Dir(logPath)
	if !utils.PathExist(logDir) {
		if err := os.MkdirAll(logDir, fs.ModeDir); err != nil {
			return "", err
		}
	}
	params["sys_path_list"] = []string{projectPath, filepath.Dir(projectPath)}
	params["log_path"] = logPath
	params["log_level"] = "INFO"

	env := make(map[string]string)
	env["project_path"] = projectPath
	if subId != "" {
		params["mod"] = subId
	} else {
		if project.IsFlow {
			flowPath := filepath.Join(projectPath, constants.MainFlow)
			params["flow_path"] = flowPath
		} else {
			params["mod"] = "main"
		}
	}

	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exePath := filepath.Dir(ex)
	env["execute_path"] = exePath
	params["environment_variables"] = env
	params["inputs"] = project.InputParam
	marshalParam, err := json.Marshal(params)
	if err != nil {
		return "", err
	}
	base64Param := base64.StdEncoding.EncodeToString(marshalParam)
	log.Logger.Logger.Info().Msg(base64Param)
	return base64Param, nil
}
func StopVisualHost() {
	if visualProcess != nil {
		_ = sys_exec.KillProcess(visualProcess)
		visualProcess = nil
	}
}
