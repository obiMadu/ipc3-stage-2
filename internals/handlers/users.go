package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/obimadu/ipc3-stage-2/internals/db"
	"github.com/obimadu/ipc3-stage-2/internals/models"
)

type jsonResponse struct {
	Status  string         `json:"status"`
	Message string         `json:"message"`
	Data    map[string]any `json:"data"`
	Error   map[string]any `json:"error,omitempty"`
}

func CreateUser(c *gin.Context) {
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
		checkUnique(c, err)
		return
	}

	c.JSON(http.StatusOK, jsonResponse{
		Status:  "success",
		Message: "User created succesfully.",
	})
}

func UpdateUser(c *gin.Context) {

	var user models.Users
	err := c.ShouldBindBodyWithJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, jsonResponse{
			Status:  "error",
			Message: "Request body not valid.",
		})
		return
	}

	username := c.Query("username")
	userID := c.Param("userID")
	if username == "" && userID != "" {
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
			checkUnique(c, err)
			return
		}

		c.JSON(http.StatusOK, jsonResponse{
			Status:  "success",
			Message: "User updated succesfully.",
		})
		return
	} else if userID == "" && username != "" {
		err = models.UpdateUserByUsername(db.DB, username, user)
		if err != nil {
			checkUnique(c, err)
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

func DeleteUser(c *gin.Context) {
	username := c.Query("username")
	if username == "" && c.Param("userID") != "" {
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
			checkRecordExists(c, err)
			return
		}

		c.JSON(http.StatusOK, jsonResponse{
			Status:  "success",
			Message: "User deleted succesfully.",
		})
		return
	}

	if username != "" && c.Param("userID") == "" {
		err := models.DeleteUserByUsername(db.DB, username)
		if err != nil {
			checkRecordExists(c, err)
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

func checkUnique(c *gin.Context, err error) {
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
			Message: "User operation failed",
			Error: gin.H{
				"error": err.Error(),
			},
		})
		return
	}
}

func checkRecordExists(c *gin.Context, err error) {
	if strings.Contains(err.Error(), "record not found") {
		c.JSON(http.StatusBadRequest, jsonResponse{
			Status:  "error",
			Message: "User does not exist",
		})
		return
	} else {
		c.JSON(http.StatusInternalServerError, jsonResponse{
			Status:  "error",
			Message: "Failed to retrieve user.",
			Error: gin.H{
				"error": err.Error(),
			},
		})
		return
	}
}
