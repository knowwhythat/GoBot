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
	"gobot/dao"
	"gobot/forms"
	"gobot/log"
	"gobot/models"
	"gobot/services/sys_exec"
	"gobot/utils"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/ahmetb/go-linq"
	uuid "github.com/google/uuid"
	"github.com/hpcloud/tail"
	"github.com/spf13/viper"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

var processMap = make(map[string]*exec.Cmd)
var debugChan chan string

func QueryProjectPage(page forms.QueryForm) (total int, resultList []*models.Project, err error) {
	if err := dao.ReadTx(dao.MainDB, func(tx dao.Tx) error {
		allProjects, err := tx.SelectAllProject()
		if err != nil {
			return err
		}

		query := linq.From(allProjects).Where(func(project interface{}) bool {
			projectVlue := project.(*models.Project)
			return strings.Contains(projectVlue.Name, page.Query)
		}).OrderByDescendingT(func(project *models.Project) int64 {
			return project.CreateTs.UnixNano()
		}).Query
		total = query.Count()

		if page.PageNum > 0 && page.PageSize > 0 {
			start, end := page.Range()
			query = query.Skip(start).Take(end - start)
		}
		projects := query.Results()

		for _, project := range projects {
			project := project.(*models.Project)
			resultList = append(resultList, project)
		}
		return nil
	}); err != nil {
		return 0, nil, err
	}

	return total, resultList, nil
}

func CountProject() (total int, err error) {
	if err := dao.ReadTx(dao.MainDB, func(tx dao.Tx) error {
		allProjects, err := tx.SelectAllProject()
		if err != nil {
			return err
		}
		total = len(allProjects)
		return nil
	}); err != nil {
		return 0, err
	}
	return total, nil
}

func QueryProjectById(id string) (result *models.Project, err error) {
	if err := dao.ReadTx(dao.MainDB, func(tx dao.Tx) error {
		result, err = tx.SelectProject(id)
		return err
	}); err != nil {
		return nil, err
	}
	return result, nil
}

