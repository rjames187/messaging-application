package users

import "fmt"

type StubStore struct {
	users  map[int]*User
	serial int
}

func NewStubStore() *StubStore {
	return &StubStore{
		users:  make(map[int]*User),
		serial: 1,
	}
}

func (s *StubStore) Insert(user *User) (*User, error) {
	_, err := s.GetByEmail(user.Email)
	if err == nil {
		return nil, fmt.Errorf("user with email %s already exists", user.Email)
	}

	_, err = s.GetByUsername(user.Username)
	if err == nil {
		return nil, fmt.Errorf("user with username %s already exists", user.Username)
	}

	user.ID = s.serial
	s.serial++
	s.users[user.ID] = user
	return user, nil
}

func (s *StubStore) GetByID(id int) (*User, error) {
	user, exists := s.users[id]
	if !exists {
		return nil, fmt.Errorf("user was not found")
	}
	return user, nil
}

func (s *StubStore) GetByEmail(email string) (*User, error) {
	for _, user := range s.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, fmt.Errorf("user was not found")
}

func (s *StubStore) GetByUsername(username string) (*User, error) {
	for _, user := range s.users {
		if user.Username == username {
			return user, nil
		}
	}
	return nil, fmt.Errorf("user was not found")
}

func (s *StubStore) Update(id int, user *User) (*User, error) {
	existingUser, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	if existingUser.Email != user.Email {
		_, err := s.GetByEmail(user.Email)
		if err == nil {
			return nil, fmt.Errorf("user with email %s already exists", user.Email)
		}
	}

	if existingUser.Username != user.Username {
		_, err := s.GetByUsername(user.Username)
		if err == nil {
			return nil, fmt.Errorf("user with username %s already exists", user.Username)
		}
	}

	// Update the existing user's fields
	existingUser.FirstName = user.FirstName
	existingUser.LastName = user.LastName
	existingUser.Email = user.Email
	existingUser.PassHash = user.PassHash

	return existingUser, nil
}
