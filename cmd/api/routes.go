package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func router() *gin.Engine {
	// make router
	mux := gin.Default()

	// use cors
	mux.Use(cors.Default())

	// routes
	// routes.users
	mux.POST("/users", createUser)

	mux.GET("/users", getUser)
	mux.GET("/users/:userID", getUser)

	mux.PUT("/users", updateUser)
	mux.PUT("users/:userID", updateUser)

	mux.DELETE("users", deleteUser)
	mux.DELETE("users/:userID", deleteUser)

	return mux
}
