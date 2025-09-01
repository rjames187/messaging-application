package users

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type MySQLStore struct {
	db *sql.DB
}

func NewMySQLStore(db *sql.DB) (MySQLStore, error) {
	if db == nil {
		return MySQLStore{}, fmt.Errorf("db must not be nil")
	}

	err := db.Ping()
	if err != nil {
		return MySQLStore{}, fmt.Errorf("error pinging db: %w", err)
	}

	return MySQLStore{db: db}, nil
}

func (s *MySQLStore) Insert(user *User) (*User, error) {
	insq := "INSERT INTO users(first_name, last_name, username, email, photo_url, pass_hash) VALUES(?,?,?,?,?,?)"
	res, err := s.db.Exec(insq, user.FirstName, user.LastName, user.Username, user.Email, user.PhotoURL, user.PassHash)
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	user.ID = int(id)
	return user, nil
}

func (s *MySQLStore) GetByID(id int) (*User, error) {
	gq := "SELECT id, first_name, last_name, username, email, photo_url, pass_hash FROM users where id = ?"
	rows, err := s.db.Query(gq, id)
	if err != nil {
		return nil, err
	}

	user := User{}
	found := rows.Next()
	if !found {
		return nil, errors.New("user was not found")
	}

	err = rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Username, &user.Email, &user.PhotoURL, &user.PassHash)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *MySQLStore) GetByEmail(email string) (*User, error) {
	gq := "SELECT id, first_name, last_name, username, email, photo_url, pass_hash FROM users where email = ?"
	rows, err := s.db.Query(gq, email)
	if err != nil {
		return nil, err
	}

	user := User{}
	found := rows.Next()
	if !found {
		return nil, errors.New("user was not found")
	}

	err = rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Username, &user.Email, &user.PhotoURL, &user.PassHash)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *MySQLStore) Update(id int, user *User) (*User, error) {
	uq := "UPDATE users SET first_name = ?, last_name = ?, username = ?, email = ?, photo_url = ?, pass_hash = ? WHERE id = ?"
	_, err := s.db.Exec(uq, user.FirstName, user.LastName, user.Username, user.Email, user.PhotoURL, user.PassHash, id)
	if err != nil {
		return nil, err
	}

	updatedUser, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	return updatedUser, nil
}
