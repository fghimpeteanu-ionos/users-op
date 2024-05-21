package controller

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	host = "localhost"
	//host     = "my-db-postgresql"
	port       = 25432
	user       = "user-op"
	password   = "user-op-pretest"
	dbname     = "user-db"
	driverName = "postgres"

	sqlQueryInsertUser = `
	INSERT INTO users (firstName, lastName, age, address, email)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING uuid`
	sqlQueryReadUser = `
	SELECT uuid, firstName, lastName, age, address, email
	FROM users
	WHERE uuid = $1
	`
)

var connInfo = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
	host, port, user, password, dbname)

type UserPersistenceImpl struct{}

func (up *UserPersistenceImpl) Persist(user *CreateUserE) (string, error) {
	dbConnection := createDBConnection()
	defer dbConnection.Close()
	testDBConnection(dbConnection)

	userId := persist(user, dbConnection)

	return userId, nil
}

func (up *UserPersistenceImpl) Read(uuid string) (*ReadUserE, error) {
	dbConnection := createDBConnection()
	defer dbConnection.Close()
	testDBConnection(dbConnection)

	readUserE := read(uuid, dbConnection)

	return readUserE, nil
}

func createDBConnection() *sql.DB {
	db, err := sql.Open(driverName, connInfo)
	if err != nil {
		panic(err)
	}
	return db
}

func testDBConnection(db *sql.DB) {
	err := db.Ping()
	if err != nil {
		panic(err)
	}
}

func persist(user *CreateUserE, db *sql.DB) string {
	userId := ""
	err := db.QueryRow(sqlQueryInsertUser,
		user.firstName,
		user.lastName,
		user.age,
		user.address,
		user.email,
	).Scan(&userId)
	if err != nil {
		panic(err)
	}
	return userId
}

func read(uuid string, db *sql.DB) *ReadUserE {
	readUserE := ReadUserE{}
	err := db.QueryRow(sqlQueryReadUser, uuid).Scan(
		&readUserE.id,
		&readUserE.firstName,
		&readUserE.lastName,
		&readUserE.age,
		&readUserE.address,
		&readUserE.email)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		panic(err)
	} else if err != nil {
		return nil
	}

	return &readUserE
}
