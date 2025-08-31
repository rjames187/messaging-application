package users

type Store interface {
	Insert(user *User) (*User, error)
	Get(id int) (*User, error)
	Update(id int, user *User) (*User, error)
}
