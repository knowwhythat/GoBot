package services

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"gobot/constants"
	"gobot/log"
	"gobot/services/sys_exec"
	"gobot/utils"
	"io/fs"
	"net"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/spf13/viper"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

var windowsInspectCommand *exec.Cmd

var windowsInspectPort = 7777

func GetElementImage(projectId, element_id string) (string, error) {
	project, err := QueryProjectById(projectId)
	if err != nil {
		return "", err
	}
	imagePath := project.Path + string(os.PathSeparator) + constants.BaseDir + string(os.PathSeparator) + constants.DevDir + string(os.PathSeparator) + constants.SnapshotDir + string(os.PathSeparator) + element_id + ".png"
	if utils.PathExist(imagePath) {
		data, err := os.ReadFile(imagePath)
		if err != nil {
			return "", err
		}
		return base64.StdEncoding.EncodeToString(data), nil
	}
	return "", errors.New("图片文件不存在")
}

func startWindowsInspectCommand() error {
	for ; windowsInspectPort < 9999; windowsInspectPort++ {
		if PortCheck(windowsInspectPort) {
			break
		}
	}
	windowsInspectCommand = sys_exec.BuildCmd(viper.GetString("python.path"), "-u", "-B", "-m", "robot_desktop.windows_inspect", strconv.Itoa(windowsInspectPort))
	var stderr bytes.Buffer
	windowsInspectCommand.Stderr = &stderr
	err := windowsInspectCommand.Run()
	if err != nil {
		errStr := stderr.String()
		log.Logger.Error(errStr)
		return errors.New(errStr)
	}
	return nil
}

func StartPickWindowsElement(ctx context.Context, projectId string) (string, error) {
	project, err := QueryProjectById(projectId)
	if err != nil {
		return "", err
	}
	imagePath := project.Path + string(os.PathSeparator) + constants.BaseDir + string(os.PathSeparator) + constants.DevDir + string(os.PathSeparator) + constants.SnapshotDir
	if !utils.PathExist(imagePath) {
		if err = os.MkdirAll(imagePath, fs.ModeDir); err != nil {
			return "", err
		}
	}
	if windowsInspectCommand == nil {
		go startWindowsInspectCommand()
	}
	for retry := 1; retry < 10; retry++ {
		conn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://127.0.0.1:%s", strconv.Itoa(windowsInspectPort)), nil)
		if err != nil {
			time.Sleep(time.Second)
			continue
		}
		defer conn.Close()
		messageId := uuid.New().String()
		sendMessage := make(map[string]interface{})
		sendMessage["method"] = "highlight_control"
		sendMessage["message_id"] = messageId
		data := make(map[string]interface{})
		data["image_path"] = imagePath
		sendMessage["data"] = data
		request, err := json.Marshal(sendMessage)
		log.Logger.Info(string(request))
		if err != nil {
			return "", err
		}
		if err := conn.WriteMessage(1, request); err != nil {
			return "", err
		}
		runtime.WindowMinimise(ctx)
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				return "", nil
			}
			resp := make(map[string]interface{})
			if err = json.Unmarshal(message, &resp); err != nil {
				return "", nil
			}
			if resp["message_id"] == messageId {
				runtime.WindowMaximise(ctx)
				return string(message), nil
			}
		}
	}
	return "", errors.New("连接服务失败")
}

func StartCheckWindowsElement(ctx context.Context, paths string) (string, error) {
	if windowsInspectCommand == nil {
		go startWindowsInspectCommand()
	}
	for retry := 1; retry < 10; retry++ {
		conn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://127.0.0.1:%s", strconv.Itoa(windowsInspectPort)), nil)
		if err != nil {
			time.Sleep(time.Second)
			continue
		}
		defer conn.Close()
		messageId := uuid.New().String()
		sendMessage := make(map[string]interface{})
		sendMessage["method"] = "check_control"
		sendMessage["message_id"] = messageId
		data := make(map[string]interface{})
		data["paths"] = paths
		sendMessage["data"] = data
		request, err := json.Marshal(sendMessage)
		log.Logger.Info(string(request))
		if err != nil {
			return "", err
		}
		if err := conn.WriteMessage(1, request); err != nil {
			return "", err
		}
		runtime.WindowMinimise(ctx)
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				return "", nil
			}
			resp := make(map[string]interface{})
			if err = json.Unmarshal(message, &resp); err != nil {
				return "", nil
			}
			if resp["message_id"] == messageId {
				runtime.WindowMaximise(ctx)
				return string(message), nil
			}
		}
	}
	return "", errors.New("连接服务失败")
}

