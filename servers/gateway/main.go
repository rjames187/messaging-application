package main

import (
	"database/sql"
	"log"
	"messaging-application/servers/gateway/handlers"
	"messaging-application/servers/gateway/models/users"
	"messaging-application/servers/gateway/sessions"
	"net/http"
	"os"

	"github.com/redis/go-redis/v9"
)

func main() {
	ADDR := os.Getenv("ADDR")
	if len(ADDR) == 0 {
		ADDR = ":443"
	}

	TLSCERT := os.Getenv("TLSCERT")
	if len(TLSCERT) == 0 {
		log.Fatal("No TLSCERT environment variable found")
	}

	TLSKEY := os.Getenv("TLSKEY")
	if len(TLSKEY) == 0 {
		log.Fatal("No TLSKEY environment variable found")
	}

	SESSIONKEY := os.Getenv("SESSIONKEY")
	if len(SESSIONKEY) == 0 {
		log.Fatal("No SESSIONKEY environment variable found")
	}

	REDISADDR := os.Getenv("REDISADDR")
	if len(REDISADDR) == 0 {
		log.Fatal("No REDISADDR environment variable found")
	}

	DSN := os.Getenv("DSN")
	if len(DSN) == 0 {
		log.Fatal("No DSN environment variable found")
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr: REDISADDR,
	})
	redisStore := sessions.NewRedisStore(redisClient, "1h")

	db, err := sql.Open("mysql", DSN)
	if err != nil {
		log.Fatalf("error opening db: %v", err)
	}
	mysqlStore, err := users.NewMySQLStore(db)
	if err != nil {
		log.Fatalf("error creating mysql store: %v", err)
	}

	hctx := handlers.NewHandlerContext(SESSIONKEY, &redisStore, &mysqlStore)

	mux := http.NewServeMux()
	mux.HandleFunc("/v1/summary", handlers.SummaryHandler)
	mux.HandleFunc("/v1/users", hctx.UsersHandler)
	mux.HandleFunc("/v1/users/{UserID}", hctx.SpecificUserHandler)
	mux.HandleFunc("/v1/sessions", hctx.SessionsHandler)
	mux.HandleFunc("/v1/sessions/{SessionID}", hctx.SpecificSessionHandler)

	handler := handlers.NewCORSHandler(mux)

	log.Printf("server is listening at %s...", ADDR)
	log.Fatal(http.ListenAndServeTLS(ADDR, TLSCERT, TLSKEY, handler))
}
