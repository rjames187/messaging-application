package users

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func generatePhotoURL(email string) string {
	cleaned := strings.TrimSpace(email)
	cleaned = strings.ToLower(cleaned)
	h := sha256.New()
	h.Write([]byte(cleaned))
	hash := hex.EncodeToString(h.Sum(nil))
	return fmt.Sprintf("https://gravatar.com/avatar/%s", hash)
}

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 13)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

type NewUser struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	Email     string `json:"email"`
}

func (nu *NewUser) Validate() error {
	if nu.Password == "" {
		return errors.New("password cannot be blank")
	}
	matches, _ := regexp.MatchString(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`, nu.Email)
	if !matches {
		return fmt.Errorf("invalid email address: %s", nu.Email)
	}
	matches, _ = regexp.MatchString(`^[^0-9]*$`, nu.FirstName)
	if !matches {
		return fmt.Errorf("invalid first name: %s", nu.FirstName)
	}
	matches, _ = regexp.MatchString(`^[^0-9]*$`, nu.LastName)
	if !matches {
		return fmt.Errorf("invalid last name: %s", nu.LastName)
	}
	return nil
}

func (nu *NewUser) ToUser() (*User, error) {
	u := User{
		FirstName: nu.FirstName,
		LastName:  nu.LastName,
		Username:  nu.Username,
		Email:     nu.Email,
	}

	u.PhotoURL = generatePhotoURL(nu.Email)

	PassHash, err := hashPassword(nu.Password)
	if err != nil {
		return &User{}, err
	}
	u.PassHash = string(PassHash)

	return &u, nil
}

type User struct {
	ID        int    `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Username  string `json:"username"`
	PassHash  string `json:"-"`
	Email     string `json:"-"`
	PhotoURL  string `json:"photoUrl"`
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

func (u *User) Authenticate(password string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(u.PassHash), []byte(password)); err != nil {
		return false
	}
	return true
}

type Updates struct {
	FirstName string
	LastName  string
	Password  string
	Email     string
}

func (u *User) ApplyUpdates(updates *Updates) error {
	if updates.FirstName != "" {
		u.FirstName = updates.FirstName
	}
	if updates.LastName != "" {
		u.LastName = updates.LastName
	}
	if updates.Password != "" {
		PassHash, err := hashPassword(updates.Password)
		if err != nil {
			return err
		}
		u.PassHash = string(PassHash)
	}
	if updates.Email != "" {
		u.PhotoURL = generatePhotoURL(updates.Email)
	}
	return nil
}

type Credentials struct {
	Email    string
	Password string
}
