package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	memdb "github.com/hashicorp/go-memdb"
	nats "github.com/nats-io/nats.go"
)

type product struct {
	Name       string `json:"name"`
	SKU        string `json:"sku"`
	StockCount int    `json:"stock_count"`
}

var usageStr = `
Usage: product-read [options] <subject> <message>

Options: 
	-port         This server port number e.g. 8080
	-s,           NATS stream server URL(s)
	-streamname   NATS stream name
	-topic        NATS stream topic
`

func usage() {
	log.Printf("%s\n", usageStr)
	os.Exit(0)
}

type JetStreamPublisher struct {
	*Connector
	StreamName string
	Topic      string
}

type Server struct {
	Jsp *JetStreamPublisher
}

var schema *memdb.DBSchema
var db *memdb.MemDB

func init() {

	schema = &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			"product": &memdb.TableSchema{
				Name: "product",
				Indexes: map[string]*memdb.IndexSchema{
					"id": &memdb.IndexSchema{
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "SKU"},
					},
				},
			},
		},
	}

	err := schema.Validate()
	if err != nil {
		log.Fatal(err)
	}

	db, err = memdb.NewMemDB(schema)
	if err != nil {
		log.Fatal(err)
	}

	txn := db.Txn(true)

	if err := txn.Insert("product", product{"Test1", "ABC232323", 100}); err != nil {
		log.Fatal(err)
	}

	if err := txn.Insert("product", product{"Test2", "ABC883388", 100}); err != nil {
		log.Fatal(err)
	}

	txn.Commit()

}

func main() {
	var port = flag.String("port", "8080", "This server port address")
	var serversUrl = flag.String("s", nats.DefaultURL, "The nats server URLs (separated by comma)")
	var streamName = flag.String("streamname", "", "Stream name")
	var streamSubject = flag.String("topic", "", "Stream subject")

	flag.Usage = usage
	flag.Parse()

	conn, err := NewConnector(*serversUrl)
	if err != nil {
		log.Fatal(err)
	}

	jsPublisher := JetStreamPublisher{
		Connector:  conn,
		StreamName: *streamName,
		Topic:      *streamSubject,
	}

	err = jsPublisher.createStream()
	if err != nil {
		log.Println("...Error on create stream.", err)
		log.Fatal(err)
	}

	server := Server{Jsp: &jsPublisher}

	http.DefaultServeMux.HandleFunc("/product", server.ProductsHandler)
	http.DefaultServeMux.HandleFunc("/stock", server.StockHandler)

	log.Println("Starting product write service on port 8080")
	log.Fatal(http.ListenAndServe(":"+*port, http.DefaultServeMux))
}

func (jsp *JetStreamPublisher) createStream() error {
	// Check if the streamname already exists; if not, create it.
	stream, err := jsp.Connector.Js.StreamInfo(jsp.StreamName)
	if err != nil {
		log.Printf("...Error on getting stream info %s\n\t%v.", jsp.StreamName, err)
	}

	if stream == nil {
		log.Printf("creating stream %q and subjects %q", jsp.StreamName, jsp.Topic)
		_, err = jsp.Connector.Js.AddStream(&nats.StreamConfig{
			Name:     jsp.StreamName,
			Subjects: []string{jsp.Topic},
		})
		if err != nil {
			log.Printf("...Error on adding new stream %s\n\t%v", jsp.StreamName, err)
			return err
		}
	}
	return nil
}

func (jsp *JetStreamPublisher) publishMessage(prod *product) error {

	msgJSON, err := json.Marshal(prod)
	if err != nil {
		log.Printf("Failed to marshal product %s\n\t%v", prod.Name, err)
		return err
	}

	_, err = jsp.Connector.Js.Publish(jsp.Topic, msgJSON)

	return err

}

func (s *Server) ProductsHandler(rw http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		s.InsertProduct(rw, r)
	}
}

func (s *Server) StockHandler(rw http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	txn := db.Txn(false)
	obj, err := txn.First("product", "id", id)
	if err != nil {
		rw.WriteHeader(http.StatusNotFound)
		return
	}

	p := obj.(product)
	fmt.Fprintf(rw, `{"quantity": %v}`, p.StockCount)
}

func (s *Server) InsertProduct(rw http.ResponseWriter, r *http.Request) {
	log.Println("/insert handler called")

	p := &product{}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	err = json.Unmarshal(data, p)
	if err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	txn := db.Txn(true)
	if err := txn.Insert("product", p); err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	txn.Commit()

	s.Jsp.publishMessage(p)
	if err != nil {
		log.Println("...Error on publis message", err)
		rw.WriteHeader(http.StatusInternalServerError)
	}

}
