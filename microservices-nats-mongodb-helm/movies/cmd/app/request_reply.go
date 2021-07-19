package main

import (
	"cinema-app/movies/pkg/models"
	"encoding/json"
	"strings"
)

func (app *application) reply_allMovies() ([]byte, error) {

	bMovies, err := app.allMovies()
	if err != nil {
		return nil, err
	}

	app.infoLog.Printf("...reply_allMovies:\n\t%v", string(bMovies))
	return bMovies, nil
}
func (app *application) reply_getMovieByID(id string) ([]byte, error) {
	movie, err := app.getMovieByID(id)
	if err != nil {
		return nil, err
	}

	bMovie, err := json.Marshal(movie)
	if err != nil {
		return nil, err
	}

	app.infoLog.Printf("...reply_getMovieByID movie:\n\t%v", string(bMovie))
	return bMovie, nil
}

func (app *application) reply_moviesByIDs(IDs string) ([]byte, error) { // a string of movie IDs separated by comma

	movieIds := strings.Split(IDs, ",")

	var movies []models.Movie
	for _, id := range movieIds {
		movie, err := app.getMovieByID(id)
		if err == nil {
			movies = append(movies, *movie)
		}
		app.errorLog.Println(err)
	}

	bMovies, err := json.Marshal(movies)
	if err != nil {
		return nil, err
	}

	app.infoLog.Printf("...reply_moviesByIDs movies:\n\t%v", string(bMovies))
	return bMovies, nil
}

func (app *application) reply_addMovie(movie string) error {
	var m models.Movie
	err := json.Unmarshal([]byte(movie), &m)
	if err != nil {
		return err
	}
	return app.insertMovie(m)
}

func (app *application) reply_deleteMovie(movieId string) error {

	return app.deleteMovie(movieId)
}
