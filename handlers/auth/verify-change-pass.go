package auth

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"stories/database"
	"stories/generated/models"
	"stories/types"
	"time"

	"stories/constants"
	"stories/handlers_utils"
	"stories/helpers"
)

type VerifyChangePassHandler struct {
	Query database.IQuery
}

func (handler VerifyChangePassHandler) IsValid() error {
	return nil
}

// Check if endpoint is enabled
func (handler VerifyChangePassHandler) IsEnabled() bool {
	return true
}

// Verify the code sent by email to change the user's password.
// Works the same way as VerifyEmailHandler: the code must exist and must not be
// older than 15 minutes. Instead of activating the user, it sends the user to
// the change password page, and when that page cannot be loaded it renders a
// "can't load change pass" message.
// You must add a parameter to Query String (URL) in this format: ?code=<CODE_SENT>
func (handler VerifyChangePassHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	var err error
	var code string
	var query url.Values
	query = request.URL.Query()
	code = query.Get("code")
	if len(code) == 0 {
		log.Println("code cannot be empty")
		handlers_utils.RespondWithError(writer, helpers.NewHttpError(http.StatusBadRequest, fmt.Sprintf(constants.RequiredFieldIsEmpty, "code")))
		return
	}
	log.Printf("code is not empty: %v", code)
	searchCode := types.NewString(code)
	taSessions, err := models.GetTaSessionByToken(handler.Query, searchCode)
	if err != nil || len(taSessions) == 0 {
		log.Printf("Error getting token %s", code)
		handlers_utils.RespondWithError(writer, helpers.NewHttpError(http.StatusBadRequest, fmt.Sprintf(InvalidToken, code)))
		return
	}
	taSession := taSessions[0]
	log.Println("Started At", taSession.StartedAt)
	limit := taSession.StartedAt.Add(time.Minute * 15)
	log.Println("Limit", limit)
	now := types.NowUTC()
	taUser, err := models.GetTaUserByPkUser(handler.Query, taSession.FkUser)
	if err != nil {
		log.Printf("Error getting user %v", taSession.FkUser)
		handlers_utils.RespondWithError(writer, helpers.NewHttpError(http.StatusBadRequest, fmt.Sprintf(InvalidToken, code)))
		return
	}
	if now.After(limit) {
		log.Printf("Token %s has expired", code)
		err = writeHTMLWithTemplate(writer, VerificationData{
			Name:            taUser.Email.String(),
			TemplateContent: ExpiredTokenChangePass,
		})
		if err != nil {
			log.Println("Error writing HTML with template error", err.Error())
			handlers_utils.ResponseWithInternalError(writer)
		}
		return
	}
	log.Println("Session found: ", taSession.PkSession)
	if !taUser.Activated.ToBool() {
		log.Println("Change Password User not activated yet email:", taUser.Email)
		handlers_utils.RespondWithError(writer, helpers.NewHttpError(http.StatusUnauthorized, UserNotActivated))
		return
	}

	urlChangePassword := os.Getenv("URL_CHANGE_PASSWORD")
	if len(urlChangePassword) == 0 {
		log.Println("URL_CHANGE_PASSWORD is not set, cannot load the change password page")
		err = writeHTMLWithTemplate(writer, VerificationData{
			Name:            taUser.Email.String(),
			TemplateContent: CannotLoadChangePass,
		})
		if err != nil {
			log.Println("Error writing HTML with template error", err.Error())
			handlers_utils.ResponseWithInternalError(writer)
		}
		return
	}
	linkChangePassword := fmt.Sprintf(EndPointVerifyEmail, urlChangePassword, code)
	log.Printf("Code %s verified for %s, redirecting to %s", code, taUser.Email, linkChangePassword)
	http.Redirect(writer, request, linkChangePassword, http.StatusSeeOther)
}
