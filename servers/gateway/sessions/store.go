package sessions

type Store interface {
	Get(sessionID string) (int, error)
	Set(sessionID string, userID int) error
	Update(sessionID string, newUserID int) error
}