package handlers

import (
	"github.com/AnthonyNixon/setsisaw/auth"
	"github.com/gin-gonic/gin"
	"net/http"
)

func RefreshToken(c *gin.Context) {

	tokenString, err := auth.RefreshToken(c)
	if err != nil {
		c.JSON(err.StatusCode(), gin.H{"error": err.Description()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}
