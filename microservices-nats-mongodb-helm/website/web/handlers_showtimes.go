package web

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"time"

	modelsMovies "cinema-app/website/pkg/movies/models"
	"cinema-app/website/pkg/showtimes/models"

	"github.com/gorilla/mux"
)

type showtimeTemplateData struct {
	ShowTime  models.ShowTime
	ShowTimes []models.ShowTime
	Movies    string
}

func (app *Application) showtimesList(w http.ResponseWriter, r *http.Request) {

	app.InfoLog.Println("Calling showtimes API...")

	std, err := app.getShowtimes()
	if err != nil {
		app.ErrorLog.Println(err.Error())
		http.Error(w, fmt.Sprintf("Internal Server Error-%v", err), 500)
		return
	}
	app.InfoLog.Println("...Showtimes:", std.ShowTimes)

	// Load template files
	files := []string{
		"ui/html/showtimes/list.page.tmpl",
		"ui/html/base.layout.tmpl",
		"ui/html/footer.partial.tmpl",
	}

	ts, err := template.ParseFS(app.Resources, files...)
	if err != nil {
		app.ErrorLog.Println(err.Error())
		http.Error(w, fmt.Sprintf("Internal Server Error-%v", err), 500)

	}

	err = ts.Execute(w, std)
	if err != nil {
		app.ErrorLog.Println(err.Error())
		http.Error(w, fmt.Sprintf("Internal Server Error-%v", err), 500)
	}
}

func (app *Application) showtimesView(w http.ResponseWriter, r *http.Request) {
	// Get id from incoming url
	vars := mux.Vars(r)
	showtimeID := vars["id"]

	// Get showtimes list from API
	app.InfoLog.Println("Calling showtimes API...showtime id", showtimeID)

	std, err := app.getShowtimeById(showtimeID)
	if err != nil {
		app.ErrorLog.Println(err.Error())
		http.Error(w, fmt.Sprintf("Internal Server Error-%v", err), 500)
		return
	}

	app.InfoLog.Println("... showtime:\n\t", std.ShowTime)

	strMovies := strings.Join(std.ShowTime.Movies, ",")

	app.InfoLog.Println("... strMovies:\n\t", strMovies)

	movies, err := app.getMoviesByIDs(strMovies)
	if err != nil {
		app.ErrorLog.Println(err.Error())
		http.Error(w, fmt.Sprintf("Internal Server Error-%v", err), 500)
		return
	}
	var movieTitles []string
	for _, movie := range movies {
		movieTitles = append(movieTitles, movie.Title)
	}

	std.Movies = strings.Join(movieTitles, ", ")

	// Load template files
	files := []string{
		"ui/html/showtimes/view.page.tmpl",
		"ui/html/base.layout.tmpl",
		"ui/html/footer.partial.tmpl",
	}

	ts, err := template.ParseFS(app.Resources, files...)
	if err != nil {
		app.ErrorLog.Println(err.Error())
		http.Error(w, fmt.Sprintf("Internal Server Error-%v", err), 500)
		return
	}

	err = ts.Execute(w, std)
	if err != nil {
		app.ErrorLog.Println(err.Error())
		http.Error(w, fmt.Sprintf("Internal Server Error-%v", err), 500)
	}
}

func (app *Application) getShowtimes() (*showtimeTemplateData, error) {

	var std showtimeTemplateData

	subj := app.Requests.Showtimes + ".list"
	payload := "Request showtimes list"

	msg, err := app.Conn.Request(subj, []byte(payload), 2*time.Second)
	if err != nil {
		if app.Conn.LastError() != nil {
			return &std, fmt.Errorf("%v for request", app.Conn.LastError())
		}
		return &std, fmt.Errorf("%v for request", err)
	}

	var showtimes []models.ShowTime

	err = json.Unmarshal(msg.Data, &showtimes)
	if err != nil {
		return &std, fmt.Errorf("error on unmarshal showtime list, %v", err)
	}

	fmt.Println("...Showtime:\n\t", showtimes)

	std.ShowTimes = showtimes
	return &std, nil
}

func (app *Application) getShowtimeById(id string) (*showtimeTemplateData, error) {

	var std showtimeTemplateData

	subj := app.Requests.Showtimes + ".get"

	payload := id

	msg, err := app.Conn.Request(subj, []byte(payload), 2*time.Second)
	if err != nil {
		if app.Conn.LastError() != nil {
			return &std, fmt.Errorf("%v for request", app.Conn.LastError())
		}
		return &std, fmt.Errorf("%v for request", err)
	}

	var showtime models.ShowTime

	err = json.Unmarshal(msg.Data, &showtime)
	if err != nil {
		return &std, fmt.Errorf("error on unmarshal showtime, %v", err)
	}

	fmt.Println("...getShowtimeById(id) Showtime:\n\t", showtime)

	std.ShowTime = showtime
	return &std, nil
}

func (app *Application) getMoviesByIDs(IDs string) ([]modelsMovies.Movie, error) {

	var movies []modelsMovies.Movie

	subj := app.Requests.Movies + ".moviesByIDs"

	payload := IDs

	msg, err := app.Conn.Request(subj, []byte(payload), 2*time.Second)
	if err != nil {
		if app.Conn.LastError() != nil {
			return movies, fmt.Errorf("%v for request", app.Conn.LastError())
		}
		return movies, fmt.Errorf("%v for request", err)
	}

	err = json.Unmarshal(msg.Data, &movies)
	if err != nil {
		return movies, fmt.Errorf("error on unmarshal showtime, %v", err)
	}

	fmt.Println("... getMoviesByIDs(IDs) movies:\n\t", movies)

	return movies, nil
}
