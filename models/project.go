package models

import (
	uuid "github.com/google/uuid"
)

type Project struct {
	Id          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Path        string    `json:"path"`
	Icon        string    `json:"icon"`
	Description string    `json:"description"`
	CreateTs    JsonTime  `json:"create_ts"`
	UpdateTs    JsonTime  `json:"update_ts"`
}
