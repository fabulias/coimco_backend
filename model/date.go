package model

import "time"

type Date struct {
	Start time.Time `json:"start" binding:"required"`
	End   time.Time `json:"end" binding:"required"`
}
