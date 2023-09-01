package main

import (
	"info-project/servers/gateway/handlers"
	"log"
	"net/http"

	"github.com/rs/cors"
)

func main() {
	addr := ":4000"

	mux := http.NewServeMux()
	mux.HandleFunc("/v1/summary", handlers.SummaryHandler)

	handler := cors.Default().Handler(mux)
	log.Printf("server is listening at %s...", addr)
	log.Fatal(http.ListenAndServe(addr, handler))
}