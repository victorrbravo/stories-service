package auth

import (
	"crypto/tls"
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

func sendEmailTLS(from string, password string, smtpHost string, smtpPort, to string, subject string, body string) error {
	// Sender email configuration

	// Authentication setup
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// SMTP client with TLS configuration
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true, // Skip certificate verification, not recommended in production
		ServerName:         smtpHost,
	}
	conn, err := tls.Dial("tcp", smtpHost+":"+smtpPort, tlsConfig)
	if err != nil {
		return err
	}
	defer conn.Close()

	// Sending email
	client, err := smtp.NewClient(conn, smtpHost)
	if err != nil {
		return err
	}

	if err := client.Auth(auth); err != nil {
		return err
	}

	if err := client.Mail(from); err != nil {
		return err
	}

	if err := client.Rcpt(to); err != nil {
		return err
	}

	w, err := client.Data()
	if err != nil {
		return err
	}
	message := fmt.Sprintf("Subject: %s\nContent-Type: text/html; charset=UTF-8\n\n%s", subject, body)
	bytesBody := []byte(message)
	_, err = w.Write(bytesBody)
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	log.Println("Email sent successfully!")
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
	log.Println("urlToSend", urlToSend, "created taUser ID:", taUser.PkUser)

	smtpConf, err := getSMTPConfiguration()
	if err != nil {
		log.Println("Error trying to get SMTP configuration", err.Error())
		handlers_utils.ResponseWithInternalError(writer)
		return
	}
	body := fmt.Sprintf(DefaultSignUpBodyEmail, handler.Name, handler.LastName, urlToSend)
	err = sendEmailTLS(smtpConf.SmtpUsername,
		smtpConf.SmtpPassword,
		smtpConf.SmtpHost,
		smtpConf.SmtpPort,
		handler.Email,
		DefaultSubjectSignUpEmail,
		body)

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
