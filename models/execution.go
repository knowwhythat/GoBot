package models

import (
	"time"

	uuid "github.com/google/uuid"
)

type Execution struct {
	Id            uuid.UUID `json:"id"`
	ProjectId     uuid.UUID `json:"project_id"`
	SubFlowId     string    `json:"sub_flow_id"`
	ExecuteResult int       `json:"execute_result"` // 0:失败;1:成功;3:手动停止
	ErrorMsg      string    `json:"error_msg"`
	StartTs       time.Time `json:"start_ts"`
	EndTs         time.Time `json:"end_ts"`
}
