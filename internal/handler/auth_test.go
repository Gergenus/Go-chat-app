package handler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Gergenus/StandardLib/internal/models"
	"github.com/Gergenus/StandardLib/internal/service"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestEchoHandlerAuth_SignUp(t *testing.T) {
	type mockBeahaviour func(s *service.AuthMock, user models.User)

	testTable := []struct {
		name                 string
		inputBody            string
		inputUser            models.User
		mock                 mockBeahaviour
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "OK",
			inputBody: `{"username": "Denis", "email": "123", "password": "123"}`,
			inputUser: models.User{
				Username: "Denis",
				Email:    "123",
				Password: "123",
			},
			mock: func(s *service.AuthMock, user models.User) {
				s.EXPECT().SignUp(user.Username, user.Email, user.Password).Return(1, nil)
			},
			expectedStatusCode:   http.StatusCreated,
			expectedResponseBody: `{"id": 1}`,
		},
		{
			name:                 "Empty fields",
			inputBody:            `{"email": "123", "password": "123"}`,
			mock:                 func(s *service.AuthMock, user models.User) {},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"error": "Invalid request payload"}`,
		},
	}

	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {

			auth := service.NewAuthMock(t)

			test.mock(auth, test.inputUser)

			hndler := NewEchoHandlerAuth(auth)

			r := echo.New()
			r.POST("/SignUp", hndler.SignUp)
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/SignUp", bytes.NewBufferString(test.inputBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			r.ServeHTTP(w, req)

			assert.JSONEq(t, test.expectedResponseBody, w.Body.String())
			assert.Equal(t, test.expectedStatusCode, w.Code)
		})
	}
}
