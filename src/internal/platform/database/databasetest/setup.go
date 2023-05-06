package databasetest

import (
	"testing"
	"time"

	"github.com/jasmanchik/garage-sale/internal/platform/database"
	"github.com/jasmanchik/garage-sale/internal/schema"
	"github.com/jmoiron/sqlx"
)

func Setup(t *testing.T) (*sqlx.DB, func()) {
	t.Helper()

	c := startContainer(t)

	db, err := database.Open(database.Config{
		User:       "postgres",
		Password:   "postgres",
		Host:       c.Host,
		Name:       "postgres",
		DisableTLS: true,
	})
	if err != nil {
		t.Fatalf("opening database connection: %s", err)
	}

	t.Log("waiting for database to be ready")

	var pingError error
	maxAttempts := 10
	for attempts := 0; attempts < maxAttempts; attempts++ {
		pingError = db.Ping()
		if pingError == nil {
			break
		}
		time.Sleep(time.Duration(attempts) * time.Second)
	}

	if pingError != nil {
		t.Fatalf("waiting for database to be ready: %s", pingError)
	}

	if err := schema.Migrate(db); err != nil {
		stopContainer(t, c)
		t.Fatalf("migrating database: %s", err)
	}

	teardown := func() {
		t.Helper()
		err := db.Close()
		if err != nil {
			return // Do not fail test if we can not close database connection.
		}
		stopContainer(t, c)
	}

	return db, teardown
}
