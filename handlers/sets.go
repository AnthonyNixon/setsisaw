package handlers

import (
	"fmt"
	"github.com/AnthonyNixon/setsisaw/auth"
	"github.com/AnthonyNixon/setsisaw/database"
	"github.com/AnthonyNixon/setsisaw/types"
	"github.com/gin-gonic/gin"
	"log"
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
		c.JSON(http.StatusForbidden, gin.H{"Error": fmt.Sprintf("UserId %s is not entitled to add a set.", claims.Username)})
		return
	}

	// If we're here, the user is authorized to add a set.

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
	newSet.UserId = userId

	unique, err := isNewSetUnique(newSet)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !unique {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Set already created"})
		return
	}

	// TODO - If set location is not a festival, make sure there's a date.

	if newSet.Metadata.Genre == "" {
		defaultGenre, _ := getArtistDefaultGenre(newSet)
		newSet.Metadata.Genre = defaultGenre
		log.Printf("Default Genre: %s", defaultGenre)
	}

	db, err := database.GetConnection()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer db.Close()

	stmt, err := db.Prepare(database.INSERT_NEW_SET)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	_, err = stmt.Exec(newSet.UserId, newSet.ArtistId, newSet.LocationId, newSet.Date, newSet.Metadata.Rating, newSet.Metadata.Genre, newSet.Metadata.Length)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, gin.H{"user_id": newSet.UserId, "artist_id": newSet.ArtistId, "location_id": newSet.LocationId, "date": newSet.Date})
}

func GetSetsForCurrentUser(c *gin.Context) {
	// Check Auth info
	claims, customErr := auth.GetUserInfo(c)
	if customErr != nil {
		c.JSON(customErr.StatusCode(), gin.H{"error": customErr.Description()})
		return
	}

	if !auth.IsEntitled(claims, "USER") {
		c.JSON(http.StatusForbidden, gin.H{"Error": fmt.Sprintf("UserId %s is not entitled to get sets.", claims.Username)})
		return
	}

	// If we're here, the user is authorized to get sets for themselves.
	db, err := database.GetConnection()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer db.Close()

	user_id, err := strconv.ParseInt(claims.Id, 10, 64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not convert user_id to int, " + err.Error()})
		return
	}

	query := fmt.Sprintf(database.GET_ALL_SETS_FOR_USER_FORMAT, user_id)

	sendSets(query, c)
}

func GetAllSets(c *gin.Context) {
	// Check Auth info
	claims, customErr := auth.GetUserInfo(c)
	if customErr != nil {
		c.JSON(customErr.StatusCode(), gin.H{"error": customErr.Description()})
		return
	}

	if !auth.IsEntitled(claims, "EDITOR") {
		c.JSON(http.StatusForbidden, gin.H{"Error": fmt.Sprintf("UserId %s is not entitled to get all sets.", claims.Username)})
		return
	}

	// If we're here, the user is authorized to get all sets.

	sendSets(database.GET_ALL_SETS, c)
}

func sendSets(query string, c *gin.Context) {
	set := types.Set{}
	sets := make([]types.Set, 0)

	db, err := database.GetConnection()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer db.Close()

	rows, err := db.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	for rows.Next() {
		err := rows.Scan(&set.Id, &set.UserId, &set.ArtistId, &set.ArtistName, &set.LocationId, &set.LocationName, &set.Date, &set.Metadata.Rating, &set.Metadata.Genre, &set.Metadata.Length)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		sets = append(sets, set)
	}
	defer rows.Close()

	c.JSON(http.StatusOK, gin.H{"sets": sets, "count": len(sets)})
}

func isNewSetUnique(newSet types.Set) (bool, error) {
	db, err := database.GetConnection()
	if err != nil {
		return false, err
	}
	defer db.Close()

	var count int
	err = db.QueryRow(database.IS_SET_UNIQUE_QUERY, newSet.UserId, newSet.ArtistId, newSet.LocationId, newSet.Date).Scan(&count)
	if err != nil {
		return false, err
	}

	return count == 0, nil
}

func getArtistDefaultGenre(newSet types.Set) (string, error) {
	db, err := database.GetConnection()
	if err != nil {
		return "", err
	}
	defer db.Close()

	var genre string
	err = db.QueryRow(database.GET_ARTIST_DEFAULT_GENRE, newSet.ArtistName, newSet.ArtistId).Scan(&genre)
	if err != nil {
		return "", err
	}

	return genre, nil
}
