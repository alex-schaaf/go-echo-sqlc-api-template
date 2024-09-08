package auth

import (
	"app/db"
	"app/lib/auth"
	"app/lib/config"
	"app/lib/test"
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	_ "github.com/mattn/go-sqlite3"
)

func TestAuthHandler_RegisterHandler(t *testing.T) {
	e, queries := test.SetupTest()
	config := &config.Config{}

	authHandler := NewAuthHandler(queries, config)
	registerJSON := `{"username":"username","email":"user@example.com","password":"password"}`
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(registerJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if err := authHandler.RegisterHandler(c); err != nil {
		t.Fatal(err)
	}
	if rec.Code != http.StatusCreated {
		t.Errorf("Expected status code %d but got %d", http.StatusCreated, rec.Code)
	}
	if rec.Body.String() != "" {
		t.Errorf("Expected empty body but got %s", rec.Body.String())
	}
}

func TestAuthHandler_LoginHandler(t *testing.T) {
	e, queries := test.SetupTest()
	ctx := context.Background()
	config := &config.Config{}

	authHandler := NewAuthHandler(queries, config)
	password := "password"
	passwordHash, _ := auth.HashPassword(password)
	_, err := queries.CreateUser(ctx, db.CreateUserParams{
		Username:     "username",
		Email:        "user@example.com",
		PasswordHash: passwordHash,
	})
	if err != nil {
		t.Fatal(err)
	}

	loginJSON := `{"email":"user@example.com","password":"password"}`
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(loginJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if err := authHandler.LoginHandler(c); err != nil {
		t.Fatal(err)
	}
	if rec.Code != http.StatusOK {
		t.Errorf("Expected status code %d but got %d", http.StatusOK, rec.Code)
	}
	// check if cookie is set
	cookie := rec.Result().Cookies()[0]
	if cookie.Name != "token" {
		t.Errorf("Expected cookie name %s but got %s", "token", cookie.Name)
	}
	if cookie.HttpOnly != true {
		t.Errorf("Expected cookie to be HttpOnly but it is not")
	}
	if cookie.Value == "" {
		t.Errorf("Expected cookie value to be set but it is empty")
	}
	if cookie.Expires.IsZero() {
		t.Errorf("Expected cookie to have an expiry time but it does not")
	}
}

func TestAuthHandler_LogoutHandler(t *testing.T) {
	e, queries := test.SetupTest()
	config := &config.Config{}

	authHandler := NewAuthHandler(queries, config)
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if err := authHandler.LogoutHandler(c); err != nil {
		t.Fatal(err)
	}
	if rec.Code != http.StatusNoContent {
		t.Errorf("Expected status code %d but got %d", http.StatusNoContent, rec.Code)
	}
	// check if cookie is set
	cookie := rec.Result().Cookies()[0]
	if cookie.Name != "token" {
		t.Errorf("Expected cookie name %s but got %s", "token", cookie.Name)
	}
	if cookie.HttpOnly != true {
		t.Errorf("Expected cookie to be HttpOnly but it is not")
	}
	if cookie.Value != "" {
		t.Errorf("Expected cookie value to be empty but it is not")
	}
	if !cookie.Expires.Equal(time.Unix(0, 0)) {
		t.Errorf("Expected cookie to have no expiry time but it does")
	}
}
