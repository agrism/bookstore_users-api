package users

import (
	"fmt"
	"github.com/agrism/bookstore_users-api/datasources/mysql/users_db"
	errors "github.com/agrism/bookstore_users-api/utils"
	"github.com/agrism/bookstore_users-api/utils/mysql_utils"
)

const (
	queryInsertUser       = "INSERT INTO users(first_name, last_name, email, date_created, status, password) VALUES(?,?,?,?,?,?);"
	queryGetUser          = "SELECT id, first_name, last_name, email, date_created, status FROM users WHERE id = ?"
	queryUpdateUser       = "UPDATE users SET first_name=?, last_name=?, email=?, status=?, password=? WHERE id=?"
	queryDeleteUser       = "DELETE FROM users WHERE id=?"
	queryFindUserByStatus = "SELECT id, first_name, last_name, email, date_created, status FROM users WHERE status=?"
)

var (
	userDB = make(map[int64]*User)
)

func (user *User) Get() *errors.RestErr {

	stmt, err := users_db.Client.Prepare(queryGetUser)

	if err != nil {
		return errors.NewBadRequestError(err.Error())
	}

	defer stmt.Close()

	result := stmt.QueryRow(user.Id)

	if scanError := result.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status);
		scanError != nil {
		return mysql_utils.ParseError(scanError)
	}

	return nil
}

func (user *User) Save() *errors.RestErr {

	stmt, err := users_db.Client.Prepare(queryInsertUser)

	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}

	defer stmt.Close()

	insertResult, saveErr := stmt.Exec(user.FirstName, user.LastName, user.Email, user.DateCreated, user.Status, user.Password)

	if saveErr != nil {
		return mysql_utils.ParseError(saveErr)
	}

	userId, err := insertResult.LastInsertId()

	if err != nil {
		return mysql_utils.ParseError(err)
	}

	user.Id = userId
	return nil
}

func (user *User) Update() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryUpdateUser)

	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}

	defer stmt.Close()

	_, updErr := stmt.Exec(user.FirstName, user.LastName, user.Email, user.Status, user.Password, user.Id)

	if updErr != nil {
		return mysql_utils.ParseError(updErr)
	}

	return nil
}

func (user *User) Delete() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryDeleteUser)

	if err != nil {
		return errors.NewBadRequestError(err.Error())
	}

	defer stmt.Close()

	if _, delError := stmt.Exec(user.Id); delError != nil {
		return mysql_utils.ParseError(delError)
	}

	return nil
}

func (user *User) FindUserByStatus(status string) ([]User, *errors.RestErr) {
	stmt, findErr := users_db.Client.Prepare(queryFindUserByStatus)
	if findErr != nil {
		return nil, errors.NewInternalServerError(findErr.Error())
	}

	defer stmt.Close()

	rows, err := stmt.Query(status)

	if err != nil {
		return nil, mysql_utils.ParseError(err)
	}
	defer rows.Close()

	result := make([]User, 0)

	for rows.Next() {
		var user User
		if err := rows.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status)
			err != nil {
			return nil, mysql_utils.ParseError(err)
		}
		result = append(result, user)
	}

	if len(result) == 0 {
		return nil, errors.NewNotFoundError(fmt.Sprintf("no users matchin status %s", status))
	}

	return result, nil
}
