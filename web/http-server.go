package web

import (
	"fmt"
	"log"
	"net/http"
	"stories/constants"
	"stories/database"
	"stories/generated/handlers/users"
	"stories/generated/models"
	"stories/handlers/auth"
	"stories/handlers_utils"
)

const (
	Authentication = "/auth"
	Users          = "/users"
)

type Server struct {
}

func CreateRouter(newQuery database.IQuery) handlers_utils.Router {
	var router = handlers_utils.NewRouter()
	newUser := models.TaUser{}
	router.AddNoAuthenticated(Authentication, "/login", auth.UserLoginHandler{Query: newQuery, User: newUser}, http.MethodPost, newQuery)

	router.AddAuthenticated(Users, "/", handlers.GetAllHandler{Query: newQuery}, http.MethodGet, newQuery)
	router.AddAuthenticated(Users, "/", handlers.CreateHandler{Query: newQuery}, http.MethodPost, newQuery)
	router.AddAuthenticated(Users, "/", handlers.DeleteHandler{Query: newQuery}, http.MethodDelete, newQuery)
	router.AddAuthenticated(Users, "/", handlers.UpdateHandler{Query: newQuery}, http.MethodPut, newQuery)
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
