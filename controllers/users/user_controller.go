package users

import (
	"fmt"
	"github.com/agrism/bookstore_users-api/domain/users"
	"github.com/agrism/bookstore_users-api/services"
	errors "github.com/agrism/bookstore_users-api/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func GetUser(context *gin.Context) {

	userId, userErr := strconv.ParseInt(context.Param("user_id"), 10, 54)

	if userErr != nil {
		restErr := errors.NewBadRequestError("Invalid user id, should be a number!")
		context.JSON(restErr.Status, restErr)
		return
	}

	user, getError := services.GetUser(userId)

	if getError != nil {
		context.JSON(getError.Status, getError)
		return
	}

	context.JSON(http.StatusOK, user)
}

func CreateUser(context *gin.Context) {
	var user users.User
	fmt.Println(user)

	if err := context.ShouldBindJSON(&user); err != nil {

		restError := errors.NewBadRequestError("invalid json body")
		context.JSON(restError.Status, restError)
		return
	}

	result, saveError := services.CreateUser(user)

	if saveError != nil {
		context.JSON(saveError.Status, saveError)
		return
	}

	context.JSON(http.StatusOK, result)
}

func SearchUser(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{
		"message": "findUser",
	})
}
