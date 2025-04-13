package handler

import "github.com/labstack/echo/v4"

func Access(c echo.Context) error {
	return c.JSON(200, map[string]string{
		"status": c.Get("name").(string),
	})
}
