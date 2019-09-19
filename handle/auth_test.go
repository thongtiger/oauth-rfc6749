package handle_test

import (
	"jwt-refresh-token/handle"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

const (
	userJSON        = `{"grant_type":"password","username":"joe","password":"password"}`
	userInvalidJSON = `{"grant_type":"fail","username":"joe","password":"password"}`
)

func TestLoginSuccess(t *testing.T) {
	t.Run("it should return httpCode 200", func(t *testing.T) {
		// Setup
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(userJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// mock
		h := handle.TokenHandle(c)
		// Assertions
		if assert.NoError(t, h) {
			assert.Equal(t, http.StatusOK, rec.Code)
			// assert.Equal(t, userJSON, rec.Body.String())
		}
	})

}
func TestLoginFail(t *testing.T) {
	t.Run("it should return httpCode 401", func(t *testing.T) {
		// Setup
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(userInvalidJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// mock
		h := handle.TokenHandle(c)
		// Assertions
		if assert.NoError(t, h) {
			assert.Equal(t, http.StatusUnauthorized, rec.Code)
			if status := rec.Code; status != http.StatusUnauthorized {
				t.Errorf("wrong code: got %v want %v", status, http.StatusOK)
			}

		}
	})
}
