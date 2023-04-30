package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ardanlabs/conf"
	"github.com/jasmanchik/garage-sale/cmd/sales-api/internal/handlers"
	"github.com/jasmanchik/garage-sale/internal/platform/database"
	"github.com/pkg/errors"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	log := log.New(os.Stdout, "SALES : ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

	log.Printf("main: started")
	defer log.Println("main: completed")

	var cfg struct {
		Web struct {
			Address      string        `conf:"default::8000"`
			ReadTimeout  time.Duration `conf:"default:5s"`
			WriteTimeout time.Duration `conf:"default:5s"`
			ShutdownTime time.Duration `conf:"default:5s"`
		}
		DB struct {
			User       string `conf:"default:db"`
			Pass       string `conf:"default:db,noprint"`
			Host       string `conf:"default:db"`
			Name       string `conf:"default:db"`
			DisableTLS bool   `conf:"default:true"`
		}
	}

	if err := conf.Parse(os.Args[1:], "SALES", &cfg); err != nil {
		if err == conf.ErrHelpWanted {
			usage, err := conf.Usage("SALES", &cfg)
			if err != nil {
				errors.Wrap(err, "generating config usage")
			}
			fmt.Println(usage)
			return nil
		}
		errors.Wrap(err, "parsing config")
	}

	out, err := conf.String(&cfg)
	if err != nil {
		errors.Wrap(err, "generating config for output")
	}
	log.Printf("main: Config :\n%v\n", out)
	db, err := database.Open(database.Config{
		User:       cfg.DB.User,
		Password:   cfg.DB.Pass,
		Host:       cfg.DB.Host,
		Name:       cfg.DB.Name,
		DisableTLS: cfg.DB.DisableTLS,
	})
	if err != nil {
		errors.Wrap(err, "unable to connect to database")
	}
	defer db.Close()

	api := http.Server{
		Addr:         cfg.Web.Address,
		Handler:      handlers.Routes(log, db),
		ReadTimeout:  cfg.Web.ReadTimeout,
		WriteTimeout: cfg.Web.WriteTimeout,
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
		errors.Wrap(err, "listening and serving")
	case <-shutdown:
		log.Println("main: Start shutdown")
		ctx, cancel := context.WithTimeout(context.Background(), cfg.Web.ShutdownTime)
		defer cancel()

		err := api.Shutdown(ctx)
		if err != nil {
			log.Printf("main : Graceful shutdown did not complete in %v : %v", cfg.Web.ShutdownTime, err)
			err = api.Close()
		}

		if err != nil {
			errors.Wrap(err, "graceful shutdown")
		}
	}

	return nil
}
