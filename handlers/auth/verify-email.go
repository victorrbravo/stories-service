package auth

import (
	"fmt"
	"html/template"
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

type VerificationData struct {
	Name            string `json:"name"`
	TemplateContent string `json:"template_name"`
	LinkResend      string `json:"link_resend"`
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
	var code string
	var query url.Values
	query = request.URL.Query()
	code = query.Get("code")
	if len(code) == 0 {
		log.Println("code cannot be empty")
		handlers_utils.RespondWithError(writer, helpers.NewHttpError(http.StatusBadRequest, fmt.Sprintf(constants.RequiredFieldIsEmpty, "code")))
		return
	}
	// Get the token from the query string
	log.Println("code is not empty")
	searchCode := types.NewString(code)
	taSessions, err := models.GetTaSessionByToken(handler.Query, searchCode)
	if err != nil {
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
		urlResend := os.Getenv("URL_RESEND_EMAIL")
		if len(urlResend) == 0 {
			log.Println("URL_RESEND_EMAIL is not set")
			handlers_utils.ResponseWithInternalError(writer)
			return
		}
		log.Printf("Token %s has expired", code)
		linkResend := fmt.Sprintf("%s?code=%s", urlResend, code)
		nameUser := fmt.Sprintf("%s", taUser.Email)
		writeHTMLWithTemplate(writer, VerificationData{Name: nameUser,
			TemplateContent: ExpiredTokenVerifyEmail,
			LinkResend:      linkResend,
		})
		return
	}
	log.Println("Session found: ", taSession.PkSession)

	if taUser.Activated.ToBool() {
		log.Printf("User %s was already activated", taUser.Email)
		writeHTMLWithTemplate(writer, VerificationData{Name: taUser.Email.String(), TemplateContent: AlreadyVerifiedEmail})
		return

	}
	log.Println("User found: ", taUser.Email, "PkUser", taUser.PkUser)
	email := taUser.Email
	err = taUser.Update(handler.Query,
		types.NilString(),
		types.NilString(),
		types.NilString(),
		types.NilString(),
		types.NilString(),
		types.NewBool(true),
		types.NilString(),
	)
	if err != nil {
		log.Printf("Error updating user %v", taSession.FkUser)
		handlers_utils.RespondWithError(writer, helpers.NewHttpError(http.StatusBadRequest, fmt.Sprintf(InvalidToken, code)))
		return
	}
	log.Printf("User %s was activated succesfully", email)
	// Redirect to the login page
	log.Println("VerificationData taUser.Email: |", email, "|")
	data := VerificationData{Name: email.String(), TemplateContent: SuccessVerifyEmail}
	err = writeHTMLWithTemplate(writer, data)
	if err != nil {
		log.Println("Error writing HTML with template error", err.Error())
		handlers_utils.ResponseWithInternalError(writer)
		return
	}
}

func writeHTMLWithTemplate(writer http.ResponseWriter, data VerificationData) error {
	parser := template.New("verify-email")
	tmplObj, err := parser.Parse(data.TemplateContent)
	if err != nil {
		return err
	}
	writer.Header().Set("Content-Type", "text/html")
	writer.WriteHeader(http.StatusOK)
	tmpl := template.Must(tmplObj, err)
	err = tmpl.Execute(writer, data)
	log.Println("writeHTMLWithTemplate err", err)
	return nil
}
