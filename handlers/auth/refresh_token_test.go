package auth

import (
	"net/http"
	"net/http/httptest"
	"os"
	"stories/constants"
	"stories/helpers"
	"strings"
	"testing"
)

func TestRefreshTokenHandler_ServeHTTP(t *testing.T) {
	// Setup env
	os.Setenv(constants.JWTAccessEnvironmentOSName, "test_access_secret")
	os.Setenv(constants.JWTRefreshEnvironmentOSName, "test_refresh_secret")
	defer os.Unsetenv(constants.JWTAccessEnvironmentOSName)
	defer os.Unsetenv(constants.JWTRefreshEnvironmentOSName)

	email := "test@example.com"
	td, _ := helpers.CreateJWTToken(email)

	payload := `{"refresh_token":"` + td.RefreshToken + `"}`
	req := httptest.NewRequest(http.MethodPost, "/refresh", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler := RefreshTokenHandler{}
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
}

func TestRefreshTokenHandler_ServeHTTP_InvalidToken(t *testing.T) {
	os.Setenv(constants.JWTRefreshEnvironmentOSName, "test_refresh_secret")
	defer os.Unsetenv(constants.JWTRefreshEnvironmentOSName)

	payload := `{"refresh_token":"invalid_token"}`
	req := httptest.NewRequest(http.MethodPost, "/refresh", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler := RefreshTokenHandler{}
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, rec.Code)
	}
}

func TestRefreshTokenHandler_ServeHTTP_EmptyToken(t *testing.T) {
	payload := `{"refresh_token":""}`
	req := httptest.NewRequest(http.MethodPost, "/refresh", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler := RefreshTokenHandler{}
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}
