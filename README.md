# setsisaw
The front and back end for SetsISaw
🐩🐩🐩🐩

## Building
To build the project, simply use docker. The command is:
`docker build -t setsisaw-api:latest`


## Running
To run the project you will need a few environment variables set. The Required variables are:
- `JWT_SIGNING_KEY`: The signing key to be used for JWT tokens. This can be any string.
- `SETSISAW_DB_HOST`: the host of the database to be used.
- `SETSISAW_DB_NAME`: the name of the database being used.
- `SETSISAW_DB_USER`: the username of the database account.
- `SETSISAW_DB_PASS`: the password for the database account.

To run once these values are set use a command like:
`docker run -p 8080:8080 -e JWT_SIGNING_KEY=... -e SETSISAW_DB_HOST=... -e SETSISAW_DB_NAME=... -e SETSISAW_DB_USER=... -e SETSISAW_DB_PASS='...' setsisaw:latest`