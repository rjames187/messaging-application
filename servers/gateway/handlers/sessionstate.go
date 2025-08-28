package handlers

import (
	"messaging-application/servers/gateway/models/users"
	"time"
)

type Session struct {
	Start time.Time  `json:"start"`
	User  users.User `json:"user"`
}
