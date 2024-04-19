package auth

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"stories/database"
	"stories/generated/models"
	"stories/types"
	"strings"

	"github.com/dgrijalva/jwt-go"

	"stories/constants"
	"stories/handlers_utils"
	"stories/helpers"
)

type VerificationData struct {
	Name string `json:"name"`
}

type VerifyEmailHandler struct {
	Query database.IQuery
}

func (handler VerifyEmailHandler) IsValid() error {
	return nil
}

// Check if endpoint is enabled
func (handler VerifyEmailHandler) IsEnabled() bool {
	return true
}

// Verify users email after registration
// **Note** registration is not completed unless the users do verify they own this email
// You must add a parameter to Query String (URL) in this format: ?code=<CODE_SENT>
func (handler VerifyEmailHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	var err error
	// var email string
	var code string
	var jwtToken jwt.MapClaims
	var query url.Values
	query = request.URL.Query()
	code = query.Get("code")
	if len(code) == 0 {
		handlers_utils.RespondWithError(writer, helpers.NewHttpError(http.StatusBadRequest, fmt.Sprintf(constants.RequiredFieldIsEmpty, "code")))
		return
	}
	jwtToken, err = helpers.TokenValid(strings.TrimSuffix(code, "\\r\\n"))
	if err != nil {
		handlers_utils.RespondWithError(writer, helpers.NewHttpError(http.StatusUnauthorized, constants.TokenIsNotValid))
		return
	}
	if _, ok := jwtToken["email"]; !ok {
		handlers_utils.RespondWithError(writer, helpers.NewHttpError(http.StatusUnauthorized, constants.TokenNotParseCorrectly))
		return
	}
	email := jwtToken["email"].(string)
	// email = jwtToken["email"].(string)
	users, err := models.GetTaUserByEmail(handler.Query, types.NewString(email))
	if err != nil {
		log.Println("models.TaUser{}.GetTaUserByEmail checking err", err)
		if err == sql.ErrNoRows {
			handlers_utils.RespondWithError(writer, helpers.NewHttpError(http.StatusUnauthorized, constants.UserNotFoundOrPasswordMismatch))
			return
		}
		handlers_utils.RespondWithError(writer, helpers.NewHttpError(http.StatusInternalServerError, err.Error()))
		return
	}

	user := users[0]

	user.Update(handler.Query,
		types.NilString(),
		types.NilString(),
		types.NilString(),
		types.NilString(),
		types.NilString(),
		types.NewBool(true),
		types.NilString(),
	)
	log.Println(fmt.Sprintf("User %s was activated succesfully"), user.Email)

	// Redirect to the login page
	data := VerificationData{Name: user.Email.String()}
	writeHTMLWithTemplate(writer, data)
}

func writeHTMLWithTemplate(writer http.ResponseWriter, data VerificationData) {
	tmpl := template.Must(template.ParseFiles("index.html"))
	writer.Header().Set("Content-Type", "text/html")
	writer.WriteHeader(http.StatusOK)
	tmpl.Execute(writer, data)
}
