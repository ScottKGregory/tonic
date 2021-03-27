package models

import (
	"time"
)

type Token struct {
	Token  string    `json:"token"`
	Expiry time.Time `json:"expiry"`
}
