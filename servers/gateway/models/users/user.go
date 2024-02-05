package models

import (
	"errors"
	"fmt"
	"regexp"
)

type NewUser struct {
	FirstName string
	LastName  string
	Password  string
	Email     string
	Username  string
}

type User struct {
	FirstName string
	LastName  string
	PassHash  string
	Email     string
	Username  string
	PhotoURL  string
}

type Credentials struct {
	Username string
	Password string
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