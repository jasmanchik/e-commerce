package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ardanlabs/conf"
	"github.com/jasmanchik/garage-sale/internal/platform/database"
	"github.com/jasmanchik/garage-sale/internal/schema"
	"github.com/pkg/errors"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	var cfg struct {
		DB struct {
			User       string `conf:"default:db"`
			Pass       string `conf:"default:db,noprint"`
			Host       string `conf:"default:db"`
			Name       string `conf:"default:db"`
			DisableTLS bool   `conf:"default:true"`
		}
		Args conf.Args
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
		errors.Wrap(err, "connecting to db")
	}
	defer db.Close()

	switch cfg.Args.Num(0) {
	case "migrate":
		if err := schema.Migrate(db); err != nil {
			errors.Wrap(err, "migrate applying error")
		} else {
			log.Println("migration complete")
		}
		return nil
	case "seed":
		if err := schema.Seed(db); err != nil {
			errors.Wrap(err, "seeding error")
		} else {
			errors.Wrap(err, "seeding complete")
		}
		return nil
	}
	return nil
}
