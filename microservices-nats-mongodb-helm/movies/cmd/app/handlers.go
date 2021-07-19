package main

import (
	"encoding/json"
	"net/http"

	"cinema-app/movies/pkg/models"

	"github.com/gorilla/mux"
)

func (app *application) all(w http.ResponseWriter, r *http.Request) {

	b, err := app.allMovies()
	if err != nil {
		app.infoLog.Println(err)
		app.serverError(w, err)
		return
	}

	// Send response back
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func (app *application) allMovies() ([]byte, error) {

	movies, err := app.movies.All()
	if err != nil {
		return nil, err
	}

	app.infoLog.Println("...movies:\n", movies)

	// Convert showtime list into json encoding
	b, err := json.Marshal(movies)
	if err != nil {
		return nil, err
	}

	app.infoLog.Println("...movies list:", string(b))

	return b, nil

}

func (app *application) findByID(w http.ResponseWriter, r *http.Request) {
	// Get id from incoming url
	vars := mux.Vars(r)
	id := vars["id"]

	m, err := app.getMovieByID(id)
	if err != nil {
		app.serverError(w, err)
	}
	// Convert movie to json encoding
	b, err := json.Marshal(m)
	if err != nil {
		app.serverError(w, err)
	}

	app.infoLog.Println("...Found a movie", *m)

	// Send response back
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func (app *application) getMovieByID(id string) (*models.Movie, error) {

	// Find movie by id
	m, err := app.movies.FindByID(id)
	if err != nil {
		if err.Error() == "ErrNoDocuments" {
			app.infoLog.Println("...Movie not found with ID", id)
			return m, err
		}
		// Any other error will send an internal server error
		app.infoLog.Println("...error on findByID() data access,", err)
		return m, err
	}

	app.infoLog.Println("...movie:\n", m)

	return m, nil
}

func (app *application) insert(w http.ResponseWriter, r *http.Request) {
	// Define movie model
	var m models.Movie
	// Get request information
	err := json.NewDecoder(r.Body).Decode(&m)
	if err != nil {
		app.infoLog.Println("...error on insert() decoding,", err)
		app.serverError(w, err)
	}

	app.infoLog.Println("...movie to insert:\n", m)

	app.insertMovie(m)

}

func (app *application) insertMovie(movie models.Movie) error {

	app.infoLog.Println("...movie to insert:\n", movie)

	insertResult, err := app.movies.Insert(movie)
	if err != nil {
		app.infoLog.Println("...error on insert() data access,", err)
		return err
	}

	app.infoLog.Printf("...New movie have been created, id=%s", insertResult.InsertedID)

	return nil
}

func (app *application) delete(w http.ResponseWriter, r *http.Request) {
	// Get id from incoming url
	vars := mux.Vars(r)
	id := vars["id"]

	app.infoLog.Println("...movie to delete:", id)

	app.deleteMovie(id)

}

func (app *application) deleteMovie(movieId string) error {

	app.infoLog.Println("...movie to delete:\n", movieId)

	// Delete movie by id
	deleteResult, err := app.movies.Delete(movieId)
	if err != nil {
		app.infoLog.Println("...error on delete() data access,", err)
		return err
	}

	app.infoLog.Printf("...Have been eliminated %d movie(s)", deleteResult.DeletedCount)

	return nil
}
