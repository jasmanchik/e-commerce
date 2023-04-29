package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jasmanchik/garage-sale/schema"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	log.Printf("main: started")
	defer log.Println("main: completed")

	db, err := openDB()
	if err != nil {
		log.Fatalf("main: unable to connect to database: %v", err)
	}
	defer db.Close()

	flag.Parse()
	log.Println(flag.Args())
	if len(flag.Args()) > 0 {
		switch flag.Args()[0] {
		case "migrate":
			if err := schema.Migrate(db); err != nil {
				log.Fatalf("migrate applying error: %s", err)
			} else {
				log.Println("migration complete")
			}
			return
		case "seed":
			if err := schema.Seed(db); err != nil {
				log.Fatalf("seeding error: %s", err)
			} else {
				log.Println("seeding complete")
			}
			return
		}
	}

	api := http.Server{
		Addr:         ":8000",
		Handler:      http.HandlerFunc(ListProducts),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	serverErrors := make(chan error, 1)

	go func() {
		log.Printf("main : API listening on %s", api.Addr)
		serverErrors <- api.ListenAndServe()
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		log.Fatalf("error: listening and serving %s", err)
	case <-shutdown:
		log.Println("main: Start shutdown")

		const timeout = 5 * time.Second
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		err := api.Shutdown(ctx)
		if err != nil {
			log.Printf("main : Graceful shutdown did not complete in %v : %v", timeout, err)
			err := api.Close()
			if err != nil {
				log.Fatalf("main: could not stop server gracefully : %v", err)
			}
		}
	}
}

func openDB() (*sqlx.DB, error) {
	q := url.Values{}
	q.Set("sslmode", "disable")
	q.Set("timezone", "utc")

	u := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD")),
		Host:     os.Getenv("DB_HOST"),
		Path:     os.Getenv("DB_NAME"),
		RawQuery: q.Encode(),
	}
	log.Println(u.String())
	return sqlx.Open("postgres", u.String())
}

type Product struct {
	Name     string `json:"name"`
	Cost     int    `json:"cost"`
	Quantity int    `json:"quantity"`
}

func ListProducts(w http.ResponseWriter, r *http.Request) {
	list := make([]Product, 0)
	list = append(list, Product{Name: "Comic Books", Cost: 50, Quantity: 42})
	list = append(list, Product{Name: "McDonalds Toys", Cost: 75, Quantity: 120})
	list = append(list, Product{Name: "Big Wheels", Cost: 500, Quantity: 2})

	data, err := json.Marshal(list)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("ListProducts: error marshalling data: %s", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(data); err != nil {
		log.Printf("ListProducts: write response: %s", err)
	}
}
