package store

import "fmt"

// ErrNoResults is a generic error of sql.ErrNoRows
var ErrNoResults = fmt.Errorf("no results returned")

// Service contains all functions to int64erface with a store
type Service interface {
	// Authentication Functions
	Auth(user User) (string, error)
	// User Functions
	CreateUser(user User) (int64, error)
	GetUser(ID int64) (User, error)
	PatchUser(user User) error
	DeleteUser(ID int64) error
	GetUsers() ([]User, error)
	CheckUsername(user User) error
	//GetUsers() ([]User, error)
	//GetUserGroup(ID int64) ([]User, error)
	//UpdateUser(user User) error
	//DeleteUser(ID int64) error
	//VerifyUser(ID int64) error
	// Resource Functions
	CreateResource(resource Resource) (int64, error)
	GetResource(ID int64) (Resource, error)
	GetResources(query map[string][]string) ([]Resource, error)
	// probably consolidating this with GetResources and query params
	//GetResourceGroup(ID int64) ([]Resource, error)
	UpdateResource(resource Resource) error
	DeleteResource(ID int64) error
	Close() error
}

// Resource is a single learning resource submitted by a user
type Resource struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	URL         string `json:"url"`
	Approved    bool   `json:"approved"`
	Submitter   int64  `json:"submitter"`
	Deleted     bool   `json:"deleted"`
}

// User Represents every user that has signed up for Instruu
type User struct {
	ID           int64  `json:"id"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	Verified     bool   `json:"verified"`
	Password     string `json:"password"`
	PasswordHash string
}

// Secure user is a struct for return values so that password will not accidentally be returned
type SecureUser struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Verified  bool   `json:"verified"`
}
