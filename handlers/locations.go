package handlers

import (
	"fmt"
	"github.com/AnthonyNixon/setsisaw/auth"
	"github.com/AnthonyNixon/setsisaw/database"
	"github.com/AnthonyNixon/setsisaw/types"
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

	// TODO: This shouldn't be an editor only function. Anyone should be able to add a location. This is just for testing claims.
	if !auth.IsEntitled(claims, "EDITOR") {
		c.JSON(http.StatusForbidden, gin.H{"Error": fmt.Sprintf("User %s is not entitled to add a location.", claims.Username)})
		return
	}

	// If we're here, the user is authorized to add an artist.

	var newLocation types.Location
	err := c.BindJSON(&newLocation)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Bad JSON Input, could not bind."})
		return
	}

	unique, err := isNewLocationUnique(newLocation)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !unique {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Location already created"})
		return
	}


	db, err := database.GetConnection()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer db.Close()

	stmt, err := db.Prepare("insert into locations (name, description, city, state, country) values(?,?,?,?,?);")
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	_, err = stmt.Exec(newLocation.Name, newLocation.Description, newLocation.City, newLocation.State, newLocation.Country)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, gin.H{"name": newLocation.Name, "description": newLocation.Description})
}

func isNewLocationUnique(newLocation types.Location) (bool, error) {
	db, err := database.GetConnection()
	if err != nil {
		return false, err
	}
	defer db.Close()

	var count int
	err = db.QueryRow("select COUNT(*) FROM locations where name = ? and city = ? and state = ? and country = ?", newLocation.Name, newLocation.City, newLocation.State, newLocation.Country).Scan(&count)
	if err != nil {
		return false, err
	}

	return count == 0, nil
}