package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	memdb "github.com/hashicorp/go-memdb"
	nats "github.com/nats-io/nats.go"
)

type productInsertedEvent struct {
	Name string `json:"name"`
	SKU  string `json:"sku"`
}

type product struct {
	Name        string `json:"name"`
	Code        string `json:"code"`
	LastUpdated string `json:"last_updated"`
}

func (p product) FromProductInsertedEvent(e productInsertedEvent) product {
	p.Name = e.Name
	p.Code = e.SKU
	p.LastUpdated = time.Now().Format(time.RFC3339)

	return p
}

var usageStr = `
Usage: product-read [options] <subject> <message>

Options: 
	-port          This server port number e.g. 8080
	-s,            NATS stream server URL(s)
	-durablename   NATS stream durable name
	-topic         NATS stream topic
`

func usage() {
	log.Printf("%s\n", usageStr)
	os.Exit(0)
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
						Indexer: &memdb.StringFieldIndex{Field: "Code"},
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
}

type Product struct {
	Name       string `json:"name"`
	SKU        string `json:"sku"`
	StockCount int    `json:"stock_count"`
}

func main() {

	log.SetFlags(0)

	var port = flag.String("port", "8080", "This server port address")
	var serversUrl = flag.String("s", nats.DefaultURL, "The nats server URLs (separated by comma)")
	var durableName = flag.String("durablename", "", "Stream durable name")
	var streamSubject = flag.String("topic", "", "Stream subject")

	flag.Usage = usage
	flag.Parse()

	conn, err := NewConnector(*serversUrl)
	if err != nil {
		log.Fatal(err)
	}

	// Create Pull based consumer with maximum 128 inflight.
	// PullMaxWaiting defines the max inflight pull requests.

	log.Println("...Subscriber subject name:", *streamSubject)
	sub, err := conn.Js.PullSubscribe(*streamSubject, *durableName, nats.PullMaxWaiting(128))
	if err != nil {
		log.Println("...Error on pull subscribe,", err)
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1200*time.Second)
	defer cancel()

	batch_size := 2

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}
			msgs, _ := sub.Fetch(batch_size, nats.Context(ctx))
			for _, msg := range msgs {
				msg.Ack()

				var prod Product
				if err := json.Unmarshal(msg.Data, &prod); err != nil {
					log.Println("...Error on unmarshal message", err)
				}

				productMessage(&prod)

				// log.Printf("...Message received\n\tMessageID: %s\n\tMessage Topic: %s\n\tMessage: %s\n",
				// 	message.MessageID, message.Topic, message.Message)

			}
		}

	}()

	http.DefaultServeMux.HandleFunc("/product", getProducts)

	log.Println("...Starting product read service on port 8081")
	log.Fatal(http.ListenAndServe(":"+*port, http.DefaultServeMux))

}

func getProducts(rw http.ResponseWriter, r *http.Request) {
	log.Println("/get handler called")

	txn := db.Txn(false)
	results, err := txn.Get("product", "id")
	if err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Println("...Product found....")

	products := make([]product, 0)
	for {
		obj := results.Next()
		if obj == nil {
			break
		}

		products = append(products, obj.(product))

		log.Println("...Product added to list:\n\t", obj.(product))
	}

	encoder := json.NewEncoder(rw)
	encoder.Encode(products)
}

func productMessage(prod *Product) { //m *nats.Msg) {
	pie := productInsertedEvent{
		Name: prod.Name,
		SKU:  prod.SKU,
	}

	p := product{}.FromProductInsertedEvent(pie)

	txn := db.Txn(true)
	if err := txn.Insert("product", p); err != nil {
		log.Println(err)
		return
	}
	txn.Commit()

	log.Println("Saved product: ", p)
}
