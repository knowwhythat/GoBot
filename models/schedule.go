package models

import (
	"time"

	uuid "github.com/google/uuid"
	"github.com/robfig/cron/v3"
)

type Schedule struct {
	Id               uuid.UUID    `json:"id"`
	ProjectId        uuid.UUID    `json:"project_id"`
	ProjectVersionId uuid.UUID    `json:"project_version_id"`
	Cron             string       `json:"cron"`
	EntryId          cron.EntryID `json:"entry_id"`
	Cmd              string       `json:"cmd"`
	Enabled          bool         `json:"enabled"`
	Description      string       `json:"description"`
	CreateTs         time.Time    `json:"create_ts"`
	UpdateTs         time.Time    `json:"update_ts"`
}
