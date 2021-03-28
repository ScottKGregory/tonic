package models

import (
	"time"
)

// Token represents an API response containing an access token
type Token struct {
	Token  string    `json:"token"`
	Expiry time.Time `json:"expiry"`
}
