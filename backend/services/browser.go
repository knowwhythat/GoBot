package services

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

var wsConn *websocket.Conn

func StartPick(ctx context.Context) (string, error) {
	if wsConn == nil {
		conn, _, err := websocket.DefaultDialer.Dial("ws://127.0.0.1:3000/ws", nil)
		if err == nil {
			wsConn = conn
		} else {
			return "", err
		}
	}
	messageId := uuid.New().String()
	sendMessage := make(map[string]string)
	sendMessage["message"] = "start"
	sendMessage["id"] = messageId
	request, err := json.Marshal(sendMessage)
	if err != nil {
		return "", err
	}
	if err := wsConn.WriteMessage(1, request); err != nil {
		_ = wsConn.Close()
		wsConn = nil
		return "", err
	}
	runtime.WindowMinimise(ctx)
	_, message, err := wsConn.ReadMessage()
	if err != nil {
		return "", nil
	}
	resp := make(map[string]interface{})
	if err = json.Unmarshal(message, &resp); err != nil {
		return "", nil
	}
	runtime.WindowMaximise(ctx)
	return string(message), nil
}

func StartCheck(ctx context.Context, frame, selector string) (string, error) {
	if wsConn == nil {
		conn, _, err := websocket.DefaultDialer.Dial("ws://127.0.0.1:3000/ws", nil)
		if err == nil {
			wsConn = conn
		} else {
			return "", err
		}
	}
	messageId := uuid.New().String()
	sendMessage := make(map[string]string)
	sendMessage["message"] = "highlight"
	sendMessage["frame"] = frame
	sendMessage["xpath"] = selector
	sendMessage["id"] = messageId
	request, err := json.Marshal(sendMessage)
	if err != nil {
		return "", err
	}
	if err := wsConn.WriteMessage(1, request); err != nil {
		return "", err
	}
	runtime.WindowMinimise(ctx)
	_, message, err := wsConn.ReadMessage()
	if err != nil {
		return "", nil
	}
	resp := make(map[string]interface{})
	if err = json.Unmarshal(message, &resp); err != nil {
		return "", nil
	}
	runtime.WindowMaximise(ctx)
	return string(message), nil
}
