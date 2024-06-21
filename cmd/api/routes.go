package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/obimadu/ipc3-stage-2/internals/db"
	"github.com/obimadu/ipc3-stage-2/internals/handlers"
)

func router() *gin.Engine {
	// make router
	mux := gin.Default()

	// use cors
	mux.Use(cors.Default())

	// ROUTES
	// API group (v1)
	api := mux.Group("/api")

	// API/USERS group
	users := api.Group("/users")
	users.POST("/", handlers.CreateUser)

	users.GET("/", func(c *gin.Context) {
		handlers.GetAll(c, db.DB)
	})
	users.GET("/:userID", func(c *gin.Context) {
		handlers.GetUserByID(c, db.DB)
	})

	users.PUT("/", handlers.UpdateUser)
	users.PUT("/:userID", handlers.UpdateUser)

	users.DELETE("/", handlers.DeleteUser)
	users.DELETE("/:userID", handlers.DeleteUser)

	return mux
}
