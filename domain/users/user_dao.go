package users

import (
	"github.com/agrism/bookstore_users-api/datasources/mysql/users_db"
	errors "github.com/agrism/bookstore_users-api/utils"
	"github.com/agrism/bookstore_users-api/utils/mysql_utils"
)

const (
	queryInsertUser = "INSERT INTO users(first_name, last_name, email, date_created) VALUES(?,?,?,?);"
	queryGetUser    = "SELECT id, first_name, last_name, email, date_created FROM users WHERE id = ?"
	queryUpdateUser = "UPDATE users SET first_name=?, last_name=?, email=? WHERE id=?"
	errorNoRows     = "no rows in result set"
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

	if scanError := result.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated);
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

	insertResult, saveErr := stmt.Exec(user.FirstName, user.LastName, user.Email, user.DateCreated)

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

	_, updErr := stmt.Exec(user.FirstName, user.LastName, user.Email, user.Id)

	if updErr != nil{
		return mysql_utils.ParseError(updErr)
	}

	return nil
}
