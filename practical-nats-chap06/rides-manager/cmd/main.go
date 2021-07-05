package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"

	"practical-nats/rides-manager/kit"
	"practical-nats/rides-manager/server"
)

func main() {

	var (
		showHelp    bool
		showVersion bool
		natsServers string
	)
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: rides-manager [options...]\n\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n")
	}

	// Setup default flags
	flag.BoolVar(&showHelp, "help", false, "Show help")
	flag.BoolVar(&showVersion, "version", false, "Show version")
	flag.StringVar(&natsServers, "s", "", "List of NATS Servers to connect. Do ") // required to pass nats URLs
	flag.Parse()

	switch {
	case showHelp:
		flag.Usage()
		os.Exit(0)
	case showVersion:
		fmt.Fprintf(os.Stderr, "NATS Rider Rides Manager Server v%s\n", server.Version)
		os.Exit(0)
	}
	log.Printf("Starting NATS Rider Rides Manager version %s", server.Version)
	log.Println("nats servers:", natsServers)

	comp := kit.NewComponent("rides-manager")
	err := comp.SetupConnectionToNATS(natsServers)
	if err != nil {
		log.Fatal(err)
	}

	s := server.Server{
		Component: comp,
	}
	err = s.SetupSubscription()
	if err != nil {
		log.Fatal(err)
	}

	runtime.Goexit()
}
