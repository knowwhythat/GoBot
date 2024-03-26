package models

import (
	"time"
)

type Schedule struct {
	Id               string    `json:"id"`
	ProjectId        string    `json:"projectId"`
	ProjectVersionId string    `json:"projectVersionId"`
	Name             string    `json:"name"`
	Cron             string    `json:"cron"`
	Enabled          bool      `json:"enabled"`
	Desc             string    `json:"desc"`
	CreateTs         time.Time `json:"createTs"`
	UpdateTs         time.Time `json:"updateTs"`
}
