package sessions

import "errors"

const SESSIONID_LENGTH = 32

func BeginSession(sessionState string, secret string, store Store) (string, error) {
	sessionToken, sessionID, err := createSessionToken(secret, SESSIONID_LENGTH)
	if err != nil {
		return "", err
	}
	err = store.Set(sessionID, sessionState)
	if err != nil {
		return "", err
	}
	return sessionToken, nil
}

func GetSessionState(sessionToken string, secret string, store Store) (string, error) {
	valid, sessionID, err := validToken(sessionToken, secret, SESSIONID_LENGTH)
	if err != nil {
		return "", err
	}
	if !valid {
		return "", errors.New("invalid session token")
	}
	userID, err := store.Get(sessionID)
	if err != nil {
		return "", err
	}
	return userID, nil
}

func EndSession(sessionToken string, store Store) error {
	sessionID, err := extractIDFromToken(sessionToken, SESSIONID_LENGTH)
	if err != nil {
		return err
	}
	err = store.Delete(sessionID)
	if err != nil {
		return err
	}
	return nil
}
