package handlers

import (
	"database/sql"
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

	unique, err := isNewLocationUnique(newLocation)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !unique {
		c.JSON(http.StatusBadRequest, gin.H{"error": "LocationId already created"})
		return
	}

	db, err := database.GetConnection()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer db.Close()

	stmt, err := db.Prepare(database.INSERT_NEW_LOCATION)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	_, err = stmt.Exec(newLocation.Name, newLocation.Description, newLocation.City, newLocation.State, newLocation.Country, newLocation.IsFestival, newLocation.Year)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
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

	location := types.Location{}
	locations := make([]types.Location, 0)

	db, err := database.GetConnection()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer db.Close()

	rows, err := db.Query(database.GET_ALL_LOCATIONS)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	for rows.Next() {
		err := rows.Scan(&location.Id, &location.Name, &location.Description, &location.City, &location.State, &location.Country, &location.IsFestival, &location.Year)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		locations = append(locations, location)
	}
	defer rows.Close()

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
	db, err := database.GetConnection()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer db.Close()

	result := db.QueryRow(database.GET_SPECIFIC_LOCATION, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed getting location - " + err.Error()})
		return
	}

	var location types.Location

	err = result.Scan(&location.Id, &location.Name, &location.Description, &location.City, &location.State, &location.Country, &location.IsFestival, &location.Year)
	if err != nil {
		// If an entry with the username does not exist, send an "Unauthorized"(401) status
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Location not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, location)

}

func isNewLocationUnique(newLocation types.Location) (bool, error) {
	db, err := database.GetConnection()
	if err != nil {
		return false, err
	}
	defer db.Close()

	var count int
	err = db.QueryRow(database.IS_LOCATION_UNIQUE_QUERY, newLocation.Name, newLocation.City, newLocation.State, newLocation.Country, newLocation.Year).Scan(&count)
	if err != nil {
		return false, err
	}

	return count == 0, nil
}
