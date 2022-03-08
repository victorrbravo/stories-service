package web

import (
	"fmt"
	"log"
	"net/http"
	"stories/constants"
	"stories/database"
	"stories/generated/handlers/users"
	"stories/handlers_utils"
)

type Server struct {
}

func CreateRouter(newQuery database.IQuery) handlers_utils.Router {
	var router = handlers_utils.NewRouter()
	router.AddNoAuthenticated("/users", "/", handlers.GetAllHandler{Query: newQuery}, http.MethodGet, newQuery)
	router.AddNoAuthenticated("/users", "/", handlers.CreateHandler{Query: newQuery}, http.MethodPost, newQuery)
	router.AddNoAuthenticated("/users", "/", handlers.DeleteHandler{Query: newQuery}, http.MethodDelete, newQuery)
	router.AddNoAuthenticated("/users", "/", handlers.UpdateHandler{Query: newQuery}, http.MethodPut, newQuery)
	return router
}

func NewServer() (Server, error) {
	return Server{}, nil
}
func (server Server) Listen(host string, port int, swaggerDir string, newQuery database.IQuery) {
	var err error
	router := CreateRouter(newQuery)
	// Setup swagger
	log.Println(fmt.Sprintf("%s starting server  @ %s:%d", constants.ApplicationName, host, port))
	// Start the http server
	err = http.ListenAndServe(fmt.Sprintf("%s:%d", host, port), router)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}
