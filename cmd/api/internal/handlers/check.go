package handlers

import (
	"context"
	// Core packages
	"net/http"

	// Internal packages
	"github.com/deezone/HydroBytes-BaseStation/internal/platform/database"
	"github.com/deezone/HydroBytes-BaseStation/internal/platform/web"

	// Third-party packages
	"github.com/jmoiron/sqlx"
)

// Check provides support for orchestration health checks.
type Check struct {
	db *sqlx.DB
}

// Health validates the service is healthy and ready to accept requests.
func (c *Check) Health(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	var health struct {
		Status string `json:"status"`
	}

	// Check if the database is ready.
	if err := database.StatusCheck(ctx, c.db); err != nil {

		// If the database is not ready we will tell the client and use a 500
		// status. Do not respond by just returning an error because further up in
		// the call stack will interpret that as an unhandled error.
		health.Status = "db not ready"
		return web.Respond(ctx, w, health, http.StatusInternalServerError)
	}

	health.Status = "ok"
	return web.Respond(ctx, w, health, http.StatusOK)
}
