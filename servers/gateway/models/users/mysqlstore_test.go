package users

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

var u *User = &User{
	FirstName: "Bob",
	LastName:  "McDonald",
	Username:  "Bobby99",
	PassHash:  "rt36456346347",
	Email:     "bobby@gmail.com",
	PhotoURL:  "https://gravatar.com/avatar/a44ytya74yfya94y",
}

var uWithID *User = &User{
	ID:        1,
	FirstName: "Bob",
	LastName:  "McDonald",
	Username:  "Bobby99",
	PassHash:  "rt36456346347",
	Email:     "bobby@gmail.com",
	PhotoURL:  "https://gravatar.com/avatar/a44ytya74yfya94y",
}

func TestShouldInsertNewUser(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	mock.ExpectExec("INSERT INTO users").WithArgs(u.FirstName, u.LastName, u.Username, u.Email, u.PhotoURL, u.PassHash).WillReturnResult(sqlmock.NewResult(1, 1))

	store := MySQLStore{db: db}
	newUser, err := store.Insert(u)
	if err != nil {
		t.Errorf("Error inserting user: %s", err)
	}
	if newUser.ID != 1 {
		t.Errorf("Expected ID to be auto-assigned to 1 but got %d", newUser.ID)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Expectations were not met: %s", err)
	}
}

func TestShouldIncrementUserID(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	mock.ExpectExec("INSERT INTO users").WithArgs(u.FirstName, u.LastName, u.Username, u.Email, u.PhotoURL, u.PassHash).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO users").WithArgs(u.FirstName, u.LastName, u.Username, u.Email, u.PhotoURL, u.PassHash).WillReturnResult(sqlmock.NewResult(2, 1))

	store := MySQLStore{db: db}
	_, err := store.Insert(u)
	if err != nil {
		t.Errorf("Error inserting user: %s", err)
	}
	newUser, err := store.Insert(u)
	if err != nil {
		t.Errorf("Error inserting user: %s", err)
	}
	if newUser.ID != 2 {
		t.Errorf("Expected ID to be auto-assigned to 1 but got %d", newUser.ID)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Expectations were not met: %s", err)
	}
}

func TestShouldReturnInsertionError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	mock.ExpectExec("INSERT INTO users").WithArgs(u.FirstName, u.LastName, u.Username, u.Email, u.PhotoURL, u.PassHash).WillReturnError(errors.New(""))

	store := MySQLStore{db: db}
	_, err := store.Insert(u)
	if err == nil {
		t.Errorf("Expected insertion error")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Expectations were not met: %s", err)
	}
}

func TestShouldSelectUserByID(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	data := sqlmock.NewRows([]string{"id", "first_name", "last_name", "username", "email", "photo_url", "pass_hash"})
	data.AddRow(uWithID.ID, u.FirstName, u.LastName, u.Username, u.Email, u.PhotoURL, u.PassHash)

	mock.ExpectQuery("SELECT").WithArgs(uWithID.ID).WillReturnRows(data)

	store := MySQLStore{db: db}
	newUser, err := store.Get(1)
	if err != nil {
		t.Errorf("Error fetching user from database: %s", err)
	}
	if newUser.ID != 1 {
		t.Errorf("Expected returned ID to be 1 but got %d", newUser.ID)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Expectations were not met: %s", err)
	}
}

func TestShouldHandleMissingUser(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	data := sqlmock.NewRows([]string{"id", "first_name", "last_name", "username", "email", "photo_url", "pass_hash"})
	mock.ExpectQuery("SELECT").WithArgs(uWithID.ID).WillReturnRows(data)

	store := MySQLStore{db: db}
	_, err := store.Get(1)
	if err == nil {
		t.Errorf("Expected Get operation to return a not found error")
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Expectations were not met: %s", err)
	}
}

func TestShouldUpdateUser(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	mock.ExpectExec("UPDATE").WithArgs(u.FirstName, u.LastName, u.Username, u.Email, u.PhotoURL, u.PassHash, uWithID.ID).WillReturnResult(sqlmock.NewResult(1, 1))

	store := MySQLStore{db: db}
	updatedUser, err := store.Update(1, u)
	if err != nil {
		t.Errorf("Error updating user: %s", err)
	}
	if updatedUser.ID != 1 {
		t.Errorf("Expected updated ID to be 1 but got %d", updatedUser.ID)
	}
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Expectations were not met: %s", err)
	}
}
