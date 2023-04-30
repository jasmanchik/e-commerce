package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/ardanlabs/conf"
	"github.com/jasmanchik/garage-sale/internal/platform/database"
	"github.com/jasmanchik/garage-sale/internal/schema"
)

func main() {

	var cfg struct {
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
				log.Fatalf("main: generating config usage: %v", err)
			}
			fmt.Println(usage)
			return
		}
		log.Fatalf("main: parsing config: %s", err)
	}

	out, err := conf.String(&cfg)
	if err != nil {
		log.Fatalf("main: generating config for output: %v", err)
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
		log.Fatalf("main: unable to connect to database: %v", err)
	}
	defer db.Close()

	flag.Parse()
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
}
