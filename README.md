# GDS Tech Assessment

## Setup

Create a `.env` file with the constants filled in. An example can be found at `.env.example`.

- `DB_USER`, `DB_PASSWORD`, `DB_NAME`, `DB_HOST`, `DB_PORT` refers to the configurations used in setting up a connection to the MySQL database

- `PORT` refers to the port used to run the Gin router. By default, the Gin router runs in localhost

After setting up, you can `go get .` and `go run .` to start the server.

## Schema

`db/init.sql` contains a script that was used to set up my local MySQL instance for local testing. It includes the schema definitions for `students`, `teachers` and an association table `teachers_students`.

## Structure

`main.go` provides the main entry point of the program.

`db/db.go` is in charge of initialising the connection to the database, as well as making queries or inserting to the database.

`controllers/handlers.go` includes the handlers for all the endpoints. `controllers/handlers_test.go` includes unit tests for the handler, whereby responses are mocked so as not to invoke the database.
