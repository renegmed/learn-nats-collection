package main

import (
	"encoding/json"
)

func (app *application) reply_allUsers() ([]byte, error) {
	// Get all user stored
	bUsers, err := app.allUsers()
	if err != nil {
		return nil, err
	}

	app.infoLog.Printf("...all users:\n\t%v", string(bUsers))
	return bUsers, nil
}

func (app *application) reply_getUser(userId string) ([]byte, error) {

	user, err := app.findUserByID(userId)
	if err != nil {
		return nil, err
	}
	//var bUser models.User
	bUser, err := json.Marshal(user)
	if err != nil {
		return bUser, nil
	}
	app.infoLog.Printf("...get user:\n\t%v", string(bUser))
	return bUser, nil
}
