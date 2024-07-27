package services

import (
	"encoding/json"
	"errors"
	"github.com/duke-git/lancet/v2/netutil"
	"github.com/google/uuid"
	"gobot/backend/forms"
	httpservice "gobot/backend/http"
	"gobot/backend/log"
	"net/http"
	"time"
)

var wxParam = ""

var wxUuid = ""

func GetQrCode() (string, error) {
	header := http.Header{}
	header.Add("Content-Type", "application/json")
	wxParam = uuid.New().String()
	req := forms.GenerateCode{
		AppToken:  "AT_7DGDlPV7P7GjSuLTvnbexDsc5NUTtml6",
		Extra:     wxParam,
		ValidTime: 1800,
	}
	jsonData, err := json.Marshal(req)
	if err != nil {
		return "", err
	}
	request := &netutil.HttpRequest{
		RawURL:  "https://wxpusher.zjiecode.com/api/fun/create/qrcode",
		Method:  "POST",
		Headers: header,
		Body:    jsonData,
	}

	httpClient := netutil.NewHttpClient()
	resp, err := httpClient.SendRequest(request)
	if err != nil || resp.StatusCode != 200 {
		return "", errors.New("net error")
	}
	var codeResp forms.WxPusherResp[forms.CodeData]
	err = httpClient.DecodeResponse(resp, &codeResp)
	if err != nil {
		return "", err
	}
	if codeResp.Code != 1000 {
		return "", errors.New(codeResp.Msg)
	}
	go queryWxPusherId(codeResp.Data.Code)
	return codeResp.Data.URL, nil
}

func Register(username, password string) error {
	if wxUuid == "" {
		return errors.New("请先用微信扫描左侧二维码")
	}
	return httpservice.Register(username, password, wxUuid)
}

func queryWxPusherId(code string) {
	for {
		time.Sleep(time.Second * 10)
		request := &netutil.HttpRequest{
			RawURL: "https://wxpusher.zjiecode.com/api/fun/scan-qrcode-uid?code=" + code,
			Method: "GET",
		}
		httpClient := netutil.NewHttpClient()
		resp, err := httpClient.SendRequest(request)
		if err != nil || resp.StatusCode != 200 {
			log.Logger.Logger.Error().Err(err).Msg("获取微信推送ID失败")
			continue
		}
		var userResp forms.WxPusherResp[string]
		err = httpClient.DecodeResponse(resp, &userResp)
		if err != nil {
			log.Logger.Logger.Error().Err(err).Msg("解析微信推送ID失败")
			continue
		}
		if userResp.Code != 1000 {
			log.Logger.Logger.Info().Msg(userResp.Msg)
			continue
		}
		if userResp.Data != "" {
			wxUuid = userResp.Data
			break
		}
	}
}
