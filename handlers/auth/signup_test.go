package auth

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"github.com/victorrbravo/stories-service/database"
	"github.com/victorrbravo/stories-service/internal/testenv"
	db "github.com/victorrbravo/stories-service/database"
	"strings"
	"testing"
)

// Credentials are read from the repository .env file (see .env.example)
var (
	DB_LOCAL_HOST string
	DB_LOCAL_PASS string
	DB_NAME       string
	DB_USER       string
)

func init() {
	testenv.Load()
	DB_LOCAL_HOST = testenv.Get("DB_HOST", "localhost")
	DB_LOCAL_PASS = testenv.Get("DB_PASS", "stories")
	DB_NAME = testenv.Get("DB_NAME", "stories")
	DB_USER = testenv.Get("DB_USER", "stories")
}

func TestSignupUserHandler_ServeHTTP_DatabaseNotInitialized(t *testing.T) {
	os.Setenv("URL_VERIFY_EMAIL", "https://doitall.work/confirm/")
	config := database.DatabaseConfig{
		Host: DB_LOCAL_HOST, Name: DB_NAME, Port: 5432, User: DB_USER,
		Password: DB_LOCAL_PASS,
	}
	var pgQuery db.PostgresQuery
	pgQuery.Initialize(&config)
	err := pgQuery.Connect()
	if err != nil {
		t.Error(err)
	}
	query := pgQuery.NewQuery()

	payload := `{"email":"victorrbravo7@gmail.com","password":"secret","password_confirmation":"secret","name":"Jane","last_name":"Doe","phone":"+1234567890","gender":"female","username":"janedoe"}`
	req := httptest.NewRequest(http.MethodPost, "/signup", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler := SignupUserHandler{}
	handler.Query = query
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, rec.Code)
	}
}

func TestSignupUserHandler_IsEnabledAndIsValid(t *testing.T) {
	handler := SignupUserHandler{}

	if !handler.IsEnabled() {
		t.Fatal("expected IsEnabled to return true")
	}

	if err := handler.IsValid(); err != nil {
		t.Fatalf("expected IsValid to return nil, got %v", err)
	}
}

func TestStringSlice_SetAndString(t *testing.T) {
	var s stringSlice

	if err := s.Set("one"); err != nil {
		t.Fatalf("unexpected error from Set: %v", err)
	}
	if err := s.Set("two"); err != nil {
		t.Fatalf("unexpected error from Set: %v", err)
	}

	if got := s.String(); got != "one two" {
		t.Fatalf("expected %q, got %q", "one two", got)
	}
}

func TestGetSMTPConfiguration_MissingVariables(t *testing.T) {
	t.Setenv("SMTP_USERNAME", "")
	t.Setenv("SMTP_PASSWORD", "")
	t.Setenv("SMTP_HOST", "")
	t.Setenv("SMTP_PORT", "")

	_, err := getSMTPConfiguration()
	if err == nil {
		t.Fatal("expected error when SMTP env vars are missing")
	}
	if !strings.Contains(err.Error(), "SMTP_USERNAME not set") {
		t.Fatalf("expected error to mention missing SMTP_USERNAME, got %q", err.Error())
	}
}

func TestBuildUrlToSend_DefaultsToPackageDefault(t *testing.T) {
	t.Setenv("URL_VERIFY_EMAIL", "")
	code := "CODE123"

	got := buildUrlToSend("test@example.com", code)
	expected := fmt.Sprintf(EndPointVerifyEmail, DefaultURLVerifyEmail, code)

	if got != expected {
		t.Fatalf("expected %q, got %q", expected, got)
	}
}

func TestBuildUrlToChangePassword_UsesEnvValue(t *testing.T) {
	envURL := "https://example.com/change-password"
	t.Setenv("URL_CHANGE_PASSWORD", envURL)
	code := "RESET123"

	got := buildUrlToChangePassword("test@example.com", code)
	expected := fmt.Sprintf(EndPointVerifyEmail, envURL, code)

	if got != expected {
		t.Fatalf("expected %q, got %q", expected, got)
	}
}
