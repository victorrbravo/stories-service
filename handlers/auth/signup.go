package auth

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"stories/database"
	"stories/generated/models"
	"stories/handlers_utils"
	"stories/helpers"
	"stories/pkg/utils"
	"stories/types"
)

type SignupUserHandler struct {
	Email                string       `json:"email" validate:"required,email"`
	Password             types.String `json:"password" validate:"required"`
	PasswordConfirmation types.String `json:"password_confirmation"  validate:"required"`
	Name                 types.String `json:"name"  validate:"required"`
	LastName             types.String `json:"last_name"  validate:"required"`
	Phone                string       `json:"phone" validate:"required,e164"`
	Gender               types.Gender `json:"gender"`
	Username             types.String `json:"username"`
	Query                database.IQuery
}

func (handler SignupUserHandler) IsValid() error {
	return nil
}

// Check if endpoint is enabled
func (handler SignupUserHandler) IsEnabled() bool {
	return true
}

func (handler SignupUserHandler) GetQuery() database.IQuery {
	log.Println("SetQuery was executed")
	return handler.Query
}

func sendEmail(smtpHost, smtpPort, smtpUsername, smtpPassword, from, to, subject, body string) error {
	// Compose the email message.
	message := fmt.Sprintf("From: %s\r\n", from)
	message += fmt.Sprintf("To: %s\r\n", to)
	message += fmt.Sprintf("Subject: %s\r\n", subject)
	message += "\r\n" + body

	// Connect to the SMTP server.
	auth := smtp.PlainAuth("", smtpUsername, smtpPassword, smtpHost)
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, []byte(message))
	if err != nil {
		return err
	}

	return nil
}

// Get currently authenticated user's information
func (handler SignupUserHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if handler.Query.GetDatabase() == nil {
		handlers_utils.RespondWithError(writer, helpers.NewHttpError(http.StatusInternalServerError, "Database not initialized"))
		return
	}
	err := handlers_utils.GetParameters(&handler, request)
	if err != nil {
		log.Println("SignupUserHandler ServeHTTP GetParameters error", err)
		handlers_utils.RespondWithError(writer, helpers.NewHttpError(http.StatusBadRequest, err.Error()))
		return
	}
	validate := validator.New(validator.WithRequiredStructEnabled())
	err = validate.Struct(handler)
	if err != nil {
		handlers_utils.RespondWithError(writer, helpers.NewHttpError(http.StatusBadRequest, err.Error()))
		return
	}
	if handler.Password != handler.PasswordConfirmation {
		handlers_utils.RespondWithError(writer, helpers.NewHttpError(http.StatusBadRequest, "Passwords do not match"))
		return
	}
	passcrypt, _ := bcrypt.GenerateFromPassword([]byte(types.AsString(handler.Password)), 0)
	password := string(passcrypt)
	code := utils.RandomString(8)
	log.Println("code: |", code, "|")
	taUser, err := models.CreateTaUser(handler.Query,
		handler.Name,
		handler.LastName,
		types.NewString(handler.Email),
		types.NewString(password),
		types.NewString("user"),
		types.NewBool(false),
		handler.Username,
	)
	if err != nil {
		handlers_utils.RespondWithError(writer, helpers.NewHttpError(http.StatusInternalServerError, err.Error()))
		return
	}
	taSession, err := models.CreateTaSession(handler.Query,
		types.NewString(code),
		taUser.PkUser,
		types.Now(),
		false,
	)
	if err != nil {
		log.Println("Error creating session", err.Error())
		handlers_utils.RespondWithError(writer, helpers.NewHttpError(http.StatusInternalServerError, err.Error()))
		return
	}
	log.Println("taSession pkSession", taSession.PkSession)
	urlToSend := buildUrlToSend(handler.Email, code)
	log.Println("urlToSend", urlToSend)
	log.Println("created taUser ID:", taUser.PkUser)
	// Send to Queue
	// Send email to user
	body := fmt.Sprintf(DefaultSignUpBodyEmail, handler.Name, handler.LastName, urlToSend)
	log.Println("SES email body: |", body, "|")
	log.Println("SES email subject: |", DefaultSubjectSignUpEmail, "|")
	err = VerifyEmailAddress(handler.Email, "us-east-1")
	if err != nil {
		log.Println("Error trying to send  email for verification", err.Error())
		handlers_utils.RespondWithError(writer, helpers.NewHttpError(http.StatusInternalServerError, err.Error()))
		return
	}
	log.Println("After Email sent successfully!: ", handler.Email)
	handlers_utils.SendJSONResponse(writer, struct {
		Message string `json:"message"`
		Email   string `json:"email"`
	}{Message: fmt.Sprintf("User %s %s created successfully", handler.Name, handler.LastName),
		Email: handler.Email}, http.StatusOK)

}

func buildUrlToSend(email string, code string) string {
	urlVerify := os.Getenv("URL_VERIFY_EMAIL")
	if len(urlVerify) == 0 {
		urlVerify = DefaultURLVerifyEmail
	}
	endpointVerify := fmt.Sprintf(EndPointVerifyEmail, urlVerify, code)
	return endpointVerify

}
