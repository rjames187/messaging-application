package users

type Store interface {
	Insert(user *User) (*User, error)
	GetByID(id int) (*User, error)
	GetByEmail(email string) (*User, error)
	Update(id int, user *User) (*User, error)
}
