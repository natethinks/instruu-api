package postgres

import (
	"database/sql"
	"fmt"

	"github.com/natethinks/instruu-api/internal/store"

	// for the postgres sql driver
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

// this doesn't satisfy the store.Service interface until ALL functions are built
type service struct {
	db *sql.DB
}

// Options holds information for connecting to a postgresql server
type Options struct {
	User, Pass string
	Host       string
	Port       int
	DBName     string
	SSLMode    string
}

func (o Options) connectionInfo() string {
	return fmt.Sprintf("host='%s' port='%d' user='%s' password='%s' dbname='%s' sslmode='%s'",
		o.Host, o.Port, o.User, o.Pass, o.DBName, o.SSLMode)
}

const usersTableCreationQuery = `
CREATE TABLE IF NOT EXISTS users (
	id          SERIAL PRIMARY KEY,
	username	varchar(256),
	email		varchar(256),
	firstName	varchar(256),
	lastName	varchar(256),
	isVerified 	BOOLEAN
)`

const resourcesTableCreationQuery = `
CREATE TABLE IF NOT EXISTS resources (
	id          SERIAL PRIMARY KEY,
	name		varchar(256),
	description text,
	url			varchar(256),
	approved	BOOLEAN NOT NULL DEFAULT FALSE,
	creator		integer references users(id),
	PRIMARY KEY (url)
)`

const tagsTableCreationQuery = `
CREATE TABLE IF NOT EXISTS tags (
	id			SERIAL PRIMARY KEY,
	name		varchar(256),
	PRIMARY KEY name
)`

const tagTableCreationQuery = `
CREATE TABLE IF NOT EXISTS tag (
	id			SERIAL PRIMARY KEY,
	resource	integer references resources(id) ON DELETE CASCADE,
	tag			integer references tags(id) ON DELETE CASCADE,
	PRIMARY KEY (resource, tag)
)`

// New connects to a postgres server with specified options and returns a store.Service
func New(options Options) (store.Service, error) {
	db, err := sql.Open("postgres", options.connectionInfo())
	if err != nil {
		return nil, errors.Wrap(err, "connecting to postgres database")
	}

	_, err = db.Exec(usersTableCreationQuery)
	if err != nil {
		return nil, errors.Wrap(err, "creating todos table")
	}

	_, err = db.Exec(resourcesTableCreationQuery)
	if err != nil {
		return nil, errors.Wrap(err, "creating resources table")
	}

	return &service{db: db}, nil
}

// User store functions
func (s *service) CreateUser(user store.User) (id int64, err error) {
	fmt.Println(user)
	err = s.db.QueryRow(
		"INSERT INTO users (username, email, firstname, lastname) VALUES ($1, $2, $3, $4) RETURNING id",
		user.Username, user.Email, user.FirstName, user.LastName).Scan(&id)
	return
}

func (s *service) GetUser(id int64) (user store.User, err error) {
	fmt.Println("s.GetUser() called")
	user = store.User{ID: id}
	err = s.db.QueryRow("SELECT username, email FROM users WHERE id = $1", id).Scan(
		&user.Username, &user.Email)
	if err == sql.ErrNoRows {
		err = store.ErrNoResults
	}
	return user, err
}

func (s *service) PatchUser(user store.User) (err error) {
	fmt.Println("s.PatchUser() called")
	fmt.Println(user)
	return nil
}

func (s *service) DeleteUser(id int64) (err error) {
	fmt.Println("s.DeleteUser() called")
	fmt.Println(id)
	return nil
}

func (s *service) GetUsers() (users []store.User, err error) {
	fmt.Println("s.GetUsers() called")

	rows, err := s.db.Query("SELECT * FROM users")
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	} else if err == sql.ErrNoRows {
		err = store.ErrNoResults
	}
	defer rows.Close()

	for rows.Next() {
		var user store.User
		if err = rows.Scan(&user); err != nil {
			return users, err
		}
		users = append(users, user)
	}
	return
}

func (s *service) Close() error {
	return s.db.Close()
}
