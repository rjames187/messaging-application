package models

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type NewUser struct {
	FirstName string
	LastName  string
	Password  string
	Email     string
}

func (nu *NewUser) Validate() error {
	matches, _ := regexp.MatchString(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`, nu.Email)
	if !matches {
		return errors.New(fmt.Sprintf("Invalid email address: %s", nu.Email))
	}
	matches, _ = regexp.MatchString(`[^0-9]+`, nu.FirstName)
	if !matches {
		return errors.New(fmt.Sprintf("Invalid first name: %s", nu.FirstName))
	}
	matches, _ = regexp.MatchString(`[^0-9]+`, nu.LastName)
	if !matches {
		return errors.New(fmt.Sprintf("Invalid last name: %s", nu.LastName))
	}
	return nil
}

func (nu *NewUser) ToUser() (*User, error) {
	u := User{
		FirstName: nu.FirstName,
		LastName: nu.LastName,
		Email: nu.Email,
	}

	cleaned := strings.TrimSpace(nu.Email)
	cleaned = strings.ToLower(cleaned)
	h := sha256.New()
	h.Write([]byte(cleaned))
	hash := string(h.Sum(nil))
	u.PhotoURL = fmt.Sprintf("https://gravatar.com/avatar/%s", hash)

	PassHash, err := bcrypt.GenerateFromPassword([]byte(nu.Password), 13)
	if err != nil {
		return &User{}, err
	}
	u.PassHash = string(PassHash)

	return &u, nil
}

type User struct {
	FirstName string
	LastName  string
	PassHash  string
	Email     string
	PhotoURL  string
}

func (u *User) FullName() string {
	if u.FirstName != "" && u.LastName != "" {
		return fmt.Sprintf("%s %s", u.FirstName, u.LastName)
	} else if u.FirstName != "" {
		return u.FirstName
	} else if u.LastName != "" {
		return u.LastName
	} else {
		return u.Email
	}
}

type Credentials struct {
	Email string
	Password string
}

