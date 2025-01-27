package api

import "github.com/labstack/echo/v4"

func Register(e *echo.Echo) {
	e.GET("/ping", getPing)

	e.GET("/imdb/:imdbid", getRatingByIMDBID)
}
