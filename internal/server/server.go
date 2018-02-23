package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/natethinks/instruu-api/internal/respond"
	"github.com/natethinks/instruu-api/internal/store"
)

// Server abstracts handlers and the store service
type Server struct {
	sto     store.Service
	handler http.Handler
}

// New creates a new server from a store and populates the handler
func New(sto store.Service) *Server {
	s := &Server{sto: sto}

	router := mux.NewRouter()

	router.Handle("/user", handlers.LoggingHandler(os.Stdout, allowedMethods(
		[]string{"OPTIONS", "GET", "POST"},
		handlers.MethodHandler{
			"GET":  http.HandlerFunc(s.getUsers),
			"POST": http.HandlerFunc(s.createUser), // created
		})))

	router.Handle("user/{id}", handlers.LoggingHandler(os.Stdout, allowedMethods(
		[]string{"OPTIONS", "GET", "PUT", "PATCH", "DELETE"},
		handlers.MethodHandler{
			"GET": http.HandlerFunc(s.getUser), // created
			//"PUT":    http.HandlerFunc(s.putUser),
			"PATCH":  http.HandlerFunc(s.patchUser),
			"DELETE": http.HandlerFunc(s.deleteUser),
		})))

	router.Handle("/resource", allowedMethods(
		[]string{"OPTIONS", "GET", "POST"},
		handlers.MethodHandler{
			// get resources will have query params since this should be reusable
			"GET":  http.HandlerFunc(s.getResources),
			"POST": http.HandlerFunc(s.createResource),
		}))

	router.Handle("/resource/{id}", allowedMethods(
		[]string{"OPTIONS", "GET", "PUT", "PATCH", "DELETE"},
		handlers.MethodHandler{
			"GET":    http.HandlerFunc(s.getResource),
			"PUT":    http.HandlerFunc(s.putResource),
			"PATCH":  http.HandlerFunc(s.patchResource),
			"DELETE": http.HandlerFunc(s.deleteResource),
		}))

	s.handler = limitBody(defaultHeaders(router))

	return s
}

// Run starts the server listening on what address is specified
func (s *Server) Run(addr string) error {
	return http.ListenAndServe(addr, s.handler)
}

func commaify(ss []string) (out string) {
	for i, s := range ss {
		out += s
		if i != len(ss)-1 {
			out += ","
		}
	}
	return
}

func allowedMethods(methods []string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Methods", commaify(methods))

		next.ServeHTTP(w, r)
	})
}

// User Functions
//
//
//

func (s *Server) getUsers(w http.ResponseWriter, r *http.Request) {
	fmt.Println("getUsers() called")
	//users, err := s.sto.GetUsers()
	//if err != nil {
	//	if err == store.ErrNoResults {
	//		users = []store.User{}
	//	} else {
	//		respond.JSON(w, err)
	//		return
	//	}
	//}

	//respond.JSON(w, users)
	return
}

func (s *Server) createUser(w http.ResponseWriter, r *http.Request) {
	var user store.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	fmt.Println(user)

	id, err := s.sto.CreateUser(user)
	if err != nil {
		respond.JSON(w, err)
		return
	}

	respond.JSON(w, map[string]int64{"id": id})
	return
}

func (s *Server) getUser(w http.ResponseWriter, r *http.Request) {
	rawID := mux.Vars(r)["id"]
	fmt.Println("s.getUser() called")

	id, err := strconv.ParseInt(rawID, 10, 64)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	user, err := s.sto.GetUser(id)
	if err != nil {
		if err == store.ErrNoResults {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		} else {
			respond.JSON(w, err)
		}
		return
	}

	respond.JSON(w, user)
	return
}

func (s *Server) putUser(w http.ResponseWriter, r *http.Request) {

	return
}

func (s *Server) patchUser(w http.ResponseWriter, r *http.Request) {

	return
}

func (s *Server) deleteUser(w http.ResponseWriter, r *http.Request) {

	return
}

// Resource Functions
//
//
//

func (s *Server) getResources(w http.ResponseWriter, r *http.Request) {

	return
}

func (s *Server) createResource(w http.ResponseWriter, r *http.Request) {

	return
}

func (s *Server) getResource(w http.ResponseWriter, r *http.Request) {

	return
}

func (s *Server) putResource(w http.ResponseWriter, r *http.Request) {

	return
}

func (s *Server) patchResource(w http.ResponseWriter, r *http.Request) {

	return
}

func (s *Server) deleteResource(w http.ResponseWriter, r *http.Request) {

	return
}

func defaultHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		next.ServeHTTP(w, r)
	})
}

func limitBody(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, 1048576)
		next.ServeHTTP(w, r)
	})
}