func GetWindowsElementList(parentId string) (string, error) {
	if windowsInspectCommand == nil {
		go startWindowsInspectCommand()
	}
	for retry := 1; retry < 10; retry++ {
		conn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://127.0.0.1:%s", strconv.Itoa(windowsInspectPort)), nil)
		if err != nil {
			time.Sleep(time.Second)
			continue
		}
		defer conn.Close()
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
		log.Logger.Info(string(request))
		if err != nil {
			return "", err
		}
		if err := conn.WriteMessage(1, request); err != nil {
			return "", err
		}
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				return "", nil
			}
			resp := make(map[string]interface{})
			if err = json.Unmarshal(message, &resp); err != nil {
				return "", nil
			}
			if resp["message_id"] == messageId {
				return string(message), nil
			}
		}
	}
	return "", errors.New("连接服务失败")
}

func HighlightCurrentElement(controlId string) error {
	if windowsInspectCommand == nil {
		go startWindowsInspectCommand()
	}
	for retry := 1; retry < 10; retry++ {
		conn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://127.0.0.1:%s", strconv.Itoa(windowsInspectPort)), nil)
		if err != nil {
			time.Sleep(time.Second)
			continue
		}
		defer conn.Close()
		messageId := uuid.New().String()
		sendMessage := make(map[string]interface{})
		sendMessage["message_id"] = messageId
		data := make(map[string]interface{})
		sendMessage["method"] = "highlight_current_control"
		data["control_id"] = controlId
		sendMessage["data"] = data
		request, err := json.Marshal(sendMessage)
		log.Logger.Info(string(request))
		if err != nil {
			return err
		}
		if err := conn.WriteMessage(1, request); err != nil {
			return err
		}
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				return nil
			}
			resp := make(map[string]interface{})
			if err = json.Unmarshal(message, &resp); err != nil {
				return nil
			}
			if resp["message_id"] == messageId {
				return nil
			}
		}
	}
	return errors.New("连接服务失败")
}

func GetSelectedWindowsElement(ctx context.Context, projectId, controlId string) (string, error) {
	project, err := QueryProjectById(projectId)
	if err != nil {
		return "", err
	}
	imagePath := project.Path + string(os.PathSeparator) + constants.BaseDir + string(os.PathSeparator) + constants.DevDir + string(os.PathSeparator) + constants.SnapshotDir
	if !utils.PathExist(imagePath) {
		if err = os.MkdirAll(imagePath, fs.ModeDir); err != nil {
			return "", err
		}
	}
	if windowsInspectCommand == nil {
		go startWindowsInspectCommand()
	}
	for retry := 1; retry < 10; retry++ {
		conn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://127.0.0.1:%s", strconv.Itoa(windowsInspectPort)), nil)
		if err != nil {
			time.Sleep(time.Second)
			continue
		}
		defer conn.Close()
		messageId := uuid.New().String()
		sendMessage := make(map[string]interface{})
		sendMessage["message_id"] = messageId
		data := make(map[string]interface{})
		sendMessage["method"] = "get_select_control"
		data["image_path"] = imagePath
		data["control_id"] = controlId
		sendMessage["data"] = data
		request, err := json.Marshal(sendMessage)
		log.Logger.Info(string(request))
		if err != nil {
			return "", err
		}
		runtime.WindowMinimise(ctx)
		if err := conn.WriteMessage(1, request); err != nil {
			return "", err
		}
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				return "", nil
			}
			resp := make(map[string]interface{})
			if err = json.Unmarshal(message, &resp); err != nil {
				return "", nil
			}
			if resp["message_id"] == messageId {
				runtime.WindowMaximise(ctx)
				return string(message), nil
			}
		}
	}
	return "", errors.New("连接服务失败")
}

func SaveWindowsElement(id, elementId, selector string) error {
	project, err := QueryProjectById(id)
	if err != nil {
		return err
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
	elementData := make([]interface{}, 0)
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

func StopWindowsInspectCommand() {
	if windowsInspectCommand != nil {
		sys_exec.KillProcess(windowsInspectCommand)
		windowsInspectCommand = nil
	}
}

func PortCheck(port int) bool {
	l, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%s", strconv.Itoa(port)))

	if err != nil {
		return false
	}
	defer l.Close()
	return true
}
