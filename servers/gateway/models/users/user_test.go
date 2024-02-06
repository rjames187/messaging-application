package models

import "testing"

func TestValidate(t *testing.T) {
	cases := []struct {
		input *NewUser
		output string
	}{
		{&NewUser{FirstName: "Bob", LastName: "Jones", Password: "password1", Email: "Bob.jones@gmail.com"}, "valid"},
		{&NewUser{Password: "password1", Email: "Bob.jones@gmail.com"}, "valid"},
		{&NewUser{Email: "Bob.jones@gmail.com"}, "invalid"},
		{&NewUser{Password: "password1", Email: "bobbyjones"}, "invalid"},
		{&NewUser{Password: "password1", Email: "bobbyjones@boby"}, "invalid"},
		{&NewUser{Password: "password1", Email: "bobby.com"}, "invalid"},
		{&NewUser{Password: "password1"}, "invalid"},
		{&NewUser{FirstName: "bob1", Password: "password1", Email: "bobyjones@bob.com"}, "invalid"},
		{&NewUser{LastName: "jones1", Password: "password1", Email: "bobyjones@bob.com"}, "invalid"},
	}

	for _, c := range cases {
		err := c.input.Validate()
		if (err == nil && c.output == "invalid") || (err != nil && c.output == "valid") {
			t.Errorf("incorrect output for `%s`: expected `%s`", c.input, c.output)
		}
	}
}