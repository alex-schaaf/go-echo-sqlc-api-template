package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
)

func TestCookieAuthMiddleware(t *testing.T) {
	e := echo.New()

	testHandler := func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	}

	e.GET("/", testHandler, CookieAuthMiddleware)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	cookie, err := GenerateTokenCookie(1)
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
