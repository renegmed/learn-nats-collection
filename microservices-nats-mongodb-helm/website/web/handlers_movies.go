package web

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"time"

	"cinema-app/website/pkg/movies/models"

	"github.com/gorilla/mux"
)

type movieTemplateData struct {
	Movie  models.Movie
	Movies []models.Movie
}

func (app *Application) moviesList(w http.ResponseWriter, r *http.Request) {

	app.InfoLog.Println("Calling movies list API....")

	mtd, err := app.getMovies()
	if err != nil {
		app.ErrorLog.Println(err.Error())
		http.Error(w, fmt.Sprintf("Internal Server Error-%v", err), 500)
		return
	}
	app.InfoLog.Println("....Movies:", mtd.Movies)

	files := []string{
		"ui/html/movies/list.page.tmpl",
		"ui/html/base.layout.tmpl",
		"ui/html/footer.partial.tmpl",
	}

	ts, err := template.ParseFS(app.Resources, files...)
	if err != nil {
		app.ErrorLog.Println(err.Error())
		http.Error(w, fmt.Sprintf("Internal Server Error-%v", err), 500)
		return
	}

	err = ts.Execute(w, mtd)
	if err != nil {
		app.ErrorLog.Println(err.Error())
		http.Error(w, fmt.Sprintf("Internal Server Error-%v", err), 500)
	}
}

func (app *Application) getMovies() (*movieTemplateData, error) {

	var mtd movieTemplateData

	subj := app.Requests.Movies + ".list"
	payload := "Request movies list"
	msg, err := app.Conn.Request(subj, []byte(payload), 2*time.Second)
	if err != nil {
		if app.Conn.LastError() != nil {
			return &mtd, fmt.Errorf("%v for request", app.Conn.LastError())
		}
		return &mtd, fmt.Errorf("%v for request", err)
	}

	var movies []models.Movie

	err = json.Unmarshal(msg.Data, &movies)
	if err != nil {
		return &mtd, fmt.Errorf("error on unmarshal movie list, %v", err)
	}

	fmt.Println("...Movies:\n\t", movies)

	mtd.Movies = movies
	return &mtd, nil
}

func (app *Application) moviesView(w http.ResponseWriter, r *http.Request) {
	// 	// Get id from incoming url
	vars := mux.Vars(r)
	movieID := vars["id"]

	app.InfoLog.Println("....Movie ID to find:", movieID)

	movie, err := app.getMovieByID(movieID)
	if err != nil {
		app.ErrorLog.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
	// Load template files
	files := []string{
		"ui/html/movies/view.page.tmpl",
		"ui/html/base.layout.tmpl",
		"ui/html/footer.partial.tmpl",
	}

	ts, err := template.ParseFS(app.Resources, files...)
	if err != nil {
		app.ErrorLog.Println(err.Error())
		http.Error(w, fmt.Sprintf("Internal Server Error-%v", err), 500)
		return
	}

	err = ts.Execute(w, movie)
	if err != nil {
		app.ErrorLog.Println(err.Error())
		http.Error(w, fmt.Sprintf("Internal Server Error-%v", err), 500)
	}
}

func (app *Application) getMovieByID(id string) (*models.Movie, error) {
	var movie models.Movie

	subj := app.Requests.Movies + ".get"
	payload := id
	msg, err := app.Conn.Request(subj, []byte(payload), 2*time.Second)
	if err != nil {
		return &movie, err
	}

	app.InfoLog.Println("...handlers_movies movie\n\t", string(msg.Data))

	err = json.Unmarshal(msg.Data, &movie)
	if err != nil {
		return &movie, err
	}

	return &movie, nil
}

func (app *Application) moviesAdd(w http.ResponseWriter, r *http.Request) {
	var movie models.Movie

	err := json.NewDecoder(r.Body).Decode(&movie)
	if err != nil {
		http.Error(w, "Internal Server Error - decoding movie", 500)
	}

	app.InfoLog.Println("...movie to add:\n", movie)

	mve, err := json.Marshal(movie)
	if err != nil {
		http.Error(w, "Internal Server Error - marshal movie", 500)
	}

	subj := app.Requests.Movies + ".add"

	_, err = app.Conn.Request(subj, mve, 2*time.Second)
	if err != nil {
		http.Error(w, "Internal Server Error - request NATS connection", 500)
	}

}

func (app *Application) moviesDelete(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	movieID := vars["id"]

	app.InfoLog.Println("...movie to delete:", movieID)

	subj := app.Requests.Movies + ".delete"
	_, err := app.Conn.Request(subj, []byte(movieID), 2*time.Second)
	if err != nil {
		http.Error(w, "Internal Server Error - request NATS connection", 500)
	}
}
