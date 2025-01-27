package movies

import (
	"encoding/json"
	"html"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	db "github.com/arthur-trt/bechdelproxy/database"
	log "github.com/arthur-trt/bechdelproxy/log"
)

func preload() (map[string]db.Movie, error) {
	var movies []db.Movie

	log.Debug("Preload all movie from local DB")
	if err := db.Conn.Find(&movies).Error; err != nil {
		return nil, err
	}

	movieMap := make(map[string]db.Movie, len(movies))
	for _, movie := range movies {
		movieMap[movie.IMDBID] = movie
	}
	return movieMap, nil
}

func fetchJSON() ([]db.Movie, error) {
	var movies []db.Movie

	if log_level := os.Getenv("LOG_LEVEL"); log_level == "DEBUG" {
		file, _ := os.ReadFile("raw_data.json")
		if err := json.Unmarshal(file, &movies); err != nil {
			return nil, err
		}
	} else {
		log.Info("Querying bechdeltest.com for movies")
		client := &http.Client{Timeout: 30 * time.Second}

		response, err := client.Get("https://bechdeltest.com/api/v1/getAllMovies")
		if err != nil {
			return nil, err
		}
		defer response.Body.Close()

		body, err := io.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal(body, &movies); err != nil {
			return nil, err
		}
	}
	return movies, nil
}

func filterJSON(raw_movies []db.Movie) []db.Movie {
	movieMap := make(map[string]db.Movie)

	for _, movie := range raw_movies {
		if len(movie.IMDBID) > 0 && !strings.HasPrefix(movie.IMDBID, "tt") {
			movie.IMDBID = "tt" + movie.IMDBID
		}
		movie.Title = html.UnescapeString(movie.Title)
		if existing, exists := movieMap[movie.IMDBID]; exists {
			if movie.BechdelID > existing.BechdelID {
				movieMap[movie.IMDBID] = movie
			}
		} else {
			movieMap[movie.IMDBID] = movie
		}
	}

	result := make([]db.Movie, 0, len(movieMap))
	for _, movie := range movieMap {
		result = append(result, movie)
	}

	return result
}

func process(existingMovies map[string]db.Movie, newMovies []db.Movie) (newMovie []db.Movie) {
	for _, movie := range newMovies {
		if existing, exists := existingMovies[movie.IMDBID]; exists {
			if movie.BechdelID > existing.BechdelID {
				existing.Title = movie.Title
				existing.BechdelID = movie.BechdelID
				existing.Rating = movie.Rating
				newMovie = append(newMovie, existing)
			}
		} else {
			newMovie = append(newMovie, movie)
		}
	}

	return newMovie
}

func Update() error {
	rawJSON, err := fetchJSON()
	log.Debug("Fetch JSON raw data")
	if err != nil {
		log.Error("Couldn't retrieve JSON: ", err)
		return err
	}

	filteredMovie := filterJSON(rawJSON)
	log.Debug("Filter movies")

	existingMovies, err := preload()
	if err != nil {
		log.Error("Couldn't retrive existing movie: ", err)
		return err
	}

	newMovie := process(existingMovies, filteredMovie)

	if err := db.InsertOrUpdateMovies(newMovie); err != nil {
		log.Error("Error while instering movies: ", err)
		return err
	}
	return nil
}
