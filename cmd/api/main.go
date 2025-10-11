package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Sasidhar222/music-api/internal/data"
	"github.com/Sasidhar222/music-api/internal/handler"
)

func main() {

	connStr := "user=postgres password=mysecretpassword host=localhost port=5432 dbname=postgres sslmode=disable"
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer db.Close()

	// Ping the database to ensure the connection is alive
	err = db.Ping()
	if err != nil {
		log.Fatalf("Unable to ping database: %v\n", err)
	}

	log.Println("Successfully connected to the database!")

	artistModel := &data.ArtistModel{DB: db}
	handlers := &handler.Handlers{Artists: artistModel}

	http.HandleFunc("/status", handlers.StatusHandler)
	http.HandleFunc("/artists", handlers.ArtistsHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	log.Printf("Listening on localhost:%s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
