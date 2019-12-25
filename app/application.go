package app

import (
	"github.com/agrism/bookstore_users-api/logger"
	"github.com/gin-gonic/gin"
)

var (
	router = gin.Default()
)

func StartApplication() {
	mapUrls()

	logger.GetLogger().Info("about to start application...")
	router.Run("localhost:8080")
}
