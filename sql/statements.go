package sql

const GET_ALL_SETS = "select sets.id, user_id, artists.id, artists.name, locations.id, locations.name " +
	"FROM sets INNER JOIN artists ON artists.id = sets.artist_id " +
	"INNER JOIN locations ON locations.id = sets.location_id;"

const GET_ALL_SETS_FOR_USER_FORMAT = "select sets.id, user_id, artists.id, artists.name, locations.id, locations.name " +
	"FROM sets INNER JOIN artists ON artists.id = sets.artist_id " +
	"INNER JOIN locations ON locations.id = sets.location_id " +
	"WHERE user_id=%d;"
