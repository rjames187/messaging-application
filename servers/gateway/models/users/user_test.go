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

func TestToUser(t *testing.T) {
	cases := []struct {
		input *NewUser
		output *User
	}{
		{
			&NewUser{FirstName: "Bob", LastName: "Jones", Password: "password1", Email: "myemailaddress@example.com"},
			&User{FirstName: "Bob", LastName: "Jones", PhotoURL: "https://gravatar.com/avatar/84059b07d4be67b806386c0aad8070a23f18836bbaae342275dc0a83414c32ee", Email: "myemailaddress@example.com"},
		},
		{
			&NewUser{Password: "password1", Email: " MyemailAddress@example.com  "},
			&User{PhotoURL: "https://gravatar.com/avatar/84059b07d4be67b806386c0aad8070a23f18836bbaae342275dc0a83414c32ee", Email: " MyemailAddress@example.com  "},
		},
	}

	for _, c := range cases {
		u, err := c.input.ToUser()
		if err != nil {
			t.Errorf("Error converting NewUser `%s` to User: %s", c.input, err)
		}
		if u.FirstName != c.output.FirstName {
			t.Errorf("Expected FirstName `%s` but got `%s`", c.output.FirstName, u.FirstName)
		}
		if u.LastName != c.output.LastName {
			t.Errorf("Expected LastName `%s` but got `%s`", c.output.LastName, u.LastName)
		}
		if u.Email != c.output.Email {
			t.Errorf("Expected Email `%s` but got `%s`", c.output.Email, u.Email)
		}
		if u.PhotoURL != c.output.PhotoURL {
			t.Errorf("Expected PhotoURL `%s` but got `%s`", c.output.PhotoURL, u.PhotoURL)
		}
	}
}

func TestFullName(t *testing.T) {
	cases := []struct{
		input *User
		output string
	}{
		{
			&User{FirstName: "Jimmy", LastName: "John", Email: "jimmyjohn3@gmail.com"},
			"Jimmy John",
		},
		{
			&User{FirstName: "Jimmy", Email: "jimmyjohn3@gmail.com"},
			"Jimmy",
		},
		{
			&User{LastName: "John", Email: "jimmyjohn3@gmail.com"},
			"John",
		},
		{
			&User{Email: "jimmyjohn3@gmail.com"},
			"jimmyjohn3@gmail.com",
		},
	}

	for _, c := range cases {
		name := c.input.FullName()
		if name != c.output {
			t.Errorf("Expected %s but got %s", c.output, name)
		}
	}
}

func TestAuthenticate(t *testing.T) {
	nu := NewUser{
		Password: "Epic1245!",
	}
	u, _ := nu.ToUser()

	cases := []struct{
		user *User
		password string
		output bool
	}{
		{u, "Epic1245!", true},
		{u, "fdgsdsg45", false},
		{u, "", false},
	}

	for _, c := range cases {
		authenticated := c.user.Authenticate(c.password)
		if authenticated != c.output {
			t.Errorf("Expected %t but got %t", c.output, authenticated)
		}
	}
}

func TestApplyUpdates(t *testing.T) {
	userFixture, _  := (&NewUser{Password: "boss88", Email: "funny@gmail.com"}).ToUser()

	cases := []struct{
		user *User
		updates *Updates
	}{
		{&User{}, &Updates{FirstName: "Bob"}},
		{&User{FirstName: "Donald"}, &Updates{FirstName: "Bob"}},
		{&User{LastName: "Mac"}, &Updates{LastName: "Thomas"}},
		{userFixture, &Updates{Password: "dfasf5ya5e6"}},
		{userFixture, &Updates{Password: "afsdg44", Email: "cool@hotmail.com"}},
	}

	for _, c := range cases {
		oldPhoto := c.user.PhotoURL
		err := c.user.ApplyUpdates(c.updates)
		if err != nil {
			t.Errorf("Error applying updates: %s", err)
		}
		if c.updates.FirstName != "" && c.user.FirstName != c.updates.FirstName {
			t.Errorf("Expected FirstName to be %s but got %s", c.updates.FirstName, c.user.FirstName)
		}
		if c.updates.LastName != "" && c.user.LastName != c.updates.LastName {
			t.Errorf("Expected LastName to be %s but got %s", c.updates.LastName, c.user.LastName)
		}
		if c.updates.Email != "" && c.user.PhotoURL == oldPhoto {
			t.Errorf("Expected photo URL to change because of email address change")
		}
		if c.updates.Password != "" && !c.user.Authenticate(c.updates.Password) {
			t.Errorf("New password %s failed to authenticate", c.updates.Password)
		}
	}
}