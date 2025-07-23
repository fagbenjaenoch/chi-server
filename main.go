package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var db *sql.DB

func main() {
	var err error
	err = godotenv.Load()
	if err != nil {
		log.Fatal("Could not load environment variables")
	}
	fmt.Println("Loaded environment variables successfully!")

	dbConnStr := os.Getenv("DB_CONN")
	db, err = sql.Open("postgres", dbConnStr)
	if err != nil {
		log.Fatal("Could not connect to db: ", err)
	}
	fmt.Println("Connected to DB successfully!")
	defer db.Close()

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World"))
	})

	http.ListenAndServe(":3000", r)
}
