package services

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"gobot/backend/constants"
	"gobot/backend/dao"
	"gobot/backend/log"
	"gobot/backend/models"
	"gobot/backend/services/sys_exec"
	"gobot/backend/utils"
	"io/fs"
	"os"
	"os/exec"
	"strings"
	"time"

	uuid "github.com/google/uuid"
	"github.com/hpcloud/tail"
	"github.com/spf13/viper"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

var processMap = make(map[string]*exec.Cmd)

func RunSubFlow(ctx context.Context, id string, subId string) error {
	params := make(map[string]interface{})
	project, err := QueryProjectById(id)
	if err != nil {
		return err
	}
	projectPath := project.Path + string(os.PathSeparator) + constants.BaseDir
	logDir := project.Path + string(os.PathSeparator) + constants.LogDir
	logPath := logDir + string(os.PathSeparator) + uuid.New().String() + ".log"
	if !utils.PathExist(logDir) {
		if err = os.MkdirAll(logDir, fs.ModeDir); err != nil {
			return err
		}
	}
	params["sys_path_list"] = []string{projectPath}
	params["log_path"] = logPath
	params["log_level"] = "DEBUG"
	params["mod"] = subId
	env := make(map[string]string)
	env["project_path"] = projectPath
	params["environment_variables"] = env
	marshalParam, err := json.Marshal(params)
	if err != nil {
		return err
	}
	base64Param := base64.StdEncoding.EncodeToString(marshalParam)
	log.Logger.Info(base64Param)
	command := sys_exec.BuildCmd(viper.GetString("python.path"), "-u", "-m", "robot_core.robot_interpreter", base64Param)
	processMap[subId] = command
	var stderr bytes.Buffer
	command.Stderr = &stderr
	execution := models.Execution{Id: uuid.New().String(), ProjectId: id, SubFlowId: subId, StartTs: time.Now()}
	background, cancel := context.WithCancel(context.Background())
	background = context.WithValue(background, constants.ExecutionId("executionId"), execution.Id)

	go func() {
		monitorLog(ctx, background, logPath)
	}()
	_, err = command.Output()
	delete(processMap, subId)
	time.AfterFunc(2*time.Second, func() {
		cancel()
	})
	execution.EndTs = time.Now()
	if err != nil {
		errStr := stderr.String()
		_, errStr, founded := strings.Cut(errStr, projectPath)
		log.Logger.Error(fmt.Sprintf("%t", founded))
		log.Logger.Error(errStr)
		if errStr == "" {
			execution.ExecuteResult = 3
		} else {
			execution.ExecuteResult = 0
		}
		return errors.New(errStr)
	} else {
		execution.ExecuteResult = 1
	}
	if err := dao.WriteTx(dao.LogDB, func(tx dao.Tx) error {
		if err := tx.InsertExecution(&execution); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func TerminateSubFlow(id string, subId string) error {
	if command, ok := processMap[subId]; ok {
		err := sys_exec.KillProcess(command)
		return err
	}
	return nil
}

func monitorLog(ctx, background context.Context, logPath string) {
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
				runtime.EventsEmit(ctx, "log_event", line.Text)
			}
		}
	}
	_ = tails.Stop()
	if executeId, ok := background.Value(constants.ExecutionId("executionId")).(string); ok {
		logs, err := os.ReadFile(logPath)
		if err != nil {
			_ = dao.WriteTx(dao.LogDB, func(tx dao.Tx) error {
				if err := tx.InsertExecutionLog(string(logs), uuid.MustParse(executeId)); err != nil {
					return err
				}
				return nil
			})
		}
	}
	_ = os.Remove(logPath)
}
