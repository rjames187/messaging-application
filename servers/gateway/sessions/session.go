package sessions

import (
	"crypto/rand"
	"encoding/base64"
)

const SESSIONID_LENGTH = 32

func BeginSession(userID int, store Store) (string, string, error) {
	secretBytes := make([]byte, 32)
	_, err := rand.Read(secretBytes); if err != nil {
		return "", "", err
	}
	secret := base64.URLEncoding.EncodeToString(secretBytes)
	sessionToken, sessionID, err := createSessionToken(secret, SESSIONID_LENGTH)
	if err != nil {
		return "", "", err
	}
	err = store.Set(sessionID, userID)
	if err != nil {
		return "", "", nil
	}
	return sessionToken, secret, nil
}