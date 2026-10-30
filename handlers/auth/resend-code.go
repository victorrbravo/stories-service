package auth

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"github.com/victorrbravo/stories-service/constants"
	"github.com/victorrbravo/stories-service/database"
	"github.com/victorrbravo/stories-service/generated/models"
	"github.com/victorrbravo/stories-service/helpers"
	"github.com/victorrbravo/stories-service/pkg/utils"
	"github.com/victorrbravo/stories-service/types"

	"github.com/victorrbravo/stories-service/handlers_utils"
)

type ResendCodeHandler struct {
	Token types.String `json:"token"`
	Query database.IQuery
}

func (handler ResendCodeHandler) IsValid() error {
	return nil
}

func (handler *ResendCodeHandler) SetQuery(query database.IQuery) {
	handler.Query = query
}

// Check if endpoint is enabled
func (handler ResendCodeHandler) IsEnabled() bool {
	return true
}

// Get currently authenticated user's information
func (handler ResendCodeHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
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
		log.Println(fmt.Sprintf("Error getting token %s", handler.Token))
		handlers_utils.RespondWithError(writer, helpers.NewHttpError(http.StatusBadRequest, fmt.Sprintf(InvalidToken, handler.Token)))
		return
	}
	taSession := taSessions[0]

	newCode := utils.RandomString(8)

	pkUser := taSession.FkUser
	err = taSession.Update(handler.Query,
		types.NewString(newCode),
		types.NilInt,
		types.Now(),
		types.NewBool(false),
	)
	if err != nil {
		log.Println(fmt.Sprintf("Error updating token %s", handler.Token))
		handlers_utils.RespondWithError(writer, helpers.NewHttpError(http.StatusBadRequest, fmt.Sprintf(InvalidToken, handler.Token)))
		return
	}
	log.Println("fkUser", pkUser)
	taUser, err := models.GetTaUserByPkUser(handler.Query, pkUser)
	if err != nil {
		log.Println(fmt.Sprintf("Error getting user %v", taSession.FkUser), "error", err.Error())
		handlers_utils.ResponseWithInternalError(writer)
		return
	}
	email := taUser.Email.String()
	body := fmt.Sprintf(`
			Please click <a href="%s"><b>here</b></a> to verify your email.</br>
		Please use this code to verify your email address on Use2 Team.
`, buildUrlToSend(email, newCode))

	err = ResendSendEmail([]string{taUser.Email.String()}, DefaultResendSubjectEmail, body)

	if err != nil {
		log.Println("Error sending email", err.Error())
		handlers_utils.RespondWithError(writer, helpers.NewHttpError(http.StatusInternalServerError, err.Error()))
		return
	}

	data := VerificationData{Name: email, TemplateContent: CodewasResendTemplate}
	err = writeHTMLWithTemplate(writer, data)
}
