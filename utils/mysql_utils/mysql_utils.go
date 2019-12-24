package mysql_utils

import (
	errors2 "github.com/agrism/bookstore_users-api/utils/errors"
	"github.com/go-sql-driver/mysql"
	"strings"
)

const (
	errorNoRows = "no rows in result set"
)

func ParseError(err error) *errors2.RestErr {
	sqlError, ok := err.(*mysql.MySQLError)
	if !ok {
		if strings.Contains(err.Error(), errorNoRows) {
			return errors2.NewNotFoundError("no record matching given id")
		}
		return errors2.NewInternalServerError("error parsing db response")
	}

	switch sqlError.Number {
	case 1062:
		return errors2.NewBadRequestError("invalid data")
	}

	return errors2.NewInternalServerError("error processing request")
}
