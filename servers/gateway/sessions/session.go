package sessions

import (
	"crypto/rand"
	"encoding/base64"
)

func BeginSession(userID int, store Store) (string, string, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", "", err
	}
	sessionID, err := createSessionID(bytes)
	if err != nil {
		return "", "", err
	}
	err = store.Set(sessionID, userID)
	if err != nil {
		return "", "", nil
	}
	secret := base64.URLEncoding.EncodeToString(bytes)
	return sessionID, secret, nil
}