package handlers

import (
	"messaging-application/servers/gateway/models/users"
	"messaging-application/servers/gateway/sessions"
)

type Context struct {
	Secret       string         `json:"secret"`
	SessionStore sessions.Store `json:"sessionStore"`
	UserStore    users.Store    `json:"userStore"`
}
