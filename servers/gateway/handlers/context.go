package handlers

import (
	"messaging-application/servers/gateway/models/users"
	"messaging-application/servers/gateway/sessions"
)

type HandlerContext struct {
	Secret       string         `json:"secret"`
	SessionStore sessions.Store `json:"sessionStore"`
	UserStore    users.Store    `json:"userStore"`
}

func NewHandlerContext(secret string, sessionStore sessions.Store, userStore users.Store) *HandlerContext {
	return &HandlerContext{
		Secret:       secret,
		SessionStore: sessionStore,
		UserStore:    userStore,
	}
}
