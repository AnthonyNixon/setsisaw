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

func NewArtist(c *gin.Context) {
	// Check Auth info
	claims, customErr := auth.GetUserInfo(c)
	if customErr != nil {
		c.JSON(customErr.StatusCode(), gin.H{"error": customErr.Description()})
		return
	}

	// TODO: This shouldn't be an editor only function. Anyone should be able to add an artist. This is just for testing claims.
	if !auth.IsEntitled(claims, "EDITOR") {
		c.JSON(http.StatusForbidden, gin.H{"Error": fmt.Sprintf("UserId %s is not entitled to add an artist.", claims.Username)})
		return
	}

	// If we're here, the user is authorized to add an artist.

	var newArtist types.Artist
	err := c.BindJSON(&newArtist)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Bad JSON Input, could not bind."})
		return
	}

	unique, err := isNewArtistUnique(newArtist)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !unique {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ArtistId already created"})
		return
	}

	db, err := database.GetConnection()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer db.Close()

	stmt, err := db.Prepare(database.INSERT_NEW_ARTIST)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	_, err = stmt.Exec(newArtist.Name, newArtist.DefaultGenre)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, gin.H{"name": newArtist.Name})
}

func GetAllArtists(c *gin.Context) {
	// Check Auth info
	claims, customErr := auth.GetUserInfo(c)
	if customErr != nil {
		c.JSON(customErr.StatusCode(), gin.H{"error": customErr.Description()})
		return
	}

	if !auth.IsEntitled(claims, "USER") {
		c.JSON(http.StatusForbidden, gin.H{"Error": fmt.Sprintf("User %s is not entitled to get all artists.", claims.Username)})
		return
	}

	// If we're here, the user is authorized to get all artists.

	artist := types.Artist{}
	artists := make([]types.Artist, 0)

	db, err := database.GetConnection()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer db.Close()

	rows, err := db.Query(database.GET_ALL_ARTISTS)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	for rows.Next() {
		err := rows.Scan(&artist.Id, &artist.Name)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		artists = append(artists, artist)
	}
	defer rows.Close()

	c.JSON(http.StatusOK, gin.H{"artists": artists, "count": len(artists)})

}

func GetArtist(c *gin.Context) {
	id := c.Param("id")

	// Check Auth info
	claims, customErr := auth.GetUserInfo(c)
	if customErr != nil {
		c.JSON(customErr.StatusCode(), gin.H{"error": customErr.Description()})
		return
	}

	if !auth.IsEntitled(claims, "USER") {
		c.JSON(http.StatusForbidden, gin.H{"Error": fmt.Sprintf("UserId %s is not entitled to get artist info.", claims.Username)})
		return
	}

	// if we made it here, we're good to go.
	db, err := database.GetConnection()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer db.Close()

	result := db.QueryRow(database.GET_SPECIFIC_ARTIST, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed getting artist - " + err.Error()})
		return
	}

	var artist types.Artist

	err = result.Scan(&artist.Id, &artist.Name)
	if err != nil {
		// If an entry with the username does not exist, send an "Unauthorized"(401) status
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Artist not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, artist)

}

func isNewArtistUnique(newArtist types.Artist) (bool, error) {
	db, err := database.GetConnection()
	if err != nil {
		return false, err
	}
	defer db.Close()

	var count int
	err = db.QueryRow("select COUNT(*) FROM artists where name = ?", newArtist.Name).Scan(&count)
	if err != nil {
		return false, err
	}

	return count == 0, nil
}
