package app

import (
	"github.com/willqiang/bookstore_users-api/controllers/ping"
	"github.com/willqiang/bookstore_users-api/controllers/users"
)

func mapUrls()  {
	router.GET("/ping", ping.Ping)

	router.POST("/users", users.Creat)
	router.GET("/users/:user_id", users.Get)
	router.PUT("/users/:user_id", users.Update)
	router.PATCH("/users/:user_id", users.Update)
	router.DELETE("/users/:user_id", users.Delete)
	router.GET("/internal/users/search", users.Search)
}
