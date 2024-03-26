package models

import (
	"time"
)

type Project struct {
	Id          string    `json:"id"`
	Name        string    `json:"name"`
	Path        string    `json:"path"`
	Icon        string    `json:"icon"`
	IsFlow      bool      `json:"isFlow"`
	Description string    `json:"description"`
	CreateTs    time.Time `json:"createTs"`
	UpdateTs    time.Time `json:"updateTs"`
}
