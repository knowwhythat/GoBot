package http

import (
	"encoding/json"
	"errors"
	"github.com/duke-git/lancet/v2/netutil"
	"gobot/backend/forms"
	"net/http"
)

const ServerUrl = "http://111.173.105.169:8090/api"

func Login(username, password string) (*forms.LoginResponse, error) {
	header := http.Header{}
	header.Add("Content-Type", "application/json")
	loginReq := forms.LoginReq{
		Username: username,
		Password: password,
	}
	jsonData, err := json.Marshal(loginReq)
	if err != nil {
		return nil, err
	}
	request := &netutil.HttpRequest{
		RawURL:  ServerUrl + "/base/app_login",
		Method:  "POST",
		Headers: header,
		Body:    jsonData,
	}

	httpClient := netutil.NewHttpClient()
	resp, err := httpClient.SendRequest(request)
	if err != nil || resp.StatusCode != 200 {
		return nil, errors.New("net error")
	}
	var loginResp forms.RespData[forms.LoginResponse]
	err = httpClient.DecodeResponse(resp, &loginResp)
	if err != nil {
		return nil, err
	}
	if loginResp.Code != 0 {
		return nil, errors.New(loginResp.Msg)
	}
	return &loginResp.Data, nil
}
