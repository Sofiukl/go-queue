package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/sofiukl/go-queue/dispatcher"
	api "github.com/sofiukl/go-queue/restapi"
)

// NWorkers - no of workers
// HTTPAddr - Host details
var (
	NWorkers = flag.Int("n", 4, "The number of workers to start")
	HTTPAddr = flag.String("http", "127.0.0.1:8000", "Address to listen for HTTP requests on")
)

func main() {
	flag.Parse()

	fmt.Println("Starting the dispatcher")
	dispatcher.StartDispatcher(*NWorkers)

	fmt.Println("Registering the work receiver")
	http.HandleFunc("/work", api.ReceiveWork)

	// Start the HTTP server!
	fmt.Println("HTTP server listening on", *HTTPAddr)
	if err := http.ListenAndServe(*HTTPAddr, nil); err != nil {
		fmt.Println(err.Error())
	}
}
