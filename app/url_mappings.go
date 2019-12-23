package app

import (
	"github.com/agrism/bookstore_users-api/controllers/ping"
	"github.com/agrism/bookstore_users-api/controllers/users"
)

func mapUrls() {
	router.GET("/ping", ping.Ping)
	//router.GET("/users/search", controllers.SearchUser)
	router.GET("/users/:user_id", users.GetUser)
	router.POST("/users", users.CreateUser)
}
