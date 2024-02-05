package models

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