package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/lib/pq"
)

var db *sql.DB

type message struct {
	Message string `json:string`
}

func main() {
	var err error
	port := os.Getenv("PORT")

	// skip := os.Getenv("skip_env_load")
	// s, _ := strconv.ParseBool(skip)

	// if !s {
	// 	fmt.Println(skip, s)
	// 	err = godotenv.Load()
	// 	if err != nil {
	// 		log.Fatal("Could not load environment variables", err)
	// 	}
	// }
	// fmt.Println("Loaded environment variables successfully!")

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
	r.Use(middleware.Heartbeat("/healthz"))
	// r.Mount("/debug", middleware.Profiler())

	r.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
		msg := message{"Hello World"}
		json.NewEncoder(w).Encode(msg)
	})

	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		if strings.HasPrefix(route, "/debug") {
			return nil
		}
		fmt.Printf("%s %s\n", method, route)
		return nil
	}

	if err := chi.Walk(r, walkFunc); err != nil {
		fmt.Printf("Walk err: %s\n", err.Error())
	}

	srv := &http.Server{
		Addr:    "0.0.0.0:" + port,
		Handler: r,
	}

	// Graceful shutdown
	go func() {
		log.Println("server started on port " + port)
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
