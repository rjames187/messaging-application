package handlers

import (
	"encoding/json"
	"messaging-application/servers/gateway/models/users"
	"time"
)

type SessionState struct {
	Start time.Time  `json:"start"`
	User  users.User `json:"user"`
}

func GetSerializedSessionState(user *users.User) (string, error) {
	state := SessionState{
		Start: time.Now(),
		User:  *user,
	}

	serialized, err := json.Marshal(state)
	if err != nil {
		return "", err
	}

	return string(serialized), nil
}
