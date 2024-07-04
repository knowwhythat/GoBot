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
	"gobot/backend/log"
	"gobot/backend/models"
	"gobot/backend/services/sys_exec"
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

var processQueue []forms.RunningInstance
var monitorMap = make(map[string]context.CancelFunc)

var runningProcess *exec.Cmd

func RunSubFlow(ctx context.Context, id, subId string) error {
	project, err := QueryProjectById(id)
	if err != nil {
		return nil
	}
	projectPath := project.Path + string(os.PathSeparator) + constants.BaseDir
	logPath := project.Path + string(os.PathSeparator) + constants.LogDir + string(os.PathSeparator) + uuid.NewString() + ".log"
	runningProcess, err = generateRunCmd(project, projectPath, logPath, subId)
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
		errStr = stderr.String()
		errStr = strutil.After(errStr, projectPath)
	}
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

func RunSequence(ctx context.Context, id, subId, triggerType string) error {
	project, err := QueryProjectById(id)
	if err != nil {
		return nil
	}
	runningInstanceId := uuid.New().String()
	projectPath := project.Path + string(os.PathSeparator) + constants.BaseDir
	logPath := project.Path + string(os.PathSeparator) + constants.LogDir + string(os.PathSeparator) + runningInstanceId + ".log"
	command, err := generateRunCmd(project, projectPath, logPath, subId)
	if err != nil {
		return err
	}
	processQueue = append(processQueue, forms.RunningInstance{
		Id:          runningInstanceId,
		ProjectId:   id,
		ProjectName: project.Name,
		StartTs:     time.Now(),
		TriggerType: triggerType,
		Process:     command,
	})
	var stderr bytes.Buffer
	command.Stderr = &stderr
	execution := models.Execution{Id: runningInstanceId, ProjectId: id, SubFlowId: subId, TriggerType: triggerType, StartTs: time.Now()}

	go func() {
		_, err = command.Output()
		processQueue = removeValue(processQueue, runningInstanceId)
		runtime.EventsEmit(ctx, "execute_event")
		errStr := ""
		execution.EndTs = time.Now()
		if err != nil {
			errStr = stderr.String()
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

func TerminateFlow(id string) error {
	for _, item := range processQueue {
		if item.Id == id {
			err := sys_exec.KillProcess(item.Process)
			if err != nil {
				return err
			}
			processQueue = removeValue(processQueue, id)
		}
	}
	return nil
}

func GetRunningFlows() (resultList []forms.RunningInstance, err error) {
	return processQueue, nil
}

func StartMonitorLog(id string, ctx context.Context) {
	result, ok := slice.FindBy(processQueue, func(index int, item forms.RunningInstance) bool {
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
		log.Logger.Error("tail file failed, err:" + err.Error())
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
			log.Logger.Info("结束监控")
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
	//_ = os.Remove(logPath)
}

func removeValue(slice []forms.RunningInstance, id string) []forms.RunningInstance {
	var temp []forms.RunningInstance
	for _, item := range slice {
		if item.Id != id {
			temp = append(temp, item)
		}
	}
	return temp
}

func generateRunCmd(project *models.Project, projectPath, logPath, subId string) (*exec.Cmd, error) {
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
	params["mod"] = subId
	env := make(map[string]string)
	env["project_path"] = projectPath
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
	log.Logger.Info(base64Param)
	return sys_exec.BuildCmd(utils.GetVenvPython(project.Path), "-u", "-m", "robot_core.robot_interpreter", base64Param), nil
}
