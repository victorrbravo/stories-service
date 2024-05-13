package auth

import (
	"database/sql"
	"errors"
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

type SendChangePwdHandler struct {
	Email types.String `json:"email"`
	Query database.IQuery
}

func (handler SendChangePwdHandler) IsValid() error {
	return nil
}

func (handler *SendChangePwdHandler) SetQuery(query database.IQuery) {
	handler.Query = query
}

// Check if endpoint is enabled
func (handler SendChangePwdHandler) IsEnabled() bool {
	return true
}

// Get currently authenticated user's information
func (handler SendChangePwdHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	err := handlers_utils.GetParameters(&handler, request)
	log.Println("SendChangePwdHandler  err", err)
	if err != nil {
		handlers_utils.RespondWithError(writer, helpers.NewHttpError(http.StatusBadRequest, err.Error()))
		return
	}
	taUsers, err := models.GetTaUserByEmail(handler.Query, handler.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Println(fmt.Sprintf("User with email %v not found", handler.Email))
			handlers_utils.RespondWithError(writer, helpers.NewHttpError(http.StatusNotFound, fmt.Sprintf("User with email %v not found", handler.Email)))
			return
		}

		log.Println(fmt.Sprintf("Error getting user email %v", handler.Email), "error", err.Error())
		handlers_utils.ResponseWithInternalError(writer)
		return
	}
	taUser := taUsers[0]
	if !taUser.Activated.ToBool() {
		log.Println("Change Password User not activated yet email:", handler.Email)
		handlers_utils.RespondWithError(
			writer,
			helpers.NewHttpError(http.StatusUnauthorized, UserNotActivated),
		)
		return
	}
	newCode := utils.RandomString(6)
	email := handler.Email.String()
	body := fmt.Sprintf(`
			Please write this code <div style="color: black; font-family: Arial; font-size: 16px;">%s</div> in Use2 App, <br/>  in order to change your current password.</br>
		Use2 app. .
`, newCode)
	smtpConf, err := getSMTPConfiguration()
	if err != nil {
		log.Println("Error trying to get SMTP configuration", err.Error())
		handlers_utils.ResponseWithInternalError(writer)
		return
	}
	_, err = models.CreateTaSession(handler.Query,
		types.NewString(newCode),
		taUser.PkUser,
		types.Now(),
		types.NewBool(false),
	)
	if err != nil {
		log.Println("Error creating session", err.Error())
		handlers_utils.RespondWithError(writer, helpers.NewHttpError(http.StatusInternalServerError, err.Error()))
		return
	}
	err = sendEmailTLS(smtpConf.SmtpUsername,
		smtpConf.SmtpPassword,
		smtpConf.SmtpHost,
		smtpConf.SmtpPort,
		email,
		DefaultChangePassword,
		body)

	if err != nil {
		log.Println("ChangePassword Error sending email", err.Error())
		handlers_utils.ResponseWithInternalError(writer)
		return
	}
	handlers_utils.SendJSONResponse(writer, struct {
		Status  int    `json:"Status"`
		Message string `json:"Message"`
	}{
		Status:  http.StatusOK,
		Message: "Code sent successfully",
	}, http.StatusOK)

}
