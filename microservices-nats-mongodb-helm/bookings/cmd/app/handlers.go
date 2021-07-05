package main

import (
	"encoding/json"
	"net/http"

	"cinema-app/bookings/pkg/models"

	"github.com/gorilla/mux"
)

func (app *application) all(w http.ResponseWriter, r *http.Request) {

	b, err := app.allBookings()
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

func (app *application) allBookings() ([]byte, error) {
	// Get all bookings stored
	bookings, err := app.bookings.All()
	if err != nil {
		return nil, err
	}

	app.infoLog.Println("...bookings:\n", bookings)

	// Convert showtime list into json encoding
	b, err := json.Marshal(bookings)
	if err != nil {
		return nil, err
	}

	app.infoLog.Println("...bookings list:", string(b))

	return b, nil
}

func (app *application) findByID(w http.ResponseWriter, r *http.Request) {
	// Get id from incoming url
	vars := mux.Vars(r)
	id := vars["id"]

	app.infoLog.Println("...booking of find:\n", id)

	booking, err := app.findBookingByID(id)
	if err != nil {
		app.infoLog.Println(err)
		app.serverError(w, err)
		return
	}
	// Convert booking to json encoding
	bBooking, err := json.Marshal(booking)
	if err != nil {
		app.infoLog.Println(err)
		app.serverError(w, err)
		return
	}

	// Send response back
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(bBooking)
}
func (app *application) findBookingByID(id string) (*models.Booking, error) {

	booking, err := app.bookings.FindByID(id)
	if err != nil {
		if err.Error() == "ErrNoDocuments" {
			app.infoLog.Println("Booking not found")

		}
		return booking, err
	}
	app.infoLog.Println("...findBookingByID booking", booking)

	return booking, nil
}

func (app *application) insert(w http.ResponseWriter, r *http.Request) {

	var m models.Booking

	err := json.NewDecoder(r.Body).Decode(&m)
	if err != nil {
		app.serverError(w, err)
	}

	app.infoLog.Println("...booking to insert:\n", m)

	// Insert new booking
	insertResult, err := app.bookings.Insert(m)
	if err != nil {
		app.serverError(w, err)
	}

	app.infoLog.Printf("...New booking have been created, id=%s", insertResult.InsertedID)
}

func (app *application) delete(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]

	app.infoLog.Println("...booking to delete:\n", id)

	deleteResult, err := app.bookings.Delete(id)
	if err != nil {
		app.serverError(w, err)
	}

	app.infoLog.Printf("...Have been deleted %d booking(s)", deleteResult.DeletedCount)
}
