//go:build windows

package virsual_desk

import (
	"bufio"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"gobot/backend/constants"
	"gobot/backend/dao"
	"gobot/backend/forms"
	"gobot/backend/global"
	"gobot/backend/log"
	"gobot/backend/models"
	"gobot/backend/utils"
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

func StartInVirtualDesk(id string) error {
	if visualProcess == nil {
		executable, err := os.Executable()
		if err != nil {
			return err
		}
		virtualExe := filepath.Join(filepath.Dir(executable), "VirtualDesk", "VirtualDesk.exe")
		visualProcess = exec.Command(virtualExe, "1024*768")
		stdout, err := visualProcess.StdoutPipe()
		if err != nil {
			return err
		}
		ctx, cancel := context.WithCancel(context.Background())
		visualChan = make(chan string, 2)
		go monitorChildSession(ctx, stdout)
		visualChan <- id
		err = visualProcess.Run()
		visualProcess = nil
		cancel()
		close(visualChan)
		_ = conn.Close()
		conn = nil
		if err != nil {
			return err
		}
	} else {
		visualChan <- id
	}

	return nil
}

func TerminateVirtualDeskFlow(id string) error {
	if err := conn.WriteMessage(1, []byte("")); err != nil {
		return err
	}
	return nil
}

func monitorChildSession(ctx context.Context, outPipe io.ReadCloser) {
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
	connectVirtualServer()
	for {
		select {
		case <-ctx.Done():
			log.Logger.Logger.Info().Msg("虚拟桌面关闭了")
			return
		case id := <-visualChan:
			log.Logger.Logger.Info().Msg(id)
			project, err := queryProjectById(id)
			if err != nil {
				return
			}
			runningInstanceId := uuid.New().String()
			global.ProcessQueue = append(global.ProcessQueue, forms.RunningInstance{
				Id:          runningInstanceId,
				ProjectId:   id,
				ProjectName: project.Name,
				StartTs:     time.Now(),
				TriggerType: "手动触发",
				IsVirtual:   true,
				Process:     nil,
			})
			execution := models.Execution{Id: runningInstanceId, ProjectId: id, TriggerType: "手动触发", StartTs: time.Now()}
			projectPath := filepath.Join(project.Path, string(os.PathSeparator), constants.BaseDir)
			logPath := filepath.Join(project.Path, string(os.PathSeparator), constants.LogDir, string(os.PathSeparator), runningInstanceId, ".log")
			executeParam, err := generateExecuteParam(project, projectPath, logPath, "main")
			if err := conn.WriteMessage(1, []byte(executeParam)); err != nil {
				log.Logger.Logger.Error().Msg(err.Error())
			}
			log.Logger.Logger.Info().Msg(execution.Id)
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

func connectVirtualServer() {
	for visualProcess != nil {
		conn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://127.0.0.1:%s", "9999"), nil)
		if err != nil {
			time.Sleep(time.Second)
			continue
		}
		go listenResponse(conn)
		break
	}
}

func listenResponse(conn *websocket.Conn) {
	for conn != nil {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Logger.Logger.Error().Err(err)
			return
		}
		resp := make(map[string]interface{})
		if err = json.Unmarshal(message, &resp); err != nil {
			log.Logger.Logger.Error().Err(err)
			return
		}
	}

}

func generateExecuteParam(project *models.Project, projectPath, logPath, subId string) (string, error) {
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
	params["mod"] = subId
	env := make(map[string]string)
	env["project_path"] = projectPath
	if project.IsFlow {
		flowPath := filepath.Join(projectPath, constants.MainFlow)
		params["flow_path"] = flowPath
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
	return base64Param, nil
}
