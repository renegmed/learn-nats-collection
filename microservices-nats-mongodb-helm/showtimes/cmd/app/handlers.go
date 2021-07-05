package main

import (
	"encoding/json"

	"net/http"
	"time"

	"cinema-app/showtimes/pkg/models"

	"github.com/gorilla/mux"
)

func (app *application) allShowtimes() ([]byte, error) {
	// Get all showtimes stored
	showtimes, err := app.showtimes.All()
	if err != nil {
		return nil, err
	}

	app.infoLog.Println("...showtimes:\n", showtimes)

	// Convert showtime list into json encoding
	b, err := json.Marshal(showtimes)
	if err != nil {
		return nil, err
	}

	app.infoLog.Println("...Showtimes list:", string(b))

	return b, nil

}

func (app *application) all(w http.ResponseWriter, r *http.Request) {

	b, err := app.allShowtimes()
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

func (app *application) findByID(w http.ResponseWriter, r *http.Request) {
	// Get id from incoming url
	vars := mux.Vars(r)
	id := vars["id"]

	showtime, err := app.findShowtimeByID(id)
	if err != nil {
		app.serverError(w, err)
	}

	// Convert showtime to json encoding
	bShowtime, err := json.Marshal(showtime)
	if err != nil {
		app.serverError(w, err)
	}
	// Send response back
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(bShowtime)
}

func (app *application) findShowtimeByID(id string) (*models.ShowTime, error) {

	app.infoLog.Println("...showtime id:\n", id)

	// Find showtime by id
	showtime, err := app.showtimes.FindByID(id)
	if err != nil {
		if err.Error() == "ErrNoDocuments" {
			app.infoLog.Println("Showtime not found")
		}
		return nil, err
	}

	app.infoLog.Println("...Found showtime\n\t", showtime)

	return showtime, nil
}

func (app *application) findByDate(w http.ResponseWriter, r *http.Request) {
	// Get id from incoming url
	vars := mux.Vars(r)
	date := vars["date"]

	app.infoLog.Println("...showtimes to find by date:\n", date)

	// Find showtime by date
	m, err := app.showtimes.FindByDate(date)
	if err != nil {
		if err.Error() == "ErrNoDocuments" {
			app.infoLog.Println("Showtime not found")
			return
		}
		// Any other error will send an internal server error
		app.serverError(w, err)
	}

	// Convert showtime to json encoding
	b, err := json.Marshal(m)
	if err != nil {
		app.serverError(w, err)
	}

	app.infoLog.Println("...Have been found a showtime by date:\n", m)

	// Send response back
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func (app *application) insert(w http.ResponseWriter, r *http.Request) {
	// Define showtime model
	var m models.ShowTime
	// Get request information
	err := json.NewDecoder(r.Body).Decode(&m)
	if err != nil {
		app.serverError(w, err)
	}

	app.infoLog.Println("...showtime to insert:\n", m)

	// Insert new showtime
	m.CreatedAt = time.Now()
	insertResult, err := app.showtimes.Insert(m)
	if err != nil {
		app.serverError(w, err)
	}

	app.infoLog.Printf("...New showtime have been created, id=%s", insertResult.InsertedID)
}

func (app *application) delete(w http.ResponseWriter, r *http.Request) {
	// Get id from incoming url
	vars := mux.Vars(r)
	id := vars["id"]

	app.infoLog.Println("...showtime to delete:\n", id)

	// Delete showtime by id
	deleteResult, err := app.showtimes.Delete(id)
	if err != nil {
		app.serverError(w, err)
	}

	app.infoLog.Printf("...Have been eliminated %d showtime(s)", deleteResult.DeletedCount)
}
