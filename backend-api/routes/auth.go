package routes

import (
	"net/http"

	"livecode-api/database"
	"livecode-api/handlers"
	"livecode-api/models"

	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	var req models.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.RegisterResponse{
			Success: false,
			Message: "An unexpected error occurred. Please try again.",
		})
		return
	}

	response, err := handlers.RegisterUserInternal(req, database.DB)

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.RegisterResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	statusCode := http.StatusCreated
	if !response.Success {
		statusCode = http.StatusOK
	}

	c.JSON(statusCode, response)
}

func CheckFieldAvailable(c *gin.Context) {
	field := c.Query("field")
	value := c.Query("value")

	if field == "" || value == "" {
		c.JSON(http.StatusBadRequest, gin.H{"available": nil})
		return
	}

	available, err := handlers.CheckFieldAvailableInternal(field, value, database.DB)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"available": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"available": available})
}
