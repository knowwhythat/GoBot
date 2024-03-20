package forms

import (
	uuid "github.com/google/uuid"
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
