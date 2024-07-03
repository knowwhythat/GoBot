package services

import (
	"bufio"
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"gobot/backend/constants"
	"gobot/backend/log"
	"gobot/backend/services/sys_exec"
	"gobot/backend/utils"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

var debugChan chan string
var debugProcess *exec.Cmd

func DebugSubFlow(ctx context.Context, id string, subId string) error {
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
	params["debug"] = true
	params["mod"] = subId
	env := make(map[string]string)
	env["project_path"] = projectPath
	params["environment_variables"] = env
	params["inputs"] = project.InputParam
	marshalParam, err := json.Marshal(params)
	if err != nil {
		return err
	}
	base64Param := base64.StdEncoding.EncodeToString(marshalParam)
	log.Logger.Info(base64Param)
	debugProcess = sys_exec.BuildCmd(utils.GetVenvPython(project.Path), "-u", "-m", "robot_core.robot_interpreter", base64Param)
	var stderr bytes.Buffer
	debugProcess.Stderr = &stderr
	inPipe, err := debugProcess.StdinPipe()
	if err != nil {
		return err
	}
	outPipe, err := debugProcess.StdoutPipe()
	if err != nil {
		return err
	}
	breakpoints := getBreakpoints(projectPath, subId)
	go dealDebug(ctx, subId, inPipe, outPipe, breakpoints)
	background, cancel := context.WithCancel(context.Background())
	go monitorLog(id, ctx, background, logPath)
	debugChan = make(chan string)
	err = debugProcess.Run()
	debugProcess = nil
	cancel()
	close(debugChan)
	debugChan = nil
	if err != nil {
		errStr := stderr.String()
		_, errStr, founded := strings.Cut(errStr, projectPath)
		log.Logger.Error(fmt.Sprintf("%t", founded))
		log.Logger.Error(errStr)
		return errors.New(errStr)
	}
	return nil
}

func DealDebugSignal(command string) error {
	debugChan <- command
	return nil
}

func getBreakpoints(projectPath, fileName string) []string {
	content, err := os.ReadFile(projectPath + string(os.PathSeparator) + ".dev" + string(os.PathSeparator) + fileName + ".flow")
	if err != nil {
		return []string{}
	}
	flow := Flow{}
	err = json.Unmarshal(content, &flow)
	if err != nil {
		return []string{}
	}
	ids := getActivityBreakpoints(flow.Sequence)
	content, err = os.ReadFile(projectPath + string(os.PathSeparator) + fileName + ".py")
	if err != nil {
		return []string{}
	}
	lines := strings.Split(string(content), "\n")
	breakpoints := make([]string, 0)
	for i := 1; i <= len(lines); i++ {
		for _, id := range ids {
			if strings.Contains(lines[i-1], id) {
				breakpoints = append(breakpoints, strconv.Itoa(i))
			}
		}
	}
	return breakpoints
}

func getActivityBreakpoints(activity Activity) []string {
	ids := make([]string, 0)
	for _, a := range activity.Children {
		if a.Breakpoint {
			ids = append(ids, a.ID)
		}
		ids = append(ids, getActivityBreakpoints(a)...)
	}
	return ids
}

func dealDebug(ctx context.Context, filename string, inPipe io.WriteCloser, outPipe io.ReadCloser, breakpoints []string) {
	defer func(inPipe io.WriteCloser, outPipe io.ReadCloser) {
		_ = inPipe.Close()
		_ = outPipe.Close()
	}(inPipe, outPipe)
	reader := bufio.NewReader(outPipe)
	writer := bufio.NewWriter(inPipe)
	output := ""
	breakpoint := 0
	first := true
	re, _ := regexp.Compile("code_map_id=\"([^\"]*)")
	for {
		line, err := reader.ReadString(')')
		if err != nil || io.EOF == err {
			break
		}
		output += line
		if strings.HasSuffix(line, "\n(Pdb)") {
			log.Logger.Info(output)
			match := re.FindAllString(output, -1)
			if len(match) > 0 {
				runtime.EventsEmit(ctx, "debug", strings.TrimPrefix(match[0], "code_map_id=\""))
			}
			if breakpoint < len(breakpoints) {
				_, err = writer.WriteString("b " + filename + ".py:" + breakpoints[breakpoint] + "\n")
				breakpoint += 1
			} else if strings.Contains(output, "--Return--") {
				_, err = writer.WriteString("c\n")
			} else {
				if first {
					_, err = writer.WriteString("c\n")
					first = false
				} else {
					if len(match) == 0 {
						_, err = writer.WriteString("n\n")
					} else {
						debugCommand := <-debugChan
						_, err = writer.WriteString(debugCommand + "\n")
					}
				}
			}
			if err != nil {
				if debugProcess != nil {
					_ = sys_exec.KillProcess(debugProcess)
				}
				break
			}
			err = writer.Flush()
			if err != nil {
				if debugProcess != nil {
					_ = sys_exec.KillProcess(debugProcess)
				}
				break
			}
			output = ""
		}
	}
}
