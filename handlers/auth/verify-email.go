package auth

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"stories/database"
	"stories/generated/models"
	"stories/types"

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
	var query url.Values
	query = request.URL.Query()
	code = query.Get("code")
	if len(code) == 0 {
		handlers_utils.RespondWithError(writer, helpers.NewHttpError(http.StatusBadRequest, fmt.Sprintf(constants.RequiredFieldIsEmpty, "code")))
		return
	}
	// Get the token from the query strin
	searchCode := types.NewString(code)
	taSessions, err := models.GetTaSessionByToken(handler.Query, searchCode)
	if err != nil {
		log.Println(fmt.Sprintf("Error getting token %s", code))
		handlers_utils.RespondWithError(writer, helpers.NewHttpError(http.StatusBadRequest, fmt.Sprintf(InvalidToken, code)))
		return

	}
	taSession := taSessions[0]
	taUser, err := models.GetTaUserByPkUser(handler.Query, taSession.FkUser)
	if err != nil {
		log.Println(fmt.Sprintf("Error getting user %d", taSession.FkUser))
		handlers_utils.RespondWithError(writer, helpers.NewHttpError(http.StatusBadRequest, fmt.Sprintf(InvalidToken, code)))
		return
	}

	err = taUser.Update(handler.Query,
		types.NilString(),
		types.NilString(),
		types.NilString(),
		types.NilString(),
		types.NilString(),
		types.NewBool(true),
	)
	if err != nil {
		log.Println(fmt.Sprintf("Error updating user %d", taSession.FkUser))
		handlers_utils.RespondWithError(writer, helpers.NewHttpError(http.StatusBadRequest, fmt.Sprintf(InvalidToken, code)))
		return
	}
	log.Println(fmt.Sprintf("User %s was activated succesfully"), taUser.Email)
	// Redirect to the login page
	data := VerificationData{Name: taUser.Email.String()}
	writeHTMLWithTemplate(writer, data)
}

func writeHTMLWithTemplate(writer http.ResponseWriter, data VerificationData) {
	tmpl := template.Must(template.ParseFiles("index.html"))
	writer.Header().Set("Content-Type", "text/html")
	writer.WriteHeader(http.StatusOK)
	tmpl.Execute(writer, data)
}
