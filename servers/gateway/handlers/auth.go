package handlers

import (
	"encoding/json"
	"messaging-application/servers/gateway/models/users"
	"messaging-application/servers/gateway/sessions"
	"net/http"
)

func (ctx *Context) UsersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Content-Type must be application/json", http.StatusUnsupportedMediaType)
		return
	}

	newUser := &users.NewUser{}
	err := json.NewDecoder(r.Body).Decode(newUser)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err = newUser.Validate()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := newUser.ToUser()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user, err = ctx.UserStore.Insert(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sessionToken, err := sessions.BeginSession(user.ID, ctx.Secret, ctx.SessionStore)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Authorization", "Bearer "+sessionToken)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}