func AddProject(name string) (result *models.Project, err error) {
	path := viper.GetString("python.path")
	dataPath := viper.GetString("data.path")
	id := uuid.New()
	projectDir := dataPath + constants.BaseDir + string(os.PathSeparator) + id.String()
	if !utils.PathExist(projectDir) {
		os.MkdirAll(projectDir, fs.ModeDir)
	}
	result = &models.Project{
		Id:          id,
		Name:        name,
		Description: "",
		Path:        projectDir,
	}
	command := sys_exec.BuildCmd(path, "-m", "venv", result.Path+"\\venv")
	output, err := command.Output()
	if err != nil {
		fmt.Println(string(output))
		return nil, err
	}
	if err := dao.WriteTx(dao.MainDB, func(tx dao.Tx) error {
		if err := tx.InsertProject(result); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return result, nil
}

func ModifyProject(id string, form models.Project) (result *models.Project, err error) {
	if err := dao.WriteTx(dao.MainDB, func(tx dao.Tx) error {
		project, err := tx.SelectProject(id)
		if err != nil {
			return err
		}
		if project == nil {
			return errors.New("project not found")
		}

		if err := tx.UpdateProject(project); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return result, nil
}

func RemoveProject(id string) (err error) {
	if err := dao.WriteTx(dao.MainDB, func(tx dao.Tx) error {
		if project, err := tx.SelectProject(id); err != nil {
			return err
		} else if project == nil {
			return errors.New("project not found")
		} else {
			if utils.PathExist(project.Path) {
				if err := os.RemoveAll(project.Path); err != nil {
					return err
				}
			}
			if err = tx.DeleteProject(id); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return err
	}

	return nil
}

func ReadMainFlow(id string) (string, string, error) {
	project, err := QueryProjectById(id)
	if err != nil {
		return "", "", err
	}
	mainFlowPath := project.Path + string(os.PathSeparator) + constants.BaseDir + string(os.PathSeparator) + constants.MainFlow
	if !utils.PathExist(mainFlowPath) {
		return "", "", nil
	}
	data, err := os.ReadFile(mainFlowPath)
	if err != nil {
		return "", "", err
	}
	return project.Name, string(data), nil
}

func SaveMainFlow(id string, data string) error {
	project, err := QueryProjectById(id)
	if err != nil {
		return err
	}
	mainFlowPath := project.Path + string(os.PathSeparator) + constants.BaseDir
	if !utils.PathExist(mainFlowPath) {
		if err = os.MkdirAll(mainFlowPath, fs.ModeDir); err != nil {
			return err
		}
	}
	err = os.WriteFile(mainFlowPath+string(os.PathSeparator)+constants.MainFlow, []byte(data), 0666)
	return err
}

func ReadSubFlow(id string, subId string) (string, error) {
	project, err := QueryProjectById(id)
	if err != nil {
		return "", err
	}
	subFlowPath := project.Path + string(os.PathSeparator) + constants.BaseDir + string(os.PathSeparator) + constants.DevDir + string(os.PathSeparator) + subId + ".flow"
	if !utils.PathExist(subFlowPath) {
		return "", nil
	}
	data, err := os.ReadFile(subFlowPath)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func SaveSubFlow(id string, subId string, data string) error {
	project, err := QueryProjectById(id)
	if err != nil {
		return err
	}
	projectPath := project.Path + string(os.PathSeparator) + constants.BaseDir
	devPath := projectPath + string(os.PathSeparator) + constants.DevDir
	if !utils.PathExist(devPath) {
		if err = os.MkdirAll(devPath, fs.ModeDir); err != nil {
			return err
		}
	}
	err = os.WriteFile(devPath+string(os.PathSeparator)+subId+".flow", []byte(data), 0666)
	if err != nil {
		return err
	}
	flow := Flow{}
	err = json.Unmarshal([]byte(data), &flow)
	if err != nil {
		return nil
	}
	err = flow.GeneratePythonCode(projectPath + string(os.PathSeparator) + subId + ".py")
	return err
}

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
	params["environment_variables"] = make(map[string]string)
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
	background, cancle := context.WithCancel(context.Background())
	go func() {
		monitorLog(ctx, background, logPath)
	}()
	_, err = command.Output()
	delete(processMap, subId)
	cancle()
	if err != nil {
		errStr := stderr.String()
		_, errStr, founded := strings.Cut(errStr, projectPath)
		log.Logger.Error(fmt.Sprintf("%t", founded))
		log.Logger.Error(errStr)
		return errors.New(errStr)
	}
	return nil
}

func DeubugSubFlow(ctx context.Context, id string, subId string) error {
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
	params["environment_variables"] = make(map[string]string)
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
	inPipe, err := command.StdinPipe()
	if err != nil {
		return err
	}
	outPipe, err := command.StdoutPipe()
	if err != nil {
		return err
	}
	breakpoints := getBreakpoints(projectPath, subId)
	go dealDebug(ctx, subId, inPipe, outPipe, breakpoints)
	background, cancle := context.WithCancel(context.Background())
	go monitorLog(ctx, background, logPath)
	debugChan = make(chan string)
	err = command.Run()
	delete(processMap, subId)
	cancle()
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

func TerminateSubFlow(id string, subId string) error {
	if command, ok := processMap[subId]; ok {
		err := sys_exec.KillProcess(command)
		return err
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
	defer inPipe.Close()
	defer outPipe.Close()
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
		output += string(line)
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
				if command, ok := processMap[filename]; ok {
					sys_exec.KillProcess(command)
				}
				break
			}
			err = writer.Flush()
			if err != nil {
				if command, ok := processMap[filename]; ok {
					sys_exec.KillProcess(command)
				}
				break
			}
			output = ""
		}
	}
}

func monitorLog(ctx, background context.Context, logPath string) {
	config := tail.Config{
		ReOpen:    true,                                 //重新打开
		Follow:    true,                                 //是否跟随
		Location:  &tail.SeekInfo{Offset: 0, Whence: 0}, //从哪儿开始读
		MustExist: false,                                //文件不存在不报错
		Poll:      true,                                 //
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
			runtime.EventsEmit(ctx, "log_event", line.Text)
		}
	}
}
