package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/obimadu/ipc3-stage-2/internals/handlers"
)

func router() *gin.Engine {
	// make router
	mux := gin.Default()

	// use cors
	mux.Use(cors.Default())

	// ROUTES
	// API group (v1)
	api := mux.Group("/api/v1")

	// API/USERS group
	users := api.Group("/users")
	users.POST("/", handlers.CreateUser)

	users.GET("/", handlers.GetUser)
	users.GET("/:userID", handlers.GetUser)

	users.PUT("/", handlers.UpdateUser)
	users.PUT("/:userID", handlers.UpdateUser)

	users.DELETE("/", handlers.DeleteUser)
	users.DELETE("/:userID", handlers.DeleteUser)

	return mux
}
