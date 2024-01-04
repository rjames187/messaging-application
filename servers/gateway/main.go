package main

import (
	"info-project/servers/gateway/handlers"
	"log"
	"net/http"
	"os"

	"github.com/rs/cors"
)

func main() {
	ADDR := os.Getenv("ADDR")
	if ADDR == "" {
		ADDR = ":443"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/v1/summary", handlers.SummaryHandler)

	handler := cors.Default().Handler(mux)
	log.Printf("server is listening at %s...", ADDR)
	log.Fatal(http.ListenAndServe(ADDR, handler))
}