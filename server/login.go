package server

import "fmt"

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	valid    bool
}

func handleLogin(input string, user *User, users []*User) string {

	cmd, args, err := parseCommand(input)
	if err != nil {
		return fmt.Sprintf("%d %s", StatusSyntaxError, "Syntax error.")
	}

	switch {
	case cmd == "USER" && args == "":
		return fmt.Sprintf("%d %s", StatusSyntaxError, "Syntax error.")
	case cmd == "USER" && args != "":
		user.Username = args
		return fmt.Sprintf("%d %s", StatusUserOk, "User name okay, need password.")
	case cmd == "PASS" && user.Username != "" && args != "":
		user.Password = args
	}

	user.Auth(users)

	if !user.valid {
		return fmt.Sprintf("%d %s", StatusNotLoggedIn, "Not logged in.")
	}
	return fmt.Sprintf("%d %s", StatusUserLoggedIn, "User logged in, proceed.")
}

func (user *User) Auth(users []*User) {
	for _, u := range users {
		if user.Username == u.Username && user.Password == u.Password {
			user.valid = true
		}
	}
}
