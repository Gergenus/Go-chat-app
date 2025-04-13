package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/Gergenus/StandardLib/internal/models"
	"github.com/Gergenus/StandardLib/internal/service"
	"github.com/labstack/echo/v4"
)

type HandlerAuth interface {
	SignUp(c echo.Context) error
	SignIn(c echo.Context) error
}

type EchoHandlerAuth struct {
	ServiceAuth service.Auth
}

func NewEchoHandlerAuth(auth service.Auth) EchoHandlerAuth {
	return EchoHandlerAuth{ServiceAuth: auth}
}

func (e *EchoHandlerAuth) SignUp(c echo.Context) error {
	var user models.User
	err := c.Bind(&user)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, map[string]string{
			"error": "Invalid request payload",
		})
	}
	id, err := e.ServiceAuth.SignUp(user.Username, user.Email, user.Password)
	if err != nil {
		if errors.Is(err, service.ErrUserAlreadyExists) {
			return echo.NewHTTPError(http.StatusBadRequest, map[string]string{
				"error": "User already exists",
			})
		}
		return echo.NewHTTPError(http.StatusInternalServerError, map[string]string{
			"error": "Internal Server Error",
		})
	}
	return c.JSON(http.StatusCreated, map[string]interface{}{
		"id": id,
	})
}

func (e *EchoHandlerAuth) SignIn(c echo.Context) error {
	var user models.User
	err := c.Bind(&user)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, map[string]string{
			"error": "Invalid request payload",
		})
	}
	token, err := e.ServiceAuth.SignIn(user.Username, user.Password)
	if err != nil {
		if errors.Is(err, service.ErrUnauthorized) {
			return echo.NewHTTPError(http.StatusUnauthorized, map[string]string{
				"error": "Password or Name is incorrect",
			})
		}
		return echo.NewHTTPError(http.StatusInternalServerError, map[string]string{
			"error": "Internal Server Error",
		})
	}
	cookie := http.Cookie{
		Name:    "Auth",
		Value:   token,
		Expires: time.Now().Add(1 * time.Hour),
	}
	c.SetCookie(&cookie)
	return c.JSON(http.StatusCreated, map[string]interface{}{
		"token": token,
	})
}
