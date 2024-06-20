package main

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/obimadu/ipc3-stage-2/internals/db"
	"github.com/obimadu/ipc3-stage-2/internals/models"
)

func getUser(c *gin.Context) {
	var idStr string
	username, exists := c.Get("username")

	if !exists && c.Param("userID") == "" {
		users, err := models.GetAll(db.DB)
		if err != nil {
			c.JSON(http.StatusInternalServerError, jsonResponse{
				Status:  "error",
				Message: "Failed to retrieve users.",
			})
			return
		}

		c.JSON(http.StatusOK, jsonResponse{
			Status:  "success",
			Message: "Retrieved all users.",
			Data: gin.H{
				"users": users,
			},
		})
		return
	} else if !exists && c.Param("userID") != "" {
		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, jsonResponse{
				Status:  "error",
				Message: "UserID must be a positive interger.",
			})
			return
		}

		user, err := models.GetUserByID(db.DB, uint(id))
		if err != nil {
			c.JSON(http.StatusInternalServerError, jsonResponse{
				Status:  "error",
				Message: "Failed to retrieve user.",
			})
			return
		}

		c.JSON(http.StatusOK, jsonResponse{
			Status:  "success",
			Message: "User retrieved succesfully.",
			Data: gin.H{
				"user": user,
			},
		})
		return
	}

	user, err := models.GetUserByUsername(db.DB, username.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, jsonResponse{
			Status:  "error",
			Message: "Failed to retrieve user.",
		})
		return
	}

	c.JSON(http.StatusOK, jsonResponse{
		Status:  "success",
		Message: "User retrieved succesfully.",
		Data: gin.H{
			"user": user,
		},
	})
}

func createUser(c *gin.Context) {
	var user models.Users

	err := c.ShouldBindBodyWithJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, jsonResponse{
			Status:  "error",
			Message: "Request body not valid.",
			Error: gin.H{
				"error": err.Error(),
			},
		})
		return
	}

	if user.Username == "" || user.Email == "" {
		c.JSON(http.StatusBadRequest, jsonResponse{
			Status:  "error",
			Message: "You must specify both a username and an email.",
		})
		return
	}

	err = models.CreateUser(db.DB, user)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") && strings.Contains(err.Error(), "email") {
			c.JSON(http.StatusBadRequest, jsonResponse{
				Status:  "error",
				Message: "Email has been taken!",
			})
			return
		} else if strings.Contains(err.Error(), "duplicate key value violates unique constraint") && strings.Contains(err.Error(), "username") {
			c.JSON(http.StatusBadRequest, jsonResponse{
				Status:  "error",
				Message: "Username has been taken!",
			})
			return
		} else {
			c.JSON(http.StatusInternalServerError, jsonResponse{
				Status:  "error",
				Message: "Failed to create user.",
				Error: gin.H{
					"error": err.Error(),
				},
			})
			return
		}
	}

	c.JSON(http.StatusOK, jsonResponse{
		Status:  "success",
		Message: "User created succesfully.",
	})
}

func updateUser(c *gin.Context) {

	var user models.Users
	err := c.ShouldBindBodyWithJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, jsonResponse{
			Status:  "error",
			Message: "Request body not valid.",
		})
		return
	}

	username, exists := c.Get("username")
	if !exists && c.Param("userID") != "" {
		idStr := c.Param("userID")
		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, jsonResponse{
				Status:  "error",
				Message: "UserID must be a positive interger.",
			})
			return
		}

		err = models.UpdateUserByID(db.DB, uint(id), user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, jsonResponse{
				Status:  "error",
				Message: "Failed to update user.",
			})
			return
		}

		c.JSON(http.StatusOK, jsonResponse{
			Status:  "success",
			Message: "User updated succesfully.",
		})
		return
	}

	if exists && c.Param("userID") == "" {
		err = models.UpdateUserByUsername(db.DB, username.(string), user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, jsonResponse{
				Status:  "error",
				Message: "Failed to update user.",
			})
			return
		}

		c.JSON(http.StatusOK, jsonResponse{
			Status:  "success",
			Message: "User updated succesfully.",
		})
		return
	}

	c.JSON(http.StatusBadRequest, jsonResponse{
		Status:  "error",
		Message: "You must/can only specify a user to update.",
	})

}

func deleteUser(c *gin.Context) {
	username, exists := c.Get("username")
	if !exists && c.Param("userID") != "" {
		idStr := c.Param("userID")
		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, jsonResponse{
				Status:  "error",
				Message: "UserID must be a positive interger.",
			})
			return
		}

		err = models.DeleteUserByID(db.DB, uint(id))
		if err != nil {
			c.JSON(http.StatusInternalServerError, jsonResponse{
				Status:  "error",
				Message: "Failed to delete user.",
			})
			return
		}

		c.JSON(http.StatusOK, jsonResponse{
			Status:  "success",
			Message: "User deleted succesfully.",
		})
		return
	}

	if exists && c.Param("userID") == "" {
		err := models.DeleteUserByUsername(db.DB, username.(string))
		if err != nil {
			c.JSON(http.StatusInternalServerError, jsonResponse{
				Status:  "error",
				Message: "Failed to delete user.",
			})
			return
		}

		c.JSON(http.StatusOK, jsonResponse{
			Status:  "success",
			Message: "User deleted succesfully.",
		})
		return
	}

	c.JSON(http.StatusBadRequest, jsonResponse{
		Status:  "error",
		Message: "You must/can only specify a user to delete.",
	})
}
