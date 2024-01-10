package sessions

const SESSIONID_LENGTH = 32

func BeginSession(userID int, store Store) (string, string, error) {
	sessionToken, sessionID, secret, err := createSessionToken(SESSIONID_LENGTH)
	if err != nil {
		return "", "", err
	}
	err = store.Set(sessionID, userID)
	if err != nil {
		return "", "", nil
	}
	return sessionToken, secret, nil
}

func GetSessionState(sessionToken string) (int, error) {

}