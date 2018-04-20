package postgres

import (
	"database/sql"
	"fmt"

	"github.com/natethinks/instruu-api/internal/auth"
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
	password	varchar(256),
	isVerified 	BOOLEAN NOT NULL DEFAULT FALSE
)`

const resourcesTableCreationQuery = `
CREATE TABLE IF NOT EXISTS resources (
	id          SERIAL PRIMARY KEY,
	name		varchar(256),
	description text,
	url			varchar(256) UNIQUE,
	approved	BOOLEAN NOT NULL DEFAULT FALSE,
	submitter	integer references users(id),
	deleted		BOOLEAN NOT NULL DEFAULT FALSE
)`

const tagsTableCreationQuery = `
CREATE TABLE IF NOT EXISTS tags (
	id			SERIAL PRIMARY KEY,
	name		varchar(256) UNIQUE
)`

const tagTableCreationQuery = `
CREATE TABLE IF NOT EXISTS tag (
	id			SERIAL PRIMARY KEY,
	resource	integer references resources(id) ON DELETE CASCADE,
	tag			integer references tags(id) ON DELETE CASCADE,
	CONSTRAINT  unq_res_tag UNIQUE(resource, tag)
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

// Authentication Functions

func (s *service) Auth(user store.User) (jwt string, err error) {
	fmt.Println("auth running")
	fmt.Println(user)
	// auth needs to query and grab the users password hash, then send it to the verify
	// password function in auth then it needs to generate a JWT and return it to the user
	err = s.db.QueryRow("SELECT password FROM users WHERE username = $1", user.Username).Scan(&user.PasswordHash)
	if err != nil {
		fmt.Println("probably can't find that user")
		return jwt, err
	}

	authSuccess := auth.VerifyPassword(user.PasswordHash, []byte(user.Password))
	if !authSuccess {
		fmt.Println("password doesn't match")
		err = errors.New("Incorrect Password")
		return jwt, err
	}

	fmt.Println(user)
	return jwt, nil
}

// User store functions
func (s *service) CreateUser(user store.User) (id int64, err error) {
	// generate password hash before storing
	user.PasswordHash = auth.GeneratePasswordHash([]byte(user.Password))
	fmt.Println(user)
	err = s.db.QueryRow(
		"INSERT INTO users (username, email, firstname, lastname, password) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		user.Username, user.Email, user.FirstName, user.LastName, user.PasswordHash).Scan(&id)
	return id, err
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

	rows, err := s.db.Query("SELECT id, username, email, firstname, lastname, isVerified FROM users")
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	} else if err == sql.ErrNoRows {
		err = store.ErrNoResults
	}
	defer rows.Close()

	for rows.Next() {
		var user store.User
		if err = rows.Scan(&user.ID, &user.Username, &user.Email, &user.FirstName, &user.LastName, &user.Verified); err != nil {
			return users, err
		}
		users = append(users, user)
	}
	return
}

func (s *service) CheckUsername(user store.User) error {
	fmt.Println(user.Username)
	var id int
	err := s.db.QueryRow("SELECT id FROM users WHERE username = $1", user.Username).Scan(&id)
	if err == sql.ErrNoRows {
		return nil
	}

	err = errors.New("Username exists, cannot create")
	return err
}

func (s *service) Close() error {
	return s.db.Close()
}

// Resource Functions

func (s *service) CreateResource(resource store.Resource) (id int64, err error) {
	fmt.Println(resource)
	err = s.db.QueryRow(
		"INSERT INTO resources (name, description, url, submitter) VALUES ($1, $2, $3, $4) RETURNING id",
		resource.ID, resource.Description, resource.URL, resource.Submitter).Scan(&id)
	return id, err
}

func (s *service) GetResource(id int64) (resource store.Resource, err error) {

	resource = store.Resource{ID: id}
	err = s.db.QueryRow("SELECT id, name, description, url, submitter, approved, submitter FROM resources WHERE id = $1 AND deleted = false", id).Scan(&resource.Name, &resource.Description, &resource.URL, &resource.Submitter)
	if err == sql.ErrNoRows {
		err = store.ErrNoResults
	}
	return resource, err
}

// GetResources might better encompass groups as well with a query param
func (s *service) GetResources(query map[string][]string) (resources []store.Resource, err error) {

	// need to do some query builder stuff and check the query params
	fmt.Println(query)

	rows, err := s.db.Query("SELECT id, name, description, url FROM resources")
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	} else if err == sql.ErrNoRows {
		err = store.ErrNoResults
	}
	defer rows.Close()

	for rows.Next() {
		var resource store.Resource
		if err = rows.Scan(&resource.ID, &resource.Name, &resource.Description, &resource.URL); err != nil {
			return resources, err
		}
		resources = append(resources, resource)
	}
	return
}

func (s *service) UpdateResource(resource store.Resource) (err error) {
	// name, description, and URL are all game to be updated. maybe not url?
	// but since i'm getting an entire resource passed to me I might as well run the
	// whole update
	fmt.Println(resource)
	res, err := s.db.Exec("UPDATE resources SET (name, description, url, approved, deleted) VALUES ($1, $2, $3, $4, $5)", resource.Name, resource.Description, resource.URL, resource.Approved, resource.Deleted)
	fmt.Println(res)
	if err != nil {
		return err
	}

	// Need to add something to track history of changes here, the original submitter should stay the same, but on this patch request the submitter value will be the editor, and that will be stored in it's own table.

	return err
}

func (s *service) DeleteResource(ID int64) (err error) {

	return err
}
