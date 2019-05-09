package utils

import (
	"database/sql"
	"fmt"
	"github.com/AnthonyNixon/setsisaw/customerrors"
	"github.com/AnthonyNixon/setsisaw/database"
	"github.com/AnthonyNixon/setsisaw/types"
	"log"
	"net/http"
)

func NewLocation(location types.Location) types.Error {
	unique, err := isNewLocationUnique(location)
	if err != nil {
		return customerrors.New(http.StatusInternalServerError, "could not determine if location is unique, "+err.Error())
	}

	if !unique {
		return customerrors.New(http.StatusConflict, fmt.Sprintf("location with ID %d already created", location.Id))
	}

	db, err := database.GetConnection()
	if err != nil {
		return customerrors.New(http.StatusInternalServerError, "could not get database connection, "+err.Error())
	}
	defer db.Close()

	stmt, err := db.Prepare(database.INSERT_NEW_LOCATION)
	if err != nil {
		return customerrors.New(http.StatusInternalServerError, "could not prepare db statement, "+err.Error())
	}

	_, err = stmt.Exec(location.Name, location.Description, location.City, location.State, location.Country, location.IsFestival, location.Year)
	if err != nil {
		return customerrors.New(http.StatusInternalServerError, "error executing insert statement, "+err.Error())
	}

	return nil
}

func GetLocation(id string) (types.Location, types.Error) {
	var location types.Location
	// if we made it here, we're good to go.
	db, err := database.GetConnection()
	if err != nil {
		return location, customerrors.New(http.StatusInternalServerError, err.Error())
	}
	defer db.Close()

	result := db.QueryRow(database.GET_SPECIFIC_LOCATION, id)
	if err != nil {
		return location, customerrors.New(http.StatusInternalServerError, "failed getting location, "+err.Error())
	}

	err = result.Scan(&location.Id, &location.Name, &location.Description, &location.City, &location.State, &location.Country, &location.IsFestival, &location.Year)
	if err != nil {
		// If an entry with the username does not exist, send an "Unauthorized"(401) status
		if err == sql.ErrNoRows {
			return location, customerrors.New(http.StatusNotFound, "location not found")
		}

		return location, customerrors.New(http.StatusInternalServerError, err.Error())
	}

	return location, nil
}

func GetAllLocations() ([]types.Location, types.Error) {
	location := types.Location{}
	locations := make([]types.Location, 0)

	db, err := database.GetConnection()
	if err != nil {
		return nil, customerrors.New(http.StatusInternalServerError, "could not get database connection, "+err.Error())
	}
	defer db.Close()

	rows, err := db.Query(database.GET_ALL_LOCATIONS)
	if err != nil {
		return nil, customerrors.New(http.StatusInternalServerError, "could not query database, "+err.Error())
	}

	for rows.Next() {
		err := rows.Scan(&location.Id, &location.Name, &location.Description, &location.City, &location.State, &location.Country, &location.IsFestival, &location.Year)
		if err != nil {
			return nil, customerrors.New(http.StatusInternalServerError, "could not scan location row, "+err.Error())
		}
		locations = append(locations, location)
	}
	rows.Close()

	return locations, nil
}

func UpdateLocation(id string, location types.Location) types.Error {
	db, err := database.GetConnection()
	if err != nil {
		return customerrors.New(http.StatusInternalServerError, "could not get database connection, "+err.Error())
	}
	defer db.Close()

	unique, err := isLocationUpdateInfoUnique(location)
	if err != nil {
		return customerrors.New(http.StatusInternalServerError, "could not determine if update is unique")
	}

	if !unique {
		return customerrors.New(http.StatusConflict, "location already exists")
	}

	stmt, err := db.Prepare(database.UPDATE_LOCATION)
	if err != nil {
		return customerrors.New(http.StatusInternalServerError, "could not prepare statement, "+err.Error())
	}

	_, err = stmt.Exec(location.Name, location.Description, location.City, location.State, location.Country, location.IsFestival, location.Year, location.Id)
	if err != nil {
		return customerrors.New(http.StatusInternalServerError, "could not execute statement, "+err.Error())
	}

	return nil
}

func isLocationUpdateInfoUnique(location types.Location) (bool, error) {
	db, err := database.GetConnection()
	if err != nil {
		return false, err
	}
	defer db.Close()

	var count int
	err = db.QueryRow(database.IS_LOCATION_UPDATE_UNIQUE, location.Id, location.Name, location.City, location.State, location.Country, location.Year).Scan(&count)
	if err != nil {
		return false, err
	}

	log.Printf("Unique count: %d", count)
	return count == 0, nil
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
