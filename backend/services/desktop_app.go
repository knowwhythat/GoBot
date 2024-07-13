package services

import (
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
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/spf13/viper"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

var windowsInspectCommand *exec.Cmd

var windowsInspectPort = 9999

var windowsInspectConn *websocket.Conn

func GetElementImage(projectId, elementId string) (string, error) {
	project, err := QueryProjectById(projectId)
	if err != nil {
		return "", err
	}
	imagePath := filepath.Join(project.Path, constants.BaseDir, constants.DevDir, constants.SnapshotDir, elementId+".png")
	if utils.PathExist(imagePath) {
		data, err := os.ReadFile(imagePath)
		if err != nil {
			return "", err
		}
		return base64.StdEncoding.EncodeToString(data), nil
	}
	return "", errors.New("图片文件不存在")
}

func startWindowsInspectCommand() {
	for ; windowsInspectPort < 9999; windowsInspectPort++ {
		if PortCheck(windowsInspectPort) {
			break
		}
	}
	windowsInspectCommand = sys_exec.BuildCmd(viper.GetString(constants.ConfigPythonPath), "-u", "-B", "-m", "robot_desktop.windows_inspect", strconv.Itoa(windowsInspectPort))
	var stderr bytes.Buffer
	windowsInspectCommand.Stderr = &stderr
	err := windowsInspectCommand.Run()
	if err != nil {
		windowsInspectCommand = nil
		errStr := stderr.String()
		log.Logger.Logger.Error().Msg(errStr)
		log.Logger.Logger.Error().Err(err)
	}
}

func checkWindowsInspectConn(ctx context.Context) error {
	if windowsInspectCommand == nil {
		//go startWindowsInspectCommand()
	}
	for i := 0; i < 5; i++ {
		if windowsInspectConn == nil {
			log.Logger.Logger.Info().Msgf("windowsInspectPort: %d", windowsInspectPort)
			conn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://127.0.0.1:%d/ws", windowsInspectPort), nil)
			if err == nil {
				windowsInspectConn = conn
				go listenInspectResponse(ctx)
				return nil
			} else {
				time.Sleep(time.Second)
				continue
			}
		} else {
			return nil
		}
	}
	return errors.New("连接拾取引擎失败")
}

func StartPickWindowsElement(ctx context.Context) error {
	err := checkWindowsInspectConn(ctx)
	if err != nil {
		return err
	}
	messageId := uuid.New().String()
	sendMessage := make(map[string]interface{})
	sendMessage["method"] = "highlight_control"
	sendMessage["message_id"] = messageId
	data := make(map[string]interface{})
	sendMessage["data"] = data
	request, err := json.Marshal(sendMessage)
	log.Logger.Logger.Info().Msg(string(request))
	if err != nil {
		return err
	}
	if err := windowsInspectConn.WriteMessage(1, request); err != nil {
		_ = windowsInspectConn.Close()
		windowsInspectConn = nil
		return err
	}
	return nil
}

func StartCheckWindowsElement(ctx context.Context, paths string) error {
	err := checkWindowsInspectConn(ctx)
	if err != nil {
		return err
	}
	messageId := uuid.New().String()
	sendMessage := make(map[string]interface{})
	sendMessage["method"] = "check_control"
	sendMessage["message_id"] = messageId
	data := make(map[string]interface{})
	data["paths"] = paths
	sendMessage["data"] = data
	request, err := json.Marshal(sendMessage)
	log.Logger.Logger.Info().Msg(string(request))
	if err != nil {
		return err
	}
	if err := windowsInspectConn.WriteMessage(1, request); err != nil {
		_ = windowsInspectConn.Close()
		windowsInspectConn = nil
		return err
	}
	return nil
}

func GetWindowsElementList(ctx context.Context, parentId string) error {
	err := checkWindowsInspectConn(ctx)
	if err != nil {
		return err
	}
	messageId := uuid.New().String()
	sendMessage := make(map[string]interface{})
	sendMessage["message_id"] = messageId
	data := make(map[string]interface{})
	if parentId == "" {
		sendMessage["method"] = "get_root"
	} else {
		sendMessage["method"] = "get_children"
		data["control_id"] = parentId
	}
	sendMessage["data"] = data
	request, err := json.Marshal(sendMessage)
	log.Logger.Logger.Info().Msg(string(request))
	if err != nil {
		return err
	}
	if err := windowsInspectConn.WriteMessage(1, request); err != nil {
		_ = windowsInspectConn.Close()
		windowsInspectConn = nil
		return err
	}
	return nil
}

func HighlightCurrentElement(ctx context.Context, controlId string) error {
	err := checkWindowsInspectConn(ctx)
	if err != nil {
		return err
	}
	messageId := uuid.New().String()
	sendMessage := make(map[string]interface{})
	sendMessage["message_id"] = messageId
	data := make(map[string]interface{})
	sendMessage["method"] = "highlight_current_control"
	data["control_id"] = controlId
	sendMessage["data"] = data
	request, err := json.Marshal(sendMessage)
	log.Logger.Logger.Info().Msg(string(request))
	if err != nil {
		return err
	}
	if err := windowsInspectConn.WriteMessage(1, request); err != nil {
		_ = windowsInspectConn.Close()
		windowsInspectConn = nil
		return err
	}
	return nil
}

func GetSelectedWindowsElement(ctx context.Context, controlId string) error {
	err := checkWindowsInspectConn(ctx)
	if err != nil {
		return err
	}
	messageId := uuid.New().String()
	sendMessage := make(map[string]interface{})
	sendMessage["message_id"] = messageId
	data := make(map[string]interface{})
	sendMessage["method"] = "get_select_control"
	data["control_id"] = controlId
	sendMessage["data"] = data
	request, err := json.Marshal(sendMessage)
	log.Logger.Logger.Info().Msg(string(request))
	if err != nil {
		return err
	}
	if err := windowsInspectConn.WriteMessage(1, request); err != nil {
		_ = windowsInspectConn.Close()
		windowsInspectConn = nil
		return err
	}
	return nil
}

func listenInspectResponse(ctx context.Context) {
	for windowsInspectConn != nil {
		_, message, err := windowsInspectConn.ReadMessage()
		if err != nil {
			log.Logger.Logger.Error().Err(err)
			return
		}
		log.Logger.Logger.Info().Msgf("拾取引擎响应: %s", string(message))
		runtime.EventsEmit(ctx, constants.WindowsEvent, string(message))
	}
}
func SaveWindowsElement(ctx context.Context, id, elementId, image, selector string) error {
	project, err := QueryProjectById(id)
	if err != nil {
		return err
	}
	if len(image) > 0 {
		imageDirPath := project.Path + string(os.PathSeparator) + constants.BaseDir + string(os.PathSeparator) + constants.DevDir + string(os.PathSeparator) + constants.SnapshotDir + string(os.PathSeparator)
		if !utils.PathExist(imageDirPath) {
			_ = os.MkdirAll(imageDirPath, os.ModePerm)
		}
		imageData, err := base64.StdEncoding.DecodeString(image)
		if err != nil {
			return err
		}
		_ = os.WriteFile(imageDirPath+string(os.PathSeparator)+elementId+".png", imageData, 0666)
	}

	selectorPath := project.Path + string(os.PathSeparator) + constants.BaseDir + string(os.PathSeparator) + constants.Selector
	selectorData := make(map[string]interface{})
	if utils.PathExist(selectorPath) {
		data, err := os.ReadFile(selectorPath)
		if err != nil {
			return err
		}
		err = json.Unmarshal(data, &selectorData)
		if err != nil {
			return err
		}
	}
	elementData := make(map[string]interface{})
	err = json.Unmarshal([]byte(selector), &elementData)
	if err != nil {
		return err
	}
	selectorData[elementId] = elementData
	saveData, err := json.Marshal(selectorData)
	if err != nil {
		return err
	}
	err = os.WriteFile(selectorPath, saveData, 0666)
	runtime.EventsEmit(ctx, "windows_element_change", "")
	return err
}

func GetWindowsElement(id, elementId string) (string, error) {
	project, err := QueryProjectById(id)
	if err != nil {
		return "", err
	}
	selectorPath := project.Path + string(os.PathSeparator) + constants.BaseDir + string(os.PathSeparator) + constants.Selector
	selectorData := make(map[string]interface{})
	if utils.PathExist(selectorPath) {
		data, err := os.ReadFile(selectorPath)
		if err != nil {
			return "", err
		}
		err = json.Unmarshal(data, &selectorData)
		if err != nil {
			return "", err
		}
	}
	if selector, ok := selectorData[elementId]; ok {
		selectorByte, err := json.Marshal(selector)
		return string(selectorByte), err
	}
	return "", nil
}

func RemoveWindowsElement(ctx context.Context, id, elementId string) error {
	project, err := QueryProjectById(id)
	if err != nil {
		return err
	}
	imagePath := project.Path + string(os.PathSeparator) + constants.BaseDir + string(os.PathSeparator) + constants.DevDir + string(os.PathSeparator) + constants.SnapshotDir + string(os.PathSeparator) + elementId + ".png"
	if utils.PathExist(imagePath) {
		_ = os.Remove(imagePath)
	}
	selectorPath := project.Path + string(os.PathSeparator) + constants.BaseDir + string(os.PathSeparator) + constants.Selector
	selectorData := make(map[string]interface{})
	if utils.PathExist(selectorPath) {
		data, err := os.ReadFile(selectorPath)
		if err != nil {
			return err
		}
		err = json.Unmarshal(data, &selectorData)
		if err != nil {
			return err
		}
	}
	if _, ok := selectorData[elementId]; ok {
		delete(selectorData, elementId)
		saveData, err := json.Marshal(selectorData)
		if err != nil {
			return err
		}
		err = os.WriteFile(selectorPath, saveData, 0666)
		runtime.EventsEmit(ctx, "windows_element_change", "")
		return err
	}
	return nil
}

func GetProjectWindowsElements(id string) (string, error) {
	project, err := QueryProjectById(id)
	if err != nil {
		return "", err
	}
	selectorPath := project.Path + string(os.PathSeparator) + constants.BaseDir + string(os.PathSeparator) + constants.Selector
	if utils.PathExist(selectorPath) {
		data, err := os.ReadFile(selectorPath)
		if err != nil {
			return "", err
		}
		return string(data), nil
	}
	return "", nil
}

func StopWindowsInspectCommand() {
	if windowsInspectCommand != nil {
		_ = sys_exec.KillProcess(windowsInspectCommand)
		windowsInspectCommand = nil
	}
}

func PortCheck(port int) bool {
	l, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%s", strconv.Itoa(port)))

	if err != nil {
		return false
	}
	defer func(l net.Listener) {
		_ = l.Close()
	}(l)
	return true
}
