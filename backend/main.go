package main

import (
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"sezzle-calculator/backend/internal/httpapi"
)

func main() {
	addr, err := httpapi.ListenAddr(os.Getenv("PORT"))
	if err != nil {
		log.Fatal(err)
	}

	server := &http.Server{
		Addr:              addr,
		Handler:           httpapi.NewMux(),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	log.Printf("calculator API listening on %s", addr)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}
