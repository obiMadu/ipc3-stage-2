package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type jsonResponse struct {
	Status  string         `json:"status"`
	Message string         `json:"message"`
	Data    map[string]any `json:"data"`
	Error   map[string]any `json:"error,omitempty"`
}

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
	users.POST("/", createUser)

	users.GET("/", getUser)
	users.GET("/:userID", getUser)

	users.PUT("/", updateUser)
	users.PUT("/:userID", updateUser)

	users.DELETE("/", deleteUser)
	users.DELETE("/:userID", deleteUser)

	return mux
}
