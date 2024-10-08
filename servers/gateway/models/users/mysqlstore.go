package users

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type MySQLStore struct {
	db *sql.DB
}

func (s *MySQLStore) Startup() error {
	dsn := fmt.Sprintf("root:%s@tcp(127.0.0.1:3306)/gateway", os.Getenv("MYSQL_ROOT_PASSWORD"))
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	if err = db.Ping(); err != nil {
		return err
	} else {
		fmt.Println("Successfully connected to the MySQL databases ...")
	}
	s.db = db
	return nil
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
	user.ID = id
	return user, nil
}

func (s *MySQLStore) Get(id int) (*User, error) {
	gq := "SELECT id,first_name,last_name,username,email,photo_url,pass_hash FROM users where id = ?"
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

func (s *MySQLStore) Update(id int, user *User) error {
	uq := "UPDATE users SET first_name = ?, last_name = ?, username = ?, email = ?, photo_url = ?, pass_hash = ? WHERE id = ?"
	_, err := s.db.Exec(uq, user.FirstName, user.LastName, user.Username, user.Email, user.PhotoURL, user.PassHash, id)
	if err != nil {
		return err
	}
	return nil
}