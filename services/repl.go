package services

import (
	"bufio"
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"gobot/constants"
	"gobot/log"
	"gobot/services/sys_exec"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/spf13/viper"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

var command *exec.Cmd
var replChan chan string
var globalImport map[string]bool
var appCtx context.Context

func startReplCommand(id string) error {
	params := make(map[string]interface{})
	project, err := QueryProjectById(id)
	if err != nil {
		return err
	}
	projectPath := project.Path + string(os.PathSeparator) + constants.BaseDir

	params["sys_path_list"] = []string{projectPath}
	params["log_level"] = "Error"
	env := make(map[string]string)
	env["project_path"] = projectPath
	params["environment_variables"] = env
	marshalParam, err := json.Marshal(params)
	if err != nil {
		return err
	}
	base64Param := base64.StdEncoding.EncodeToString(marshalParam)
	log.Logger.Info(base64Param)
	command = sys_exec.BuildCmd(viper.GetString("python.path"), "-u", "-B", "-m", "robot_core.robot_debugger", base64Param)

	var stderr bytes.Buffer
	command.Stderr = &stderr
	inPipe, err := command.StdinPipe()
	if err != nil {
		return err
	}
	outPipe, err := command.StdoutPipe()
	if err != nil {
		return err
	}
	go dealRepl(inPipe, outPipe)
	err = command.Run()
	if err != nil {
		errStr := stderr.String()
		_, errStr, founded := strings.Cut(errStr, projectPath)
		log.Logger.Error(fmt.Sprintf("%t", founded))
		log.Logger.Error(errStr)
		return errors.New(errStr)
	}
	return nil
}

func dealRepl(inPipe io.WriteCloser, outPipe io.ReadCloser) error {
	writer := bufio.NewWriter(inPipe)
	out, err := waitToWrite(outPipe)
	log.Logger.Info(out)
	if err != nil {
		return err
	}
	re, _ := regexp.Compile("code_map_id=\"([^\"]*)")
	for {
		code := <-replChan
		log.Logger.Info(strings.Trim(code, "\n"))
		_, err = writer.WriteString(base64.StdEncoding.EncodeToString([]byte(strings.Trim(code, "\n"))) + "\n")
		writer.Flush()
		if err != nil {
			return err
		}
		out, err = waitToWrite(outPipe)
		match := re.FindAllString(code, -1)
		out = strings.Trim(out, ">")
		if len(match) > 0 {
			runtime.EventsEmit(appCtx, "repl", strings.TrimPrefix(match[0], "code_map_id=\""), out)
		}
		if len(out) > 0 {
			log.Logger.Info(strings.Trim(out, "\n"))
		}
		if err != nil {
			return err
		}
	}
}

func waitToWrite(r io.Reader) (result string, err error) {
	reader := bufio.NewReader(r)
	for {
		out, err := reader.ReadString('>')
		if err != nil {
			return result, err
		}
		result += out
		if strings.HasSuffix(result, ">>>") {
			return result, err
		}
	}
}

func RunActivity(ctx context.Context, id, req string) error {
	appCtx = ctx
	if command == nil {
		globalImport = map[string]bool{}
		replChan = make(chan string)
		go startReplCommand(id)
	}
	activity := Activity{}
	err := json.Unmarshal([]byte(req), &activity)
	if err != nil {
		return err
	}
	namespace := make(map[string]string)
	code := activity.GeneratePythonCode(namespace, 0)
	for key := range namespace {
		if _, ok := globalImport[key]; !ok {
			replChan <- "import " + key + "\n"
		}
		globalImport[key] = true
	}
	replChan <- code
	return nil
}

func RestartReplCommand(id string) error {
	if command != nil {
		if err := command.Process.Kill(); err != nil {
			return nil
		}
		for k := range globalImport {
			delete(globalImport, k)
		}
	}
	globalImport = map[string]bool{}
	replChan = make(chan string)
	go startReplCommand(id)
	return nil
}
