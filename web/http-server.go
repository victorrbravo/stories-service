package web

import (
	"fmt"
	"log"
	"net/http"

	"github.com/victorrbravo/stories-service/constants"
	"github.com/victorrbravo/stories-service/database"
	handlers "github.com/victorrbravo/stories-service/generated/handlers/users"
	"github.com/victorrbravo/stories-service/generated/models"
	"github.com/victorrbravo/stories-service/handlers/auth"
	"github.com/victorrbravo/stories-service/handlers/social"
	"github.com/victorrbravo/stories-service/handlers_utils"
)

const (
	Authentication = "/auth"
	Users          = "/users"
	Social         = "/social"
)

type Server struct{}

func CreateRouter(newQuery database.IQuery) handlers_utils.Router {
	router := handlers_utils.NewRouter()
	newUser := models.TaUser{}
	router.AddNoAuthenticated(Authentication, "/login", auth.UserLoginHandler{Query: newQuery, User: newUser},
		http.MethodPost, newQuery)

	router.AddAuthenticated(Social, "/add-step/{nameId}", social.AddStepHandler{Query: newQuery}, http.MethodPost, newQuery)
	router.AddAuthenticated(Users, "/", handlers.GetAllHandler{Query: newQuery}, http.MethodGet, newQuery)
	router.AddAuthenticated(Users, "/", handlers.CreateHandler{Query: newQuery}, http.MethodPost, newQuery)
	router.AddAuthenticated(Users, "/", handlers.DeleteHandler{Query: newQuery}, http.MethodDelete, newQuery)
	return router
}

func NewServer() (Server, error) {
	return Server{}, nil
}

func (server Server) Listen(host string, port int, swaggerDir string, newQuery database.IQuery) {
	var err error
	router := CreateRouter(newQuery)
	// Setup swagger
	log.Printf("%s starting server  @ %s:%d", constants.ApplicationName, host, port)
	// Start the http server
	err = http.ListenAndServe(fmt.Sprintf("%s:%d", host, port), router)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}
