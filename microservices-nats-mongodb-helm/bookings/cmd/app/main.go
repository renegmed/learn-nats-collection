package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"cinema-app/bookings/pkg/models/mongodb"

	nats "github.com/nats-io/nats.go"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	bookings *mongodb.BookingModel
	conn     *nats.Conn
}

func main() {

	// Define command-line flags
	serverAddr := flag.String("serveraddr", "", "HTTP server network address")
	serversUrl := flag.String("s", nats.DefaultURL, "The nats server URLs (separated by comma)")
	topic := flag.String("t", "", "NATS topic")
	queueName := flag.String("q", "", "NATS queue name")
	mongoURI := flag.String("mongouri", "mongodb://localhost:27017", "Database hostname url")
	mongoDatabase := flag.String("databasename", "bookings", "Database name")
	enableCredentials := flag.Bool("enablecredentials", true, "Enable the use of credentials for mongo connection")
	flag.Parse()

	// Create logger for writing information and error messages.
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	fmt.Println(".... enable credentials?", *enableCredentials)

	opts := []nats.Option{nats.Name("NATS Bookings")}

	// Connect to NATS
	nconn, err := nats.Connect(*serversUrl, opts...)
	if err != nil {

		log.Fatal(err)
	}
	defer nconn.Close()

	log.Println("...Got NATS connection")

	// Create mongo client configuration
	co := options.Client().ApplyURI(*mongoURI)

	if *enableCredentials {
		co.Auth = &options.Credential{
			Username: os.Getenv("MONGODB_USERNAME"),
			Password: os.Getenv("MONGODB_PASSWORD"),
		}
	}

	//infoLog.Printf("... mongo client %v", co)

	// Establish database connection
	client, err := mongo.NewClient(co)
	if err != nil {
		errLog.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		errLog.Fatal(err)
	}

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	infoLog.Printf("Database connection established")

	// Initialize a new instance of application containing the dependencies.
	app := &application{
		infoLog:  infoLog,
		errorLog: errLog,
		bookings: &mongodb.BookingModel{
			C: client.Database(*mongoDatabase).Collection("bookings"),
		},
		conn: nconn,
	}

	nconn.QueueSubscribe(*topic+".list", *queueName+"_list", func(msg *nats.Msg) {

		infoLog.Println("...QueueSubscribe called - topic", *topic)

		bBookings, err := app.reply_allBookings()
		if err == nil {
			msg.Respond(bBookings)
			return
		}

		app.errorLog.Println("...Error from reply all Bookings,", err)
	})

	nconn.QueueSubscribe(*topic+".get", *queueName+"_get", func(msg *nats.Msg) {

		infoLog.Println("...QueueSubscribe called - topic", *topic)
		app.infoLog.Printf("...Subject: %s  Data: %s", msg.Subject, string(msg.Data))
		bBooking, err := app.reply_getBookingById(string(msg.Data))
		if err == nil {
			msg.Respond(bBooking)
			return
		}

		app.errorLog.Println("...Error from reply get booking by ID,", err)
	})

	nconn.Flush()

	if err := nconn.LastError(); err != nil {
		log.Fatal(err)
	}

	log.Printf("Listening on topic: %s", *topic)

	srv := &http.Server{
		Addr:         *serverAddr,
		ErrorLog:     errLog,
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	infoLog.Printf("Starting server on %s", *serverAddr)
	err = srv.ListenAndServe()
	errLog.Fatal(err)

	// Setup the interrupt handler to drain so we don't miss
	// requests when scaling down.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	log.Println()
	log.Printf("Draining...")
	nconn.Drain()
	log.Fatalf("Exiting")
}
