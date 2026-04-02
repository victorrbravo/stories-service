package auth

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"stories/database"
	"stories/generated/models"
	"stories/handlers_utils"
	"stories/helpers"
	"stories/pkg/utils"
	"stories/types"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/resend/resend-go/v3"
	"golang.org/x/crypto/bcrypt"
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

type SMTPConfiguration struct {
	SmtpUsername string
	SmtpPassword string
	SmtpHost     string
	SmtpPort     string
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

type stringSlice []string

func (s *stringSlice) String() string {
	return strings.Join(*s, " ")
}

func (s *stringSlice) Set(value string) error {
	*s = append(*s, value)
	return nil
}

func resendSendEmail() {
	apiKey := os.Getenv("RESEND_API_KEY")
	client := resend.NewClient(apiKey)
	var toMails stringSlice
	flag.Var(&toMails, "m", "Add a mail")
	flag.Parse()
	params := &resend.SendEmailRequest{
		From:    "support@doitall.work",
		To:      toMails,
		Subject: "Bienvenido a Doitall",
		Html:    "<h1>Bienvenido a Doitall, 🔥</h1>",
	}
	client.Emails.Send(params)
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
		types.NewBool(false),
	)
	if err != nil {
		log.Println("Error creating session", err.Error())
		handlers_utils.RespondWithError(writer, helpers.NewHttpError(http.StatusInternalServerError, err.Error()))
		return
	}
	log.Println("taSession pkSession", taSession.PkSession)
	urlToSend := buildUrlToSend(handler.Email, code)
	log.Println("urlToSend", urlToSend, "created taUser ID:", taUser.PkUser)

	body := fmt.Sprintf(DefaultSignUpBodyEmail, handler.Name, handler.LastName, urlToSend)
	log.Println("body", body)
	//err = SendEmail()

	if err != nil {
		log.Println("Error sending email", err.Error())
		handlers_utils.RespondWithError(writer, helpers.NewHttpError(http.StatusInternalServerError, err.Error()))
		return
	}
	log.Println("After Email sent successfully!: ", handler.Email)
	handlers_utils.SendJSONResponse(writer, struct {
		Status  int    `json:"Status"`
		Message string `json:"Message"`
		Email   string `json:"Email"`
	}{Status: 200,
		Message: fmt.Sprintf("User %s %s created successfully", handler.Name, handler.LastName),
		Email:   handler.Email}, http.StatusOK)

}

func getSMTPConfiguration() (SMTPConfiguration, error) {
	smtpConf := SMTPConfiguration{}

	smtpConf.SmtpUsername = os.Getenv("SMTP_USERNAME")
	smtpConf.SmtpPassword = os.Getenv("SMTP_PASSWORD")
	smtpConf.SmtpHost = os.Getenv("SMTP_HOST")
	smtpConf.SmtpPort = os.Getenv("SMTP_PORT")
	if len(smtpConf.SmtpUsername) == 0 {
		return smtpConf, fmt.Errorf("SMTP_USERNAME not set")
	}
	if len(smtpConf.SmtpPassword) == 0 {
		return smtpConf, fmt.Errorf("SMTP_PASSWORD not set")
	}
	if len(smtpConf.SmtpHost) == 0 {
		return smtpConf, fmt.Errorf("SMTP_HOST not set")
	}
	if len(smtpConf.SmtpPort) == 0 {
		return smtpConf, fmt.Errorf("SMTP_PORT not set")
	}
	return smtpConf, nil
}

func buildUrlToSend(email string, code string) string {
	urlVerify := os.Getenv("URL_VERIFY_EMAIL")
	if len(urlVerify) == 0 {
		urlVerify = DefaultURLVerifyEmail
	}
	endpointVerify := fmt.Sprintf(EndPointVerifyEmail, urlVerify, code)
	return endpointVerify
}

func buildUrlToChangePassword(email string, code string) string {
	urlVerify := os.Getenv("URL_CHANGE_PASSWORD")
	if len(urlVerify) == 0 {
		urlVerify = DefaultURLVerifyEmail
	}
	endpointVerify := fmt.Sprintf(EndPointVerifyEmail, urlVerify, code)
	return endpointVerify

}
