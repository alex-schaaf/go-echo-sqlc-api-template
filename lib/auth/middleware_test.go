package auth

import (
	"app/lib/config"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
)

func TestCookieAuthMiddleware(t *testing.T) {
	e := echo.New()
	config := config.Config{JWT_SECRET: "jwtSecret"}

	testHandler := func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	}

	e.GET("/", testHandler, CookieAuthMiddleware(config.JWT_SECRET))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	cookie, err := GenerateTokenCookie(config.JWT_SECRET, 1)
	if err != nil {
		t.Fatal(err)
	}
	req.AddCookie(&cookie)
	rec := httptest.NewRecorder()
	e.NewContext(req, rec)

	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status code 200, got %d", rec.Code)
	}
}
