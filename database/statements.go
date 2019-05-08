package database

// Sets
const GET_ALL_SETS = "select sets.id, user_id, artists.id, artists.name, locations.id, locations.name " +
	"FROM sets INNER JOIN artists ON artists.id = sets.artist_id " +
	"INNER JOIN locations ON locations.id = sets.location_id;"
const GET_ALL_SETS_FOR_USER_FORMAT = "select sets.id, user_id, artists.id, artists.name, locations.id, locations.name " +
	"FROM sets INNER JOIN artists ON artists.id = sets.artist_id " +
	"INNER JOIN locations ON locations.id = sets.location_id " +
	"WHERE user_id=%d;"

// Users
const GET_SPECIFIC_USER = `select id, username, email, IFNULL(first_name,""), IFNULL(last_name,""), role FROM users where id = ?;`
const GET_ALL_USERS = `select id, username, email, IFNULL(first_name,""), IFNULL(last_name,""), role FROM users;`

// Artists
const GET_ALL_ARTISTS = "select id, name FROM artists;"
const GET_SPECIFIC_ARTIST = "select id, name FROM artists where id = ?;"

// Locations
const GET_ALL_LOCATIONS = `select id, name, IFNULL(description,""), IFNULL(city,""), IFNULL(state,""), IFNULL(country,"") FROM locations;`
const GET_SPECIFIC_LOCATION = `select id, name, IFNULL(description,""), IFNULL(city,""), IFNULL(state,""), IFNULL(country,"") FROM locations WHERE id = ?;`
