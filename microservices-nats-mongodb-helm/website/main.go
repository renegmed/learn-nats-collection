package main

import (
	"embed"
	"flag"

	"log"
	"net/http"
	"os"
	"time"

	"cinema-app/website/web"

	nats "github.com/nats-io/nats.go"
)

var (
	//go:embed ui/*
	res embed.FS
)

func main() {

	// Define command-line flags
	serverAddr := flag.String("serveraddr", "", "HTTP server network address")
	serversUrl := flag.String("s", nats.DefaultURL, "The nats server URLs (separated by comma)")
	usersTopic := flag.String("users", "", "User request subject")
	moviesTopic := flag.String("movies", "", "Movies request subject")
	showtimesTopic := flag.String("showtimes", "", "Showtimes request subject")
	bookingsTopic := flag.String("bookings", "", "Bookings request subject")
	flag.Parse()

	// Create logger for writing information and error messages.
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	opts := []nats.Option{nats.Name("NATS Sample Requestor")}
	// Connect to NATS
	nc, err := nats.Connect(*serversUrl, opts...)
	if err != nil {

		log.Fatal(err)
	}
	defer nc.Close()

	log.Println("...Go NATS connection")

	// Initialize a new instance of application containing the dependencies.
	app := &web.Application{
		InfoLog:  infoLog,
		ErrorLog: errLog,
		Requests: web.Requests{
			Users:     *usersTopic,
			Movies:    *moviesTopic,
			Showtimes: *showtimesTopic,
			Bookings:  *bookingsTopic,
		},
		Resources: &res,
		Conn:      nc,
	}

	// Initialize a new http.Server struct.
	// serverURI := fmt.Sprintf("%s:%d", *serverAddr, *serverPort)
	srv := &http.Server{
		Addr:         *serverAddr,
		ErrorLog:     errLog,
		Handler:      app.Routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	infoLog.Printf("Starting server on %s", *serverAddr)
	err = srv.ListenAndServe()
	errLog.Fatal(err)
}
