package auth

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"stories/database"
	"stories/generated/models"
	"stories/types"
)

type fakeQuery struct {
	sessionResult []models.TaSession
	userResult    models.TaUser
}

func (q *fakeQuery) QueryOne(query string, scanner database.Scanner, values ...interface{}) (interface{}, error) {
	if query == models.TaUserByPkUserQuery {
		return q.userResult, nil
	}
	return nil, nil
}

func (q *fakeQuery) QueryMany(query string, scanner database.Scanner, values ...interface{}) ([]interface{}, error) {
	if query == models.TaSessionByTokenQuery {
		rows := make([]interface{}, len(q.sessionResult))
		for i, session := range q.sessionResult {
			rows[i] = session
		}
		return rows, nil
	}
	return nil, nil
}

func (q *fakeQuery) QueryManyPaged(query string, scanner database.Scanner, page int, limit int, values ...interface{}) (types.Page, error) {
	return types.Page{}, nil
}

func (q *fakeQuery) Exec(query string, values ...interface{}) error {
	return nil
}

func (q *fakeQuery) Commit() error {
	return nil
}

func (q *fakeQuery) Rollback() {}

func (q *fakeQuery) Connect() error {
	return nil
}

func (q *fakeQuery) Initialize(newConfig *database.DatabaseConfig) bool {
	return true
}

func (q *fakeQuery) GetDatabase() *database.Connection {
	return nil
}

func (q *fakeQuery) SetSqlDB(db *sql.DB) {}

func (q *fakeQuery) NewQuery() database.IQuery {
	return q
}

func TestVerifyEmailHandler_ServeHTTP_Success(t *testing.T) {
	email := "user@example.com"
	code := "VALID_TOKEN"
	startedAt := types.DateTime{DateTime: time.Now().UTC().Add(-10 * time.Minute), NotNil: true}

	fake := &fakeQuery{
		sessionResult: []models.TaSession{{
			PkSession: types.NewInt(1),
			Token:     types.NewString(code),
			FkUser:    types.NewInt(42),
			StartedAt: startedAt,
			Invalid:   types.NewBool(false),
		}},
		userResult: models.TaUser{
			PkUser:      types.NewInt(42),
			Email:       types.NewString(email),
			Activated:   types.NewBool(false),
			Firstname:   types.NilString(),
			Lastname:    types.NilString(),
			Password:    types.NilString(),
			AccountType: types.NilString(),
			Username:    types.NilString(),
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/verify-email?code="+code, nil)
	rec := httptest.NewRecorder()

	handler := VerifyEmailHandler{Query: fake}
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	body := rec.Body.String()
	if !strings.Contains(body, "activated") {
		t.Fatalf("expected response body to include activated text, got %q", body)
	}
	if !strings.Contains(body, email) {
		t.Fatalf("expected response body to include user email, got %q", body)
	}
}

func TestVerifyEmailHandler_ServeHTTP_MissingCode(t *testing.T) {
	fake := &fakeQuery{}
	req := httptest.NewRequest(http.MethodGet, "/verify-email", nil)
	rec := httptest.NewRecorder()

	handler := VerifyEmailHandler{Query: fake}
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}

	if !strings.Contains(rec.Body.String(), "required field `code' cannot be empty") {
		t.Fatalf("expected error response about missing code, got %q", rec.Body.String())
	}
}
