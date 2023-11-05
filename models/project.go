package models

import (
	"time"

	uuid "github.com/google/uuid"
)

type Project struct {
	Id          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Path        string    `json:"path"`
	Icon        string    `json:"icon"`
	Description string    `json:"description"`
	CreateTs    time.Time `json:"create_ts"`
	UpdateTs    time.Time `json:"update_ts"`
}
