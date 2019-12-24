package users

import (
	"github.com/agrism/bookstore_users-api/domain/users"
	"github.com/agrism/bookstore_users-api/services"
	"github.com/agrism/bookstore_users-api/utils/errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func getUserId(userIdParam string) (int64, *errors.RestErr) {
	userId, userErr := strconv.ParseInt(userIdParam, 10, 54)

	if userErr != nil {
		return 0, errors.NewBadRequestError("Invalid user id, should be a number!")
	}
	return userId, nil
}

func Get(ctx *gin.Context) {

	userId, idError := getUserId(ctx.Param("user_id"))
	if idError != nil {
		ctx.JSON(idError.Status, idError)
	}

	user, getError := services.GetUser(userId)

	if getError != nil {
		ctx.JSON(getError.Status, getError)
		return
	}

	ctx.JSON(http.StatusOK, user.Marshall(ctx.GetHeader("X-Public") == "true"))
}

func Create(ctx *gin.Context) {
	var user users.User

	if err := ctx.ShouldBindJSON(&user); err != nil {

		restError := errors.NewBadRequestError("invalid json body")
		ctx.JSON(restError.Status, restError)
		return
	}

	result, saveError := services.CreateUser(user)

	if saveError != nil {
		ctx.JSON(saveError.Status, saveError)
		return
	}

	ctx.JSON(http.StatusOK, result.Marshall(ctx.GetHeader("X-Public") == "true"))
}

func Update(ctx *gin.Context) {
	userId, idError := getUserId(ctx.Param("user_id"))
	if idError != nil {
		ctx.JSON(idError.Status, idError)
	}

	var user users.User

	if err := ctx.ShouldBindJSON(&user); err != nil {
		restError := errors.NewBadRequestError("invalid json body")
		ctx.JSON(restError.Status, restError)
		return
	}

	user.Id = userId

	isPartial := ctx.Request.Method == http.MethodPatch

	result, updateError := services.UpdateUser(isPartial, user)

	if updateError != nil {
		ctx.JSON(updateError.Status, updateError)
		return
	}

	ctx.JSON(http.StatusOK, result.Marshall(ctx.GetHeader("X-Public") == "true"))
}

func Delete(ctx *gin.Context) {

	userId, idError := getUserId(ctx.Param("user_id"))
	if idError != nil {
		ctx.JSON(idError.Status, idError)
	}

	var user users.User
	user.Id = userId

	if deleteError := services.DeleteUser(user); deleteError != nil {
		ctx.JSON(deleteError.Status, deleteError)
		return
	}

	ctx.JSON(http.StatusOK, map[string]string{"status": "deleted"})
}

func FindByStatus(ctx *gin.Context) {

	status := ctx.Query("status")

	users, err := services.FindByStatus(status)

	if err != nil {
		ctx.JSON(err.Status, err)
		return
	}

	ctx.JSON(http.StatusOK, users.Marshal(ctx.GetHeader("X-Public") == "true"))
}
