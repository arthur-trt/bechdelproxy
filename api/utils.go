package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func getPing (c echo.Context) error {
	return c.String(http.StatusOK, "Pong")
}
