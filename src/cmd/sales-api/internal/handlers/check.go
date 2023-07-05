package handlers

import (
	"github.com/jasmanchik/garage-sale/internal/platform/database"
	"github.com/jasmanchik/garage-sale/internal/platform/web"
	"github.com/jmoiron/sqlx"
	"net/http"
)

type Check struct {
	DB *sqlx.DB
}

// Health Response with a 200 OK if the service is ready fo traffic
func (c *Check) Health(w http.ResponseWriter, r *http.Request) error {

	var health struct {
		Status string `json:"status"`
	}

	if err := database.StatusCheck(r.Context(), c.DB); err != nil {
		health.Status = "db is not ready"
		return web.Response(w, health, http.StatusInternalServerError)
	}

	health.Status = "OK"
	return web.Response(w, health, http.StatusOK)
}
