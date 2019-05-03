package handlers

import (
	"github.com/AnthonyNixon/setsisaw/auth"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AuthCheck(c *gin.Context) {
	username, err := auth.GetUsernameFromAuthHeader(c)
	if err != nil {
		c.JSON(err.StatusCode(), gin.H{"error": err.Description()})
		return
	}

	c.String(http.StatusOK, username)
}
