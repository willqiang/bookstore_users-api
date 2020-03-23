package app

import (
	"github.com/gin-gonic/gin"
	"github.com/willqiang/bookstore_users-api/logger"
)

var(
	router = gin.Default()
)

func StartApplication()  {
	mapUrls()

	logger.Info("Starting the application")
	router.Run(":8080")
}
