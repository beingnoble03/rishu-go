package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Ok(c *gin.Context) {
	message := "Hi! This is a Go Project."

	c.JSON(http.StatusOK, gin.H{
		"message": message,
	})
}
