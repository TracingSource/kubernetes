package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

var port = flag.Int("port", 8080, "Port number.")

func main() {
	flag.Parse()
	resourceConsumerHandler := NewResourceConsumerHandler()
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), resourceConsumerHandler))
}
