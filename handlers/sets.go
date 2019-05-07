package handlers

import (
	"fmt"
	"github.com/AnthonyNixon/setsisaw/auth"
	"github.com/AnthonyNixon/setsisaw/database"
	"github.com/AnthonyNixon/setsisaw/types"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func NewSet(c *gin.Context) {
	// Check Auth info
	claims, customErr := auth.GetUserInfo(c)
	if customErr != nil {
		c.JSON(customErr.StatusCode(), gin.H{"error": customErr.Description()})
		return
	}

	if !auth.IsEntitled(claims, "USER") {
		c.JSON(http.StatusForbidden, gin.H{"Error": fmt.Sprintf("User %s is not entitled to add a set.", claims.Username)})
		return
	}

	// If we're here, the user is authorized to add an artist.

	var newSet types.Set
	err := c.BindJSON(&newSet)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Bad JSON Input, could not bind."})
		return
	}

	userId, err := strconv.Atoi(claims.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not convert user id to an int"})
		return
	}
	newSet.User = userId

	unique, err := isNewSetUnique(newSet)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !unique {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Set already created"})
		return
	}


	db, err := database.GetConnection()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer db.Close()

	stmt, err := db.Prepare("insert into sets (user_id, artist_id, location_id) values(?,?,?);")
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	_, err = stmt.Exec(newSet.User, newSet.Artist, newSet.Location)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, gin.H{"user_id": newSet.User, "artist_id": newSet.Artist, "location_id": newSet.Location})
}

func isNewSetUnique(newSet types.Set) (bool, error) {
	db, err := database.GetConnection()
	if err != nil {
		return false, err
	}
	defer db.Close()

	var count int
	err = db.QueryRow("select COUNT(*) FROM sets where user_id = ? and artist_id = ? and location_id = ?", newSet.User, newSet.Artist, newSet.Location).Scan(&count)
	if err != nil {
		return false, err
	}

	return count == 0, nil
}