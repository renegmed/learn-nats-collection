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

	"cinema-app/showtimes/pkg/models/mongodb"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	nats "github.com/nats-io/nats.go"
)

type application struct {
	errorLog  *log.Logger
	infoLog   *log.Logger
	showtimes *mongodb.ShowTimeModel
	conn      *nats.Conn
}

func main() {

	serverAddr := flag.String("serveraddr", "", "HTTP server network address")
	serversUrl := flag.String("s", nats.DefaultURL, "The nats server URLs (separated by comma)")
	topic := flag.String("t", "", "NATS topic")
	queueName := flag.String("q", "", "NATS queue name")
	mongoURI := flag.String("mongouri", "mongodb://localhost:27019", "Database hostname url")
	mongoDatabse := flag.String("databasename", "showtimes", "Database name")
	enableCredentials := flag.Bool("enablecredentials", true, "Enable the use of credentials for mongo connection")
	flag.Parse()

	// Create logger for writing information and error messages.
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	opts := []nats.Option{nats.Name("NATS Showtimes")}

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
		showtimes: &mongodb.ShowTimeModel{
			C: client.Database(*mongoDatabse).Collection("showtimes"),
		},
		conn: nconn,
	}

	nconn.QueueSubscribe(*topic+".list", *queueName+"_list", func(msg *nats.Msg) {

		infoLog.Println("...QueueSubscribe called - topic", *topic)

		bShowtimes, err := app.reply_allShowtimes()
		if err == nil {
			msg.Respond(bShowtimes)
			return
		}

		app.errorLog.Println("...Error from reply all showtimes,", err)
	})

	nconn.QueueSubscribe(*topic+".get", *queueName+"_get", func(msg *nats.Msg) {

		infoLog.Println("...QueueSubscribe called - topic", *topic+".get")
		app.infoLog.Printf("...Subject: %s  Data: %s", msg.Subject, string(msg.Data))
		bShowtime, err := app.reply_getShowtime(string(msg.Data))
		if err == nil {
			msg.Respond(bShowtime)
			return
		}

		app.errorLog.Println("...Error from reply get showtime,", err)
	})

	nconn.QueueSubscribe(*topic+".add", *queueName+"_add", func(msg *nats.Msg) {

		infoLog.Println("...QueueSubscribe called - topic", *topic+".add")
		app.infoLog.Printf("...Subject: %s  Data: %s", msg.Subject, string(msg.Data))
		bShowTime, err := app.reply_addShowTime(string(msg.Data))
		if err == nil {
			msg.Respond(bShowTime)
			return
		}

		app.errorLog.Println("...Error from reply add showtime,", err)
	})

	nconn.QueueSubscribe(*topic+".delete", *queueName+"_delete", func(msg *nats.Msg) {

		infoLog.Println("...QueueSubscribe called - topic", *topic+".delete")
		app.infoLog.Printf("...Subject: %s  Data: %s", msg.Subject, string(msg.Data))
		err := app.reply_deleteShowTime(string(msg.Data))
		if err == nil {
			msg.Respond([]byte(fmt.Sprintf("ShowTime %s was deleted", msg.Data)))
			return
		}
		app.errorLog.Println("...Error from reply delete showtime,", err)
	})
	nconn.Flush()

	if err := nconn.LastError(); err != nil {
		log.Fatal(err)
	}

	log.Printf("Listening on topic: %s", *topic)

	// Initialize a new http.Server struct.
	//serverURI := fmt.Sprintf("%s:%d", *serverAddr, *serverPort)
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
