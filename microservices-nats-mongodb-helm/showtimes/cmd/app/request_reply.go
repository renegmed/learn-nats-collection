package main

import "encoding/json"

func (app *application) reply_allShowtimes() ([]byte, error) {
	// Get all user stored
	bShowtimes, err := app.allShowtimes()
	if err != nil {
		return nil, err
	}

	app.infoLog.Printf("...reply_allShowtimes:\n\t%v", string(bShowtimes))
	return bShowtimes, nil
}

func (app *application) reply_getShowtime(id string) ([]byte, error) {

	showtime, err := app.findShowtimeByID(id)
	if err != nil {
		return nil, err
	}
	app.infoLog.Printf("...reply_getShowtime showtime:\n\t%v", showtime)

	bShowtime, err := json.Marshal(showtime)
	if err != nil {
		return bShowtime, nil
	}
	app.infoLog.Printf("...reply_getShowtime:\n\t%v", string(bShowtime))
	return bShowtime, nil
}
