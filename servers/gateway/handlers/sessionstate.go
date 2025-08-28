package handlers

import (
	"messaging-application/servers/gateway/models/users"
	"time"
)

type SessionState struct {
	Start time.Time  `json:"start"`
	User  users.User `json:"user"`
}
