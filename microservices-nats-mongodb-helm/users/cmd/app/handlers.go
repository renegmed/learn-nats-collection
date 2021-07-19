package main

import (
	"encoding/json"
	"net/http"

	"cinema-app/users/pkg/models"

	"github.com/gorilla/mux"
)

func (app *application) allUsers() ([]byte, error) { //([]*models.User, error) {
	// Get all user stored
	users, err := app.users.All()
	if err != nil {
		return nil, err
	}

	app.infoLog.Println("...all users:\n", users)

	// Convert user list into json encoding
	b, err := json.Marshal(users)
	if err != nil {
		return nil, err
	}
	return b, nil

	//app.infoLog.Println("...Users list\n\t", users)
}
func (app *application) all(w http.ResponseWriter, r *http.Request) {

	bUsers, err := app.allUsers()
	if err != nil {
		app.infoLog.Println(err)
		app.serverError(w, err)
		return
	}
	// Send response back
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(bUsers)
}

func (app *application) findByID(w http.ResponseWriter, r *http.Request) {
	// Get id from incoming url
	vars := mux.Vars(r)
	id := vars["id"]

	app.infoLog.Println("...user to find:\n", id)
	// Find user by id
	m, err := app.users.FindByID(id)
	if err != nil {
		if err.Error() == "ErrNoDocuments" {
			app.infoLog.Println("User not found")
			return
		}
		// Any other error will send an internal server error
		app.serverError(w, err)
	}

	// Convert user to json encoding
	b, err := json.Marshal(m)
	if err != nil {
		app.serverError(w, err)
	}

	app.infoLog.Println("...Have been found a user\n", m)

	// Send response back
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func (app *application) findUserByID(id string) (*models.User, error) {

	app.infoLog.Println("...user to find:\n", id)
	// Find user by id
	user, err := app.users.FindByID(id)
	if err != nil {
		if err.Error() == "ErrNoDocuments" {
			app.infoLog.Println("User not found")
		}
		return user, err
	}

	return user, nil

}

func (app *application) insert(w http.ResponseWriter, r *http.Request) {
	// Define user model
	var u models.User
	// Get request information
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		app.serverError(w, err)
	}

	// // Insert new user
	// insertResult, err := app.users.Insert(u)
	// if err != nil {
	// 	app.serverError(w, err)
	// }

	// app.infoLog.Printf("New user have been created, id=%s", insertResult.InsertedID)

	err = app.insertUser(u)
	if err != nil {
		app.serverError(w, err)
	}

}

func (app *application) insertUser(user models.User) error {

	app.infoLog.Println("...user to insert:\n", user)

	// Insert new user
	insertResult, err := app.users.Insert(user)
	if err != nil {
		app.infoLog.Println("...Error on insert user,", err)
		return err
	}

	app.infoLog.Printf("New user has been created, id=%s", insertResult.InsertedID)

	return nil
}

func (app *application) delete(w http.ResponseWriter, r *http.Request) {
	// Get id from incoming url
	vars := mux.Vars(r)
	id := vars["id"]

	app.infoLog.Println("...user ID to delete:\n", id)

	err := app.deleteUser(id)
	if err != nil {
		app.serverError(w, err)
	}

}

func (app *application) deleteUser(userId string) error {

	app.infoLog.Println("...user to delete:\n", userId)

	deleteResult, err := app.users.Delete(userId)
	if err != nil {
		app.infoLog.Println("...Error on delete user,", err)
		return err
	}

	app.infoLog.Printf("Number of users deleted =%d", deleteResult.DeletedCount)

	return nil
}
