package web

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"time"

	modelsBooking "cinema-app/website/pkg/bookings/models"
	modelsShowTime "cinema-app/website/pkg/showtimes/models"

	modelsUser "cinema-app/website/pkg/users/models"

	"github.com/gorilla/mux"
)

type bookingTemplateData struct {
	Booking      modelsBooking.Booking
	Bookings     []modelsBooking.Booking
	BookingData  bookingData
	BookingsData []bookingData
}

type bookingData struct {
	ID           string
	UserFullName string
	ShowTimeDate string
}

func (app *Application) loadBookingData(btd *bookingTemplateData, isList bool) {
	// Clean booking data
	btd.BookingsData = []bookingData{}
	btd.BookingData = bookingData{}

	// Load booking data
	if isList {
		for _, b := range btd.Bookings {

			// Load user data
			user, err := app.loadUserData(b.UserID)
			if err != nil {
				app.ErrorLog.Println(err)
				continue
			}

			// Load showtime data
			showtime, err := app.loadShowtimeData(b.ShowtimeID)
			if err != nil {
				app.ErrorLog.Println(err)
				continue
			}

			bookingData := bookingData{
				ID:           b.ID.Hex(),
				UserFullName: fmt.Sprintf("%s %s", user.Name, user.LastName),
				ShowTimeDate: showtime.Date,
			}
			btd.BookingsData = append(btd.BookingsData, bookingData)
			app.InfoLog.Println(b.UserID)
		}
	} else {
		b := btd.Booking

		user, err := app.loadUserData(b.UserID)
		if err != nil {
			app.ErrorLog.Println(err)
			return
		}

		showtime, err := app.loadShowtimeData(b.ShowtimeID)
		if err != nil {
			app.ErrorLog.Println(err)
			return
		}
		btd.BookingData = bookingData{
			ID:           b.ID.Hex(),
			UserFullName: fmt.Sprintf("%s %s", user.Name, user.LastName),
			ShowTimeDate: showtime.Date,
		}
	}
}

func (app *Application) loadUserData(id string) (modelsUser.User, error) {
	var user modelsUser.User

	subj := app.Requests.Users + ".get"
	payload := id
	msg, err := app.Conn.Request(subj, []byte(payload), 2*time.Second)
	if err != nil {
		app.ErrorLog.Println(err)
		return user, err
	}

	app.InfoLog.Println("...loadUserData user\n\t", string(msg.Data))

	err = json.Unmarshal(msg.Data, &user)
	if err != nil {
		app.ErrorLog.Println(err)
		return user, err
	}

	app.InfoLog.Println("...loadUserData user\n\t", user)

	return user, nil
}

func (app *Application) loadShowtimeData(id string) (modelsShowTime.ShowTime, error) {
	var showtime modelsShowTime.ShowTime

	subj := app.Requests.Showtimes + ".get"
	payload := id
	msg, err := app.Conn.Request(subj, []byte(payload), 2*time.Second)
	if err != nil {
		app.ErrorLog.Println(err)
		return showtime, err
	}
	app.InfoLog.Println("...loadShowtimeData showtime data\n\t", string(msg.Data))

	err = json.Unmarshal(msg.Data, &showtime)
	if err != nil {
		app.ErrorLog.Println(err)
		return showtime, err
	}

	app.InfoLog.Println("...loadShowtimeData showtime\n\t", showtime)

	return showtime, nil
}

func (app *Application) bookingsList(w http.ResponseWriter, r *http.Request) {

	app.InfoLog.Println("Calling bookings API...")

	btd, err := app.getBookings()
	if err != nil {
		app.ErrorLog.Println(err.Error())
		http.Error(w, fmt.Sprintf("Internal Server Error-%v", err), 500)
		return
	}
	app.InfoLog.Println("....Bookings:", btd.Bookings)

	app.loadBookingData(btd, true)

	// Load template files
	files := []string{
		"ui/html/bookings/list.page.tmpl",
		"ui/html/base.layout.tmpl",
		"ui/html/footer.partial.tmpl",
	}

	ts, err := template.ParseFS(app.Resources, files...)
	if err != nil {
		app.ErrorLog.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}

	err = ts.Execute(w, btd)
	if err != nil {
		app.ErrorLog.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
}

func (app *Application) getBookings() (*bookingTemplateData, error) {

	var btd bookingTemplateData

	subj := app.Requests.Bookings + ".list"
	payload := "Request bookings list"
	msg, err := app.Conn.Request(subj, []byte(payload), 2*time.Second)
	if err != nil {
		if app.Conn.LastError() != nil {
			return &btd, fmt.Errorf("%v for request", app.Conn.LastError())
		}
		return &btd, fmt.Errorf("%v for request", err)
	}

	var bookings []modelsBooking.Booking

	err = json.Unmarshal(msg.Data, &bookings)
	if err != nil {
		return &btd, fmt.Errorf("error on unmarshal user list, %v", err)
	}

	fmt.Println("...Bookings:\n\t", bookings)

	btd.Bookings = bookings
	return &btd, nil
}

func (app *Application) bookingsView(w http.ResponseWriter, r *http.Request) {
	// Get id from incoming url
	vars := mux.Vars(r)
	bookingID := vars["id"]

	booking, err := app.getBookingById(bookingID)
	if err != nil {
		app.ErrorLog.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}
	// 	// Get bookings list from API
	var td bookingTemplateData
	td.Booking = booking

	app.loadBookingData(&td, false)

	// Load template files
	files := []string{
		"ui/html/bookings/view.page.tmpl",
		"ui/html/base.layout.tmpl",
		"ui/html/footer.partial.tmpl",
	}

	ts, err := template.ParseFS(app.Resources, files...)
	if err != nil {
		app.ErrorLog.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}

	err = ts.Execute(w, td)
	if err != nil {
		app.ErrorLog.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
}

func (app *Application) getBookingById(id string) (modelsBooking.Booking, error) {
	var booking modelsBooking.Booking

	subj := app.Requests.Bookings + ".get"
	payload := id
	msg, err := app.Conn.Request(subj, []byte(payload), 2*time.Second)
	if err != nil {
		app.ErrorLog.Println(err)
		return booking, err
	}
	app.InfoLog.Println("...getBookingById booking msg data\n\t", string(msg.Data))

	err = json.Unmarshal(msg.Data, &booking)
	if err != nil {
		app.ErrorLog.Println(err)
		return booking, err
	}

	app.InfoLog.Println("...getBookingById booking\n\t", booking)

	return booking, nil
}
