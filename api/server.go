package api

import (
	"os"

	echoLog "github.com/labstack/gommon/log"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/echo/v4"
)

func New() *echo.Echo {
	e := echo.New()

	if env := os.Getenv("LOG_LEVEL"); env == "DEBUG" {
		e.Logger.SetLevel(echoLog.DEBUG)
	}

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowMethods: []string{echo.GET},
	}))

	return e
}
