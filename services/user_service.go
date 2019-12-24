package services

import (
	"github.com/agrism/bookstore_users-api/domain/users"
	errors "github.com/agrism/bookstore_users-api/utils"
	"github.com/agrism/bookstore_users-api/utils/date_utils"
)

func CreateUser(user users.User) (*users.User, *errors.RestErr) {
	if err := user.Validate(); err != nil {
		return nil, err
	}

	user.DateCreated = date_utils.GetNowString()

	if err := user.Save(); err != nil {
		return nil, err
	}

	return &user, nil
}

func GetUser(userId int64) (*users.User, *errors.RestErr) {

	result := &users.User{Id: userId}

	if err := result.Get(); err != nil {
		return nil, err
	}

	return result, nil
}

func UpdateUser(isPartial bool, user users.User) (*users.User, *errors.RestErr) {
	current, err := GetUser(user.Id)

	if err != nil {
		return nil, err
	}

	if err := user.Validate(); err != nil {
		return nil, err
	}

	if isPartial {
		if user.FirstName != "" {
			current.FirstName = user.FirstName
		}

		if user.LastName != "" {
			current.LastName = user.LastName
		}

		if user.Email != "" {
			current.Email = user.Email
		}

		if user.Status != "" {
			current.Status = user.Status
		}

		if user.Password != "" {
			current.Password = user.Password
		}
	} else {
		current.FirstName = user.FirstName
		current.LastName = user.LastName
		current.Email = user.Email
		current.Status = user.Status
		current.Password = user.Password
	}

	if err := current.Update(); err != nil {
		return nil, err
	}

	sanitized :=current.SanitizeUser()

	return &sanitized, nil
}

func DeleteUser(user users.User) *errors.RestErr {

	result := &users.User{Id: user.Id}

	if err := result.Delete(); err != nil {
		return err
	}

	return nil
}

func FindByStatus(status string) ([]users.User, *errors.RestErr) {
	dao := &users.User{}
	return dao.FindUserByStatus(status)
}
