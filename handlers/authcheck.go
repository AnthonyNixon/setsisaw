package handlers

import (
	"github.com/AnthonyNixon/setsisaw/auth"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AuthCheck(c *gin.Context) {
	username, err := auth.GetUsernameFromAuthHeader(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.String(http.StatusOK, username)
}
