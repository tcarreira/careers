# superhero
SuperHero API - Go (inspired by superheroapi.com)

This is being made in the context of https://github.com/levpay/careers#desafio


# Running

using docker-compose: 

```
docker-compose up --build -d
``` 

then access http://localhost:8080

If you have a PostgreSQL instance running on `postgres:password@localhost:5432/postgres`, you may use 
```
go build && ./superhero
``` 


# Testing

Run unit tests without a PostgeSQL database (on every save): 
```
go test -v ./...
```

## testing with a database

For performing full tests (with a real database), run it with docker:
```
 docker run --rm --name pgsql -d -p 5432:5432 -e POSTGRES_PASSWORD=password postgres:12-alpine

 go test -v ./... -tags sql

 docker stop pgsql
```

If you have a PostgtreSQL instance you want to run test against:

```
DB_HOST=db DB_USER=user DB_PASS=pass go test -v ./... -tags sql
``` 


# Features

- [X] Create new Super(Hero/Vilan)
- [X] Get Super list
- [X] Get Super(Heroes) list
- [X] Get Super(Vilans) list
- [X] Search by name
- [X] Search by uuid
- [X] Delete Super
- [X] Super Groups


# TODO

- [X] Create Repository and setup Go
- [X] GIN hello world
- [X] Setup docker (it will be easier for Database integration)
- [X] Add Postgres and go-pg
- [X] Setup command-line admin (create db-schema)
- [X] Create new Super (Hero/Vilan)
    - [X] Super structs
    - [X] POST
    - [X] GET
- [X] GET filters (by type, name, uuid)
- [X] DELETE Super
- [X] Implement Groups
    - [X] Group struct
    - [X] POST New Group
    - [X] GET Group
- [X] Specifics
    - [X] GET super retrieves List of Groups
    - [X] GET super retrieves number of relatives
    - [X] Create a Super(Hero/Vilan) by name only