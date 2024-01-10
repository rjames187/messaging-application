package sessions

const SESSIONID_LENGTH = 32

func BeginSession(userID int, secret string, store Store) (string, error) {
	sessionToken, sessionID, err := createSessionToken(secret, SESSIONID_LENGTH)
	if err != nil {
		return "", err
	}
	err = store.Set(sessionID, userID)
	if err != nil {
		return "", nil
	}
	return sessionToken, nil
}

func GetSessionState(sessionToken string) (int, error) {

}