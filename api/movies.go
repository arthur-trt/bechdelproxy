package api

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/arthur-trt/bechdelproxy/database"
	log "github.com/arthur-trt/bechdelproxy/log"
)

func getRatingByIMDBID (c echo.Context) error {
	log.Info("Querying movie by IMDB ID")
	imdbid := c.Param("imdbid")
	var result database.Movie

	log.Debug("IMDB ID: " + imdbid)

	if err := database.Conn.Where("imdb_id = ?", imdbid).First(&result).Error; err != nil {
		log.Warn("IMDB ID: " + imdbid + " not found in database")
		return echo.NewHTTPError(http.StatusNotFound, "Movie not in database")
	}

	return c.JSON(http.StatusOK, result)
}
