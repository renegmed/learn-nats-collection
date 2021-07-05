package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	nats "github.com/nats-io/nats.go"
)

func usage() {
	log.Printf("Usage: publisher [-s server] \n")
	flag.PrintDefaults()
}

func showUsageAndExit(exitcode int) {
	usage()
	os.Exit(exitcode)
}

func main() {

	var serversUrl = flag.String("s", nats.DefaultURL, "The nats server URLs (separated by comma)")
	var streamName = flag.String("streamname", "", "Stream name")
	var streamSubject = flag.String("topic", "", "Stream subject")
	var serverPort = flag.String("port", "", "This server port")
	var showHelp = flag.Bool("h", false, "Show help message")

	log.SetFlags(0)
	flag.Usage = usage
	flag.Parse()

	if *showHelp {
		showUsageAndExit(0)
	}

	log.Println("...nats server URL", *serversUrl)

	// Connect to NATS
	nc, err := nats.Connect(*serversUrl)
	if err != nil {
		log.Println("...Error on nats connection,", err)
		log.Fatal(err)
	}

	// Creates JetStreamContext
	js, err := nc.JetStream()
	checkErr(err)

	log.Println("...Got nats jet stream connection")

	// Creates stream
	err = createStream(js, *streamName, []string{*streamSubject})
	//err = createStream(js, streamName, []string{"ORDER.created", "ORDER.cancelled"})
	checkErr(err)

	log.Println("...stream created with stream name -", streamName, " and subjects:\n\t", streamSubject)

	srv := server{
		js: js,
	}

	// Serve HTTP
	r := mux.NewRouter()
	r.HandleFunc("/publish", srv.HandlePublishMessage)

	log.Printf("Starting HTTP server on '%s'", *serverPort)

	err = http.ListenAndServe(fmt.Sprintf(":%s", *serverPort), r)
	checkErr(err)

}

// createStream creates a stream by using JetStreamContext
func createStream(js nats.JetStreamContext, streamName string, streamSubjects []string) error {
	// Check if the ORDERS stream already exists; if not, create it.
	stream, err := js.StreamInfo(streamName)
	if err != nil {
		log.Println(err)
	}

	if stream == nil {
		log.Printf("creating stream %q and subjects %q", streamName, streamSubjects)
		_, err = js.AddStream(&nats.StreamConfig{
			Name:     streamName,
			Subjects: streamSubjects,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
