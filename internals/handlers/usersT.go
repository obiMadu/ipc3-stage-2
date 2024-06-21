package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/obimadu/ipc3-stage-2/internals/models"
	"gorm.io/gorm"
)

func GetAll(c *gin.Context, db *gorm.DB) {

	if username := c.Query("username"); username != "" {
		user, err := models.GetUserByUsername(db, username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, jsonResponse{
				Status:  "error",
				Message: "Failed to retrieve users.",
				Error: gin.H{
					"error": err.Error(),
				}})
			return
		}

		c.JSON(http.StatusOK, jsonResponse{
			Status:  "success",
			Message: "User retrieved successfully.",
			Data: gin.H{
				"user": user,
			},
		})
		return
	}

	users, err := models.GetAll(db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, jsonResponse{
			Status:  "error",
			Message: "Failed to retrieve users.",
			Error: gin.H{
				"error": err.Error(),
			}})
		return
	}

	c.JSON(http.StatusOK, jsonResponse{
		Status:  "success",
		Message: "Retrieved all users.",
		Data: gin.H{
			"users": users,
		},
	})
}

func GetUserByID(c *gin.Context, db *gorm.DB) {
	userID := c.Param("userID")
	id, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, jsonResponse{
			Status:  "error",
			Message: "UserID must be a positive interger.",
		})
		return
	}

	user, err := models.GetUserByID(db, uint(id))
	if err != nil {
		checkRecordExists(c, err)
		return
	}

	c.JSON(http.StatusOK, jsonResponse{
		Status:  "success",
		Message: "User retrieved successfully.",
		Data: gin.H{
			"user": user,
		},
	})
}
