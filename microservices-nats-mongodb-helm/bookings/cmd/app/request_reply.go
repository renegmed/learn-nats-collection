package main

import "encoding/json"

func (app *application) reply_allBookings() ([]byte, error) {
	// Get all user stored
	bBookings, err := app.allBookings()
	if err != nil {
		return nil, err
	}

	app.infoLog.Printf("...all Bookings:\n\t%v", string(bBookings))
	return bBookings, nil
}

func (app *application) reply_getBookingById(id string) ([]byte, error) {
	booking, err := app.findBookingByID(id)
	if err != nil {
		return nil, err
	}

	bBooking, err := json.Marshal(booking)
	if err != nil {
		return bBooking, nil
	}
	app.infoLog.Printf("...get user:\n\t%v", string(bBooking))
	return bBooking, nil
}
