package users

import (
	errors2 "github.com/agrism/bookstore_users-api/utils/errors"
	"strings"
)

type User struct {
	Id          int64  `json:"id`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	DateCreated string `json:"date_created"`
	Status      string `json:"status"`
	Password    string `json:"password"`
}

type  Users []User

func (user *User) Validate() *errors2.RestErr {

	user.FirstName = strings.TrimSpace(user.FirstName)
	user.LastName = strings.TrimSpace(user.LastName)

	user.Email = strings.TrimSpace(strings.ToLower(user.Email))

	if user.Email == "" {
		return errors2.NewBadRequestError("invalid email address")
	}

	user.Password = strings.TrimSpace(user.Password)
	if user.Password == "" {
		return errors2.NewBadRequestError("invalid password")
	}

	return nil
}

func (user *User) SanitizeUser() User {
	return User{user.Id,
		user.FirstName,
		user.LastName,
		user.Email,
		user.DateCreated,
		user.Status,
		"",
	}
}
