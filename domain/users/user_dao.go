package users

import (
	"fmt"
	"github.com/agrism/bookstore_users-api/datasources/mysql/users_db"
	"github.com/agrism/bookstore_users-api/logger"
	"github.com/agrism/bookstore_users-api/utils/errors"
)

const (
	queryInsertUser       = "INSERT INTO users(first_name, last_name, email, date_created, status, password) VALUES(?,?,?,?,?,?);"
	queryGetUser          = "SELECT id, first_name, last_name, email, date_created, status FROM users WHERE id = ?"
	queryUpdateUser       = "UPDATE users SET first_name=?, last_name=?, email=?, status=?, password=? WHERE id=?"
	queryDeleteUser       = "DELETE FROM users WHERE id=?"
	queryFindUserByStatus = "SELECT id, first_name, last_name, email, date_created, status FROM users WHERE status=?"
	StatusActive          = "active"
)

func (user *User) Get() *errors.RestErr {

	stmt, err := users_db.Client.Prepare(queryGetUser)

	if err != nil {
		logger.Error("error when trying to prepare get user statement", err)
		return errors.NewInternalServerError("database error")
	}

	defer stmt.Close()

	result := stmt.QueryRow(user.Id)

	if getError := result.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status);
		getError != nil {
		logger.Error("error when trying to to get user by id", getError)
		return errors.NewInternalServerError("database error")
	}

	return nil
}

func (user *User) Save() *errors.RestErr {

	stmt, err := users_db.Client.Prepare(queryInsertUser)

	if err != nil {
		logger.Error("error when trying to prepare save user statement", err)
		return errors.NewInternalServerError("database error")
	}

	defer stmt.Close()

	insertResult, saveErr := stmt.Exec(user.FirstName, user.LastName, user.Email, user.DateCreated, user.Status, user.Password)

	if saveErr != nil {
		logger.Error("error when trying to save user statement", saveErr)
		return errors.NewInternalServerError("database error")
	}

	userId, err := insertResult.LastInsertId()

	if err != nil {
		logger.Error("error when trying to get last inserted id after creating user", err)
		return errors.NewInternalServerError("database error")
	}

	user.Id = userId
	return nil
}

func (user *User) Update() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryUpdateUser)

	if err != nil {
		logger.Error("error when trying to prepare statement for update user", err)
		return errors.NewInternalServerError("database error")
	}

	defer stmt.Close()

	_, updErr := stmt.Exec(user.FirstName, user.LastName, user.Email, user.Status, user.Password, user.Id)

	if updErr != nil {
		logger.Error("error when trying to update user", updErr)
		return errors.NewInternalServerError("database error")
	}

	return nil
}

func (user *User) Delete() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryDeleteUser)

	if err != nil {
		logger.Error("error when trying to prepare delete for delete user", err)
		return errors.NewInternalServerError("database error")
	}

	defer stmt.Close()

	if _, delError := stmt.Exec(user.Id); delError != nil {
		logger.Error("error when trying to  delete user", delError)
		return errors.NewInternalServerError("database error")
	}

	return nil
}

func (user *User) FindUserByStatus(status string) ([]User, *errors.RestErr) {
	stmt, findErr := users_db.Client.Prepare(queryFindUserByStatus)
	if findErr != nil {
		logger.Error("error when trying to prepare find statement user", findErr)
		return nil, errors.NewInternalServerError("database error")
	}

	defer stmt.Close()

	rows, err := stmt.Query(status)

	if err != nil {
		logger.Error("error when trying to execute find statement", err)
		return nil, errors.NewInternalServerError("database error")
	}
	defer rows.Close()

	result := make([]User, 0)

	for rows.Next() {
		var user User
		if err := rows.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status)
			err != nil {
			logger.Error("error when trying to scan user in find by Status", err)
			return nil, errors.NewInternalServerError("database error")
		}
		result = append(result, user)
	}

	if len(result) == 0 {
		return nil, errors.NewNotFoundError(fmt.Sprintf("no users matchin status %s", status))
	}

	return result, nil
}
