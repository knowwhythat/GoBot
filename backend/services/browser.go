package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"gobot/backend/constants"
	"gobot/backend/log"
	"gobot/backend/services/env"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var wsConn *websocket.Conn

func StartPick(ctx context.Context) error {
	if checkConn(ctx) != nil {
		return errors.New("连接浏览器插件失败，请确定是否正确配置浏览器插件")
	}
	messageId := uuid.New().String()
	sendMessage := make(map[string]string)
	sendMessage["message"] = "start"
	sendMessage["id"] = messageId
	request, err := json.Marshal(sendMessage)
	if err != nil {
		return err
	}
	log.Logger.Logger.Info().Msgf("向浏览器插件发送: %s", request)
	err = wsConn.WriteMessage(1, request)
	if err != nil {
		_ = wsConn.Close()
		wsConn = nil
		if checkConn(ctx) != nil {
			return errors.New("连接浏览器插件失败，请确定是否正确配置浏览器插件")
		}
		err = wsConn.WriteMessage(1, request)
	}
	return err
}

func StartCheck(ctx context.Context, frame, selector string) error {
	if checkConn(ctx) != nil {
		return errors.New("连接浏览器插件失败，请确定是否正确配置浏览器插件")
	}
	messageId := uuid.New().String()
	sendMessage := make(map[string]string)
	sendMessage["message"] = "highlight"
	sendMessage["frame"] = frame
	sendMessage["xpath"] = selector
	sendMessage["id"] = messageId
	request, err := json.Marshal(sendMessage)
	if err != nil {
		return err
	}
	log.Logger.Logger.Info().Msgf("向浏览器插件发送: %s", request)
	err = wsConn.WriteMessage(1, request)
	if err != nil {
		_ = wsConn.Close()
		wsConn = nil
		if checkConn(ctx) != nil {
			return errors.New("连接浏览器插件失败，请确定是否正确配置浏览器插件")
		}
		err = wsConn.WriteMessage(1, request)
	}
	return err
}

func checkConn(ctx context.Context) error {
	for i := 0; i < 3; i++ {
		if wsConn == nil {
			data, err := os.ReadFile(filepath.Join(env.GetExecutablePath(), "chrome", "nativehost_port"))
			if err != nil {
				continue
			}
			log.Logger.Logger.Info().Msgf("nativehost_port: %s", string(data))
			conn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://127.0.0.1:%s/ws", strings.TrimSpace(string(data))), nil)
			if err == nil {
				wsConn = conn
				go listenResponse(ctx)
				return nil
			} else {
				continue
			}
		} else {
			return nil
		}
	}
	return errors.New("连接浏览器插件失败")
}

func listenResponse(ctx context.Context) {
	for wsConn != nil {
		_, message, err := wsConn.ReadMessage()
		if err != nil {
			return
		}
		log.Logger.Logger.Info().Msgf("浏览器插件响应: %s", string(message))
		runtime.EventsEmit(ctx, constants.BrowserEvent, string(message))
	}
}
