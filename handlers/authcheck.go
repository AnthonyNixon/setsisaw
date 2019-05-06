package handlers

import (
	"github.com/AnthonyNixon/setsisaw/auth"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AuthCheck(c *gin.Context) {
	claims, err := auth.GetUserInfo(c)
	if err != nil {
		c.JSON(err.StatusCode(), gin.H{"error": err.Description()})
		return
	}

	c.JSON(http.StatusOK, claims)
}
