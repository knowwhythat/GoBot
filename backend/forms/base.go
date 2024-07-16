package forms

import (
	uuid "github.com/google/uuid"
	"os/exec"
	"time"
)

type BaseForm struct {
}

type PageForm struct {
	BaseForm

	PageNum  int `form:"page_num" json:"page_num"`
	PageSize int `form:"page_size" json:"page_size"`
}

type QueryForm struct {
	PageForm

	Query string `form:"query" json:"query"`
}

func (page *PageForm) Range() (start int, end int) {
	if page.PageNum <= 0 || page.PageSize <= 0 {
		return 0, 0
	}
	start = (page.PageNum - 1) * page.PageSize
	return start, start + page.PageSize - 1
}

type BatchForm struct {
	BaseForm

	Ids []uuid.UUID `form:"ids" json:"ids"`
}

type LoginForm struct {
	Username   string `form:"username" json:"username"`
	Pwd        string `form:"pwd" json:"pwd"`
	RememberMe bool   `form:"rememberMe" json:"rememberMe"`
	AutoLogin  bool   `form:"autoLogin" json:"autoLogin"`
}

type ExecutionForm struct {
	Id            string    `json:"id"`
	ProjectId     string    `json:"projectId"`
	ProjectName   string    `json:"projectName"`
	SubFlowId     string    `json:"subFlowId"`
	TriggerType   string    `json:"triggerType"`   //触发方式,手动;定时任务
	ExecuteResult int       `json:"executeResult"` // 0:失败;1:成功;3:手动停止
	ErrorMsg      string    `json:"errorMsg"`
	StartTs       time.Time `json:"startTs"`
	EndTs         time.Time `json:"endTs"`
}

type RunningInstance struct {
	Id          string    `form:"id" json:"id"`
	ProjectId   string    `form:"projectId" json:"projectId"`
	ProjectName string    `form:"projectName" json:"projectName"`
	StartTs     time.Time `form:"startTs" json:"startTs"`
	TriggerType string    `form:"triggerType" json:"triggerType"`
	IsVirtual   bool      `form:"isVirtual" json:"isVirtual"`
	Process     *exec.Cmd `form:"-" json:"-"`
}

type PipPackage struct {
	Name         string `form:"name" json:"name"`
	UpgradePip   bool   `form:"upgradePip" json:"upgradePip"`
	Upgrade      bool   `form:"upgrade" json:"upgrade"`
	InputVersion bool   `form:"inputVersion" json:"inputVersion"`
	Version      string `form:"version" json:"version"`
	InputMirror  bool   `form:"inputMirror" json:"inputMirror"`
	MirrorUrl    string `form:"mirrorUrl" json:"mirrorUrl"`
}

type Package struct {
	Name    string `form:"name" json:"name"`
	Version string `form:"version" json:"version"`
}

type ImageData struct {
	Id     string  `json:"id"`
	Name   string  `json:"name"`
	Image  string  `json:"image"`
	Left   float32 `json:"left"`
	Top    float32 `json:"top"`
	Right  float32 `json:"right"`
	Bottom float32 `json:"bottom"`
}
