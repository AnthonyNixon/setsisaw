package handlers

import (
	"fmt"
	"github.com/AnthonyNixon/setsisaw/auth"
	"github.com/AnthonyNixon/setsisaw/types"
	"github.com/AnthonyNixon/setsisaw/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func NewLocation(c *gin.Context) {
	// Check Auth info
	claims, customErr := auth.GetUserInfo(c)
	if customErr != nil {
		c.JSON(customErr.StatusCode(), gin.H{"error": customErr.Description()})
		return
	}

	if !auth.IsEntitled(claims, "EDITOR") {
		c.JSON(http.StatusForbidden, gin.H{"Error": fmt.Sprintf("UserId %s is not entitled to add a location.", claims.Username)})
		return
	}

	// If we're here, the user is authorized to add an artist.

	var newLocation types.Location
	err := c.BindJSON(&newLocation)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Bad JSON Input, could not bind."})
		return
	}

	customErr = utils.NewLocation(newLocation)
	if customErr != nil {
		c.JSON(customErr.StatusCode(), customErr.Description())
		return
	}

	c.JSON(http.StatusCreated, newLocation)
}

func GetAllLocations(c *gin.Context) {
	// Check Auth info
	claims, customErr := auth.GetUserInfo(c)
	if customErr != nil {
		c.JSON(customErr.StatusCode(), gin.H{"error": customErr.Description()})
		return
	}

	if !auth.IsEntitled(claims, "USER") {
		c.JSON(http.StatusForbidden, gin.H{"Error": fmt.Sprintf("User %s is not entitled to get all locations.", claims.Username)})
		return
	}

	// If we're here, the user is authorized to get all artists.

	locations, customErr := utils.GetAllLocations()
	if customErr != nil {
		c.JSON(customErr.StatusCode(), customErr.Description())
		return
	}

	c.JSON(http.StatusOK, gin.H{"locations": locations, "count": len(locations)})

}

func GetLocation(c *gin.Context) {
	id := c.Param("id")

	// Check Auth info
	claims, customErr := auth.GetUserInfo(c)
	if customErr != nil {
		c.JSON(customErr.StatusCode(), gin.H{"error": customErr.Description()})
		return
	}

	if !auth.IsEntitled(claims, "USER") {
		c.JSON(http.StatusForbidden, gin.H{"Error": fmt.Sprintf("UserId %s is not entitled to get location info.", claims.Username)})
		return
	}

	// if we made it here, we're good to go.
	location, err := utils.GetLocation(id)
	if err != nil {
		c.JSON(err.StatusCode(), gin.H{"error": err.Description()})
		return
	}

	c.JSON(http.StatusOK, location)

}

func UpdateLocation(c *gin.Context) {
	id := c.Param("id")

	// Check Auth info
	claims, customErr := auth.GetUserInfo(c)
	if customErr != nil {
		c.JSON(customErr.StatusCode(), gin.H{"error": customErr.Description()})
		return
	}

	if !auth.IsEntitled(claims, "EDITOR") {
		c.JSON(http.StatusForbidden, gin.H{"Error": fmt.Sprintf("User %s is not entitled to edit location info.", claims.Username)})
		return
	}

	location, customErr := utils.GetLocation(id)
	if customErr != nil {
		c.JSON(customErr.StatusCode(), gin.H{"error": customErr.Description()})
		return
	}

	err := c.BindJSON(&location)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "could not bind location JSON", "details": err.Error()})
		return
	}

	// If we're here, the user is authorized to edit location information
	customErr = utils.UpdateLocation(id, location)
	if customErr != nil {
		c.JSON(customErr.StatusCode(), customErr.Description())
		return
	}

	c.JSON(http.StatusOK, location)
}
