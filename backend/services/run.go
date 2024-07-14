package services

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/duke-git/lancet/v2/slice"
	"github.com/duke-git/lancet/v2/strutil"
	"gobot/backend/constants"
	"gobot/backend/dao"
	"gobot/backend/forms"
	"gobot/backend/global"
	"gobot/backend/log"
	"gobot/backend/models"
	"gobot/backend/services/sys_exec"
	"gobot/backend/services/virtual_desk"
	"gobot/backend/utils"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hpcloud/tail"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

var monitorMap = make(map[string]context.CancelFunc)

var runningProcess *exec.Cmd

func RunSubFlow(ctx context.Context, id, subId string) error {
	project, err := QueryProjectById(id)
	if err != nil {
		return nil
	}
	projectPath := project.Path + string(os.PathSeparator) + constants.BaseDir
	instanceId := uuid.NewString()
	logPath := filepath.Join(project.Path, constants.LogDir, instanceId+".log")
	runningProcess, err = generateRunCmd(project, instanceId, subId)
	if err != nil {
		return err
	}
	var stderr bytes.Buffer
	runningProcess.Stderr = &stderr
	background, cancel := context.WithCancel(context.Background())

	go func() {
		monitorLog(id, ctx, background, logPath)
	}()
	_, err = runningProcess.Output()
	runningProcess = nil
	time.AfterFunc(2*time.Second, func() {
		cancel()
	})
	errStr := ""
	if err != nil {
		log.Logger.Logger.Error().Err(err).Msg("流程块运行失败")
		errStr = stderr.String()
		errStr = strutil.After(errStr, projectPath)
	}
	_ = os.Remove(logPath)
	runtime.EventsEmit(ctx, "execute_event")
	if errStr != "" {
		return errors.New(errStr)
	}
	return nil
}

func TerminateSubFlow() error {
	if runningProcess != nil {
		err := sys_exec.KillProcess(runningProcess)
		if err != nil {
			return err
		}
		runningProcess = nil
	}
	if debugProcess != nil {
		err := sys_exec.KillProcess(debugProcess)
		if err != nil {
			return err
		}
		debugProcess = nil
	}
	if replProcess != nil {
		err := sys_exec.KillProcess(replProcess)
		if err != nil {
			return err
		}
		replProcess = nil
	}
	return nil
}

func RunProject(ctx context.Context, id, triggerType string) error {
	project, err := QueryProjectById(id)
	if err != nil {
		return nil
	}
	runningInstanceId := uuid.New().String()
	logPath := filepath.Join(project.Path, constants.LogDir, runningInstanceId+".log")
	command, err := generateRunCmd(project, runningInstanceId, "")
	if err != nil {
		return err
	}
	global.ProcessQueue = append(global.ProcessQueue, forms.RunningInstance{
		Id:          runningInstanceId,
		ProjectId:   id,
		ProjectName: project.Name,
		StartTs:     time.Now(),
		TriggerType: triggerType,
		IsVirtual:   false,
		Process:     command,
	})
	var stderr bytes.Buffer
	command.Stderr = &stderr
	execution := models.Execution{Id: runningInstanceId, ProjectId: id, TriggerType: triggerType, StartTs: time.Now()}

	go func() {
		_, err = command.Output()
		global.ProcessQueue = slice.Filter(global.ProcessQueue, func(index int, item forms.RunningInstance) bool {
			return item.Id != runningInstanceId
		})
		runtime.EventsEmit(ctx, "execute_event")
		errStr := ""
		execution.EndTs = time.Now()
		if err != nil {
			log.Logger.Logger.Error().Err(err).Msg("项目运行失败")
			errStr = stderr.String()
			log.Logger.Logger.Error().Msg(errStr)
			if errStr == "" {
				execution.ExecuteResult = 3
			} else {
				execution.ExecuteResult = 0
			}
			execution.ErrorMsg = errStr
		} else {
			execution.ExecuteResult = 1
		}
		_ = dao.WriteTx(dao.MainDB, func(tx dao.Tx) error {
			if err := tx.InsertExecution(&execution); err != nil {
				return err
			}
			return nil
		})
		logs, err := os.ReadFile(logPath)
		logs = append(logs, []byte(errStr)...)
		if err == nil {
			_ = dao.WriteTx(dao.LogDB, func(tx dao.Tx) error {
				if err := tx.InsertExecutionLog(string(logs), execution.Id); err != nil {
					return err
				}
				return nil
			})
		}
		_ = os.Remove(logPath)
	}()

	return nil
}

