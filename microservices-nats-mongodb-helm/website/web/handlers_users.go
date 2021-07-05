package web

import (
	"encoding/json"
	"fmt"

	"net/http"
	"text/template"
	"time"

	"cinema-app/website/pkg/users/models"

	"github.com/gorilla/mux"
)

type userTemplateData struct {
	User  models.User
	Users []models.User
}

func (app *Application) usersList(w http.ResponseWriter, r *http.Request) {

	utd, err := app.getUsers()
	if err != nil {
		app.ErrorLog.Println(err.Error())
		http.Error(w, fmt.Sprintf("Internal Server Error-%v", err), 500)
		return
	}
	app.InfoLog.Println("....Users:", utd.Users)

	// Load template files
	files := []string{
		"ui/html/users/list.page.tmpl",
		"ui/html/base.layout.tmpl",
		"ui/html/footer.partial.tmpl",
	}

	ts, err := template.ParseFS(app.Resources, files...)
	if err != nil {
		app.ErrorLog.Println(err.Error())
		http.Error(w, fmt.Sprintf("Internal Server Error-%v", err), 500)
		return
	}

	err = ts.Execute(w, utd)
	if err != nil {
		app.ErrorLog.Println(err.Error())
		http.Error(w, fmt.Sprintf("Internal Server Error-%v", err), 500)
	}
}

func (app *Application) getUsers() (*userTemplateData, error) {

	var utd userTemplateData

	subj := app.Requests.Users + ".list"
	payload := "Request users list"
	msg, err := app.Conn.Request(subj, []byte(payload), 2*time.Second)
	if err != nil {
		if app.Conn.LastError() != nil {
			return &utd, fmt.Errorf("%v for request", app.Conn.LastError())
		}
		return &utd, fmt.Errorf("%v for request", err)
	}

	var users []models.User

	err = json.Unmarshal(msg.Data, &users)
	if err != nil {
		return &utd, fmt.Errorf("Error on unmarshal user list, %v", err)
	}

	fmt.Println("...Users:\n\t", users)

	utd.Users = users
	return &utd, nil
}

func (app *Application) usersView(w http.ResponseWriter, r *http.Request) {
	// Get id from incoming url
	vars := mux.Vars(r)
	userID := vars["id"]

	// Get users list from API
	app.InfoLog.Println("Calling usersView API...")

	user, err := app.getUserByID(userID)
	if err != nil {
		app.ErrorLog.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}

	app.InfoLog.Println("...UsersView user:\n\t", user)

	// Load template files
	files := []string{
		"ui/html/users/view.page.tmpl",
		"ui/html/base.layout.tmpl",
		"ui/html/footer.partial.tmpl",
	}

	ts, err := template.ParseFS(app.Resources, files...)
	if err != nil {
		app.ErrorLog.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}

	err = ts.Execute(w, *user)
	if err != nil {
		app.ErrorLog.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
}

func (app *Application) getUserByID(id string) (*models.User, error) {
	var user models.User

	subj := app.Requests.Users + ".get"
	payload := id
	msg, err := app.Conn.Request(subj, []byte(payload), 2*time.Second)
	if err != nil {
		return &user, err
	}

	app.InfoLog.Println("...handlers_users user\n\t", string(msg.Data))

	err = json.Unmarshal(msg.Data, &user)
	if err != nil {
		return &user, err
	}

	return &user, nil
}
