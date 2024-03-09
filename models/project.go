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
	IsFlow      bool      `json:"isFlow"`
	Description string    `json:"description"`
	CreateTs    time.Time `json:"createTs"`
	UpdateTs    time.Time `json:"updateTs"`
}
