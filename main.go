package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

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
	r.Use(middleware.Recoverer)
	r.Use(middleware.Heartbeat("/"))
	// r.Mount("/debug", middleware.Profiler())

	r.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World"))
	})

	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		fmt.Printf("%s %s\n", method, route)
		return nil
	}

	if err := chi.Walk(r, walkFunc); err != nil {
		fmt.Printf("Logging err: %s\n", err.Error())
	}

	srv := &http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: r,
	}

	// Clean up
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	srv.Shutdown(ctx)
	log.Println("shutting down server")
	os.Exit(0)
}
