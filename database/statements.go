package database

// Sets
const GET_ALL_SETS = "select sets.id, user_id, artists.id, artists.name, locations.id, locations.name, sets.date, sets.rating, sets.genre, sets.length, sets.notes" +
	"FROM sets INNER JOIN artists ON artists.id = sets.artist_id " +
	"INNER JOIN locations ON locations.id = sets.location_id;"
const GET_ALL_SETS_FOR_USER_FORMAT = "select sets.id, user_id, artists.id, artists.name, locations.id, locations.name, sets.date, sets.rating, sets.genre, sets.length, sets.notes " +
	"FROM sets INNER JOIN artists ON artists.id = sets.artist_id " +
	"INNER JOIN locations ON locations.id = sets.location_id " +
	"WHERE user_id=%d;"
const IS_SET_UNIQUE_QUERY = "select COUNT(*) FROM sets where user_id = ? and artist_id = ? and location_id = ? and date = ?"
const INSERT_NEW_SET = `insert into sets (user_id, artist_id, location_id, date, rating, genre, length, notes) values(?,?,?,?,?,?,?,?);`

// Users
const GET_SPECIFIC_USER = `select id, username, email, IFNULL(first_name,""), IFNULL(last_name,""), role FROM users where id = ?;`
const GET_ALL_USERS = `select id, username, email, IFNULL(first_name,""), IFNULL(last_name,""), role FROM users;`

// Artists
const GET_ALL_ARTISTS = "select id, name FROM artists;"
const GET_SPECIFIC_ARTIST = "select id, name FROM artists where id = ?;"
const INSERT_NEW_ARTIST = `insert into artists (name, default_genre) values(?,?);`
const GET_ARTIST_DEFAULT_GENRE = `select default_genre FROM artists where name = ? or id = ?;`

// Locations
const GET_ALL_LOCATIONS = `select id, name, IFNULL(description,""), IFNULL(city,""), IFNULL(state,""), IFNULL(country,""), is_festival, IFNULL(year, 0000) FROM locations;`
const GET_SPECIFIC_LOCATION = `select id, name, IFNULL(description,""), IFNULL(city,""), IFNULL(state,""), IFNULL(country,""), is_festival, IFNULL(year, 0000) FROM locations WHERE id = ?;`
const INSERT_NEW_LOCATION = `insert into locations (name, description, city, state, country, is_festival, year) values(?,?,?,?,?,?,?);`
const IS_LOCATION_UNIQUE_QUERY = `select COUNT(*) FROM locations where name = ? and city = ? and state = ? and country = ? and IF(is_festival = TRUE, year = ?, true );`