func TerminateFlow(ctx context.Context, id string) error {
	result, ok := slice.FindBy(global.ProcessQueue, func(index int, item forms.RunningInstance) bool {
		return item.Id == id
	})
	if ok {
		if !result.IsVirtual {
			err := sys_exec.KillProcess(result.Process)
			if err != nil {
				return err
			}
			global.ProcessQueue = slice.Filter(global.ProcessQueue, func(index int, item forms.RunningInstance) bool {
				return item.Id != id
			})
		} else {
			err := virtual_desk.TerminateVirtualDeskFlow(ctx, id)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func GetRunningFlows() (resultList []forms.RunningInstance, err error) {
	return global.ProcessQueue, nil
}

func StartMonitorLog(id string, ctx context.Context) {
	result, ok := slice.FindBy(global.ProcessQueue, func(index int, item forms.RunningInstance) bool {
		return item.Id == id
	})
	if !ok {
		return
	}
	project, err := QueryProjectById(result.ProjectId)
	if err != nil {
		return
	}
	logPath := project.Path + string(os.PathSeparator) + constants.LogDir + string(os.PathSeparator) + id + ".log"
	background, cancel := context.WithCancel(context.Background())
	monitorMap[id] = cancel
	go monitorLog(id, ctx, background, logPath)
}

func StopMonitorLog(id string) {
	if cancel, ok := monitorMap[id]; ok {
		cancel()
		delete(monitorMap, id)
	}
}

func monitorLog(id string, ctx, background context.Context, logPath string) {
	config := tail.Config{
		ReOpen:    true,                                 //重新打开
		Follow:    true,                                 //是否跟随
		Location:  &tail.SeekInfo{Offset: 0, Whence: 0}, //从哪儿开始读
		MustExist: false,                                //文件不存在不报错
		Poll:      true,                                 //
		Logger:    tail.DiscardingLogger,
	}
	tails, err := tail.TailFile(logPath, config)
	if err != nil {
		log.Logger.Logger.Error().Err(err)
		return
	}

	var (
		line *tail.Line
		ok   bool
	)
loop:
	for {
		select {
		case <-background.Done():
			log.Logger.Logger.Info().Msg("结束监控")
			break loop
		case line, ok = <-tails.Lines:
			if !ok {
				time.Sleep(time.Second)
				continue
			}
			if !strings.Contains(line.Text, "DEBUG") {
				runtime.EventsEmit(ctx, "log_event", map[string]string{"id": id, "text": line.Text})
			}
		}
	}
	_ = tails.Stop()
}

func generateRunCmd(project *models.Project, instanceId, subId string) (*exec.Cmd, error) {
	projectPath := filepath.Join(project.Path, constants.BaseDir)
	logPath := filepath.Join(project.Path, constants.LogDir, instanceId+".log")
	module := "robot_core.robot_interpreter"
	params := make(map[string]interface{})
	logDir := filepath.Dir(logPath)
	if !utils.PathExist(logDir) {
		if err := os.MkdirAll(logDir, fs.ModeDir); err != nil {
			return nil, err
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
			module = "robot_core.flow_executor"
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
		return nil, err
	}
	base64Param := base64.StdEncoding.EncodeToString(marshalParam)
	log.Logger.Logger.Info().Msg(base64Param)
	return sys_exec.BuildCmd(utils.GetVenvPython(project.Path), "-u", "-m", module, base64Param), nil
}
