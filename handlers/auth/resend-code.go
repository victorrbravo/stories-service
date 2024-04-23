package auth

import (
	"fmt"
	"log"
	"net/http"
	"stories/database"
	"stories/generated/models"
	"stories/helpers"
	"stories/pkg/utils"
	"stories/types"

	"stories/handlers_utils"
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
func (handler ResendCodeHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request, user ModelUser) {

	err := handlers_utils.GetParameters(&handler, request)
	if err != nil {
		log.Println("ResendCodeHandler ServeHTTP GetParameters error", err)
		handlers_utils.RespondWithError(writer, helpers.NewHttpError(http.StatusBadRequest, err.Error()))
		return
	}
	taSessions, err := models.GetTaSessionByToken(handler.Query, handler.Token)
	if err != nil {
		log.Println(fmt.Sprintf("Error getting token %s", handler.Token))
		handlers_utils.RespondWithError(writer, helpers.NewHttpError(http.StatusBadRequest, fmt.Sprintf(InvalidToken, handler.Token)))
		return
	}
	taSession := taSessions[0]

	code := utils.RandomString(8)

	err = taSession.Update(handler.Query,
		types.NewString(code),
		types.NilInt,
		types.Now(),
		false,
	)
	if err != nil {
		log.Println(fmt.Sprintf("Error updating token %s", handler.Token))
		handlers_utils.RespondWithError(writer, helpers.NewHttpError(http.StatusBadRequest, fmt.Sprintf(InvalidToken, handler.Token)))
		return
	}
	body := ""
	taUser, err := models.GetTaUserByPkUser(handler.Query, taSession.FkUser)

	if err != nil {
		log.Println(fmt.Sprintf("Error getting user %d", taSession.FkUser))
		handlers_utils.RespondWithError(writer, helpers.NewHttpError(http.StatusBadRequest, fmt.Sprintf(InvalidToken, handler.Token)))
		return
	}
	SendEmail(SenderSignupEmail,
		taUser.Email.String(),
		DefaultSubjectSignUpEmail,
		body)

	handlers_utils.SendJSONResponse(writer, user, http.StatusOK)
}
