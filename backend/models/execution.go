package models

import (
	"time"
)

type Execution struct {
	Id            string    `json:"id"`
	ProjectId     string    `json:"projectId"`
	SubFlowId     string    `json:"subFlowId"`
	ExecuteResult int       `json:"executeResult"` // 0:失败;1:成功;3:手动停止
	ErrorMsg      string    `json:"errorMsg"`
	StartTs       time.Time `json:"startTs"`
	EndTs         time.Time `json:"endTs"`
}
