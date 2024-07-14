package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"gobot/backend/config"
	"gobot/backend/constants"
	"gobot/backend/forms"
	"gobot/backend/log"
	"gobot/backend/services/sys_exec"
	"gobot/backend/utils"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"
)

var wsConn *websocket.Conn

var monitorMap = make(map[string]*exec.Cmd)

func main() {
	if err := config.InitConfig("./config.yml"); err != nil {
		println("Init config error:" + err.Error())
		panic(err)
	}
	log.Init("virtual_host")
	log.Logger.Info("Start virtual host")
	port := 3500
	for ; port < 4000; port++ {
		if utils.PortCheck(port) {
			break
		}
	}
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exePath := filepath.Dir(ex)
	err = os.WriteFile(filepath.Join(exePath, "virtual_host_port"), []byte(strconv.Itoa(port)), os.ModePerm)
	if err != nil {
		panic(err)
	}
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", serveWs)
	server := &http.Server{
		Addr:              fmt.Sprintf("127.0.0.1:%s", strconv.Itoa(port)),
		ReadHeaderTimeout: 3 * time.Second,
	}
	log.Logger.Logger.Info().Msgf("Start virtual host on port %d", port)
	err = server.ListenAndServe()
	if err != nil {
		log.Logger.Logger.Error().Msg("服务启动失败")
		return
	}
}

var upgrader = websocket.Upgrader{CheckOrigin: checkOrigin}

func checkOrigin(r *http.Request) bool {
	return true
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	_, err := http.ResponseWriter.Write(w, []byte("success"))
	if err != nil {
		return
	}
}

func serveWs(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Logger.Logger.Error().AnErr("upgrade:", err)
		return
	}
	if wsConn != nil {
		err := wsConn.Close()
		if err != nil {
		}
	}
	wsConn = conn
	go listenRequest()
	log.Logger.Logger.Info().Msg("有新连接接入")
}

func listenRequest() {
	for {
		if wsConn != nil {
			_, message, err := wsConn.ReadMessage()
			log.Logger.Logger.Info().Msg("read from control:" + string(message))
			if err != nil {
				log.Logger.Logger.Error().Err(err)
				_ = wsConn.Close()
				wsConn = nil
				break
			}
			data := &forms.Transfer{}
			err = json.Unmarshal(message, data)
			if err != nil {
				log.Logger.Logger.Error().Err(err)
				continue
			}
			switch data.Action {
			case constants.StartAction:
				startFLow := &forms.StartFlow{}
				err = json.Unmarshal([]byte(data.Data), startFLow)
				if err != nil {
					log.Logger.Logger.Error().Err(err)
					continue
				}
				go StartInVirtualDesk(startFLow)
			case constants.StopAction:
				stopFlow := &forms.StopFlow{}
				err = json.Unmarshal([]byte(data.Data), stopFlow)
				if err != nil {
					log.Logger.Logger.Error().Err(err)
					continue
				}
				go TerminateVirtualDeskFlow(stopFlow.InstanceId)
			}
		} else {
			break
		}
	}
}

func StartInVirtualDesk(data *forms.StartFlow) {
	proc := sys_exec.BuildCmd(data.ExePath, "-u", "-m", data.Module, data.Param)
	monitorMap[data.InstanceId] = proc
	log.Logger.Logger.Info().Msgf("启动进程%s", data.InstanceId)
	var stderr bytes.Buffer
	proc.Stderr = &stderr
	_, err := proc.Output()
	log.Logger.Logger.Info().Msgf("%s执行完成", data.InstanceId)
	delete(monitorMap, data.InstanceId)
	if err != nil {
		log.Logger.Logger.Error().Err(err)
		errStr := stderr.String()
		log.Logger.Logger.Error().Msg(errStr)
		writeResponseMessage(constants.ExecuteResultAction, data.InstanceId, err.Error(), false)
	} else {
		writeResponseMessage(constants.ExecuteResultAction, data.InstanceId, "success", true)
	}
}

func TerminateVirtualDeskFlow(instanceId string) {
	if proc, ok := monitorMap[instanceId]; ok {
		err := sys_exec.KillProcess(proc)
		if err != nil {
			log.Logger.Logger.Error().Err(err)
			writeResponseMessage(constants.StopResultAction, instanceId, err.Error(), false)
		} else {
			log.Logger.Logger.Info().Msgf("成功终止进程%s", instanceId)
			writeResponseMessage(constants.StopResultAction, instanceId, "success", true)
		}
	} else {
		log.Logger.Logger.Error().Msgf("没有找到对应的进程%s", instanceId)
		writeResponseMessage(constants.StopResultAction, instanceId, "没有找到对应的进程", false)
	}
}

func writeResponseMessage(action string, instanceId string, message string, success bool) {
	resp := forms.Transfer{Action: action}
	data := forms.Result{
		InstanceId: instanceId,
		Message:    message,
		Success:    success,
	}
	respData, err := json.Marshal(data)
	if err != nil {
		log.Logger.Logger.Error().Err(err)
		return
	}
	resp.Data = string(respData)
	respMessage, err := json.Marshal(resp)
	if err != nil {
		log.Logger.Logger.Error().Err(err)
		return
	}
	err = wsConn.WriteMessage(websocket.TextMessage, respMessage)
	if err != nil {
		log.Logger.Logger.Error().Err(err)
		return
	}
}
