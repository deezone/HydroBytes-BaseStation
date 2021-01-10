package main

import (
	// Core packages
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	// Third-party packages
	"github.com/pkg/errors"

	// Internal packages
	"github.com/deezone/HydroBytes-BaseStation/cmd/api/internal/handlers"
	"github.com/deezone/HydroBytes-BaseStation/internal/platform/conf"
	"github.com/deezone/HydroBytes-BaseStation/internal/platform/database"

)

// Main entry point for program.
func main() {

	// Only make termination calls (log.Fatalf()) in main() to allow all defers to complete before shutdown in the
	// case of an error.
	if err := run(); err != nil {
		log.Printf("error: shutting down: %s", err)
		os.Exit(1)
	}
}

// Main application logic.
func run() error {

	// Define constants that includes their type to ensure a "kind" decarlation is not used resulting in much
	// larger memory space use. See https://education.ardanlabs.com/courses/take/ultimate-syntax/lessons/13570526-constants-pt-2
	// ex: readTimeout time.Duration   = 5 * time.Second

	// =========================================================================
	// Configuration

	var cfg struct {
		Web struct {
			Address         string        `conf:"default:localhost:8000"`
			ReadTimeout     time.Duration `conf:"default:5s"`
			WriteTimeout    time.Duration `conf:"default:5s"`
			ShutdownTimeout time.Duration `conf:"default:5s"`
		}
		DB struct {
			User       string `conf:"default:postgres"`
			Password   string `conf:"default:postgres,noprint"`
			Host       string `conf:"default:localhost"`
			Name       string `conf:"default:postgres"`
			DisableTLS bool   `conf:"default:false"`
		}
	}

	if err := conf.Parse(os.Args[1:], "STATIONS", &cfg); err != nil {
		if err == conf.ErrHelpWanted {
			usage, err := conf.Usage("STATIONS", &cfg)
			if err != nil {
				return errors.Wrap(err, "generating config usage")
			}
			fmt.Println(usage)
			return nil
		}
		return errors.Wrap(err, "parsing config")
	}

	// =========================================================================
	// App Starting

	log.Printf("main : Started")
	defer log.Println("main : Completed")

	out, err := conf.String(&cfg)
	if err != nil {
		return errors.Wrap(err, "generating config for output")
	}
	log.Printf("main : Config :\n%v\n", out)

	// =========================================================================
	// Start Database

	db, err := database.Open(database.Config{
		User:       cfg.DB.User,
		Password:   cfg.DB.Password,
		Host:       cfg.DB.Host,
		Name:       cfg.DB.Name,
		DisableTLS: cfg.DB.DisableTLS,
	})
	if err != nil {
		return errors.Wrap(err, "connecting to db")
	}
	defer db.Close()

	// =========================================================================
	// Start API Service

	// Create copy of service (ss) to allow passing method (ss.List) to map to handler
	stationTypeHandler := handlers.StationTypes{DB: db}

	/**
	 * Convert the ListStationTypes function to a type that implements http.Handler
	 * See https://education.ardanlabs.com/courses/take/ultimate-syntax/lessons/13570357-type-conversions for details
	 * on "types" and conversions.
	 *
	 * Details on using HandlerFunc which is an adapter
	 * "to allow the use of ordinary functions as HTTP handlers"
	 * https://golang.org/pkg/net/http/#HandlerFunc
	 *
	 * Note this is an implementation of the same functionality that http.ListenAndServe() provides but is made
	 * available in a channel.
	 * https://golang.org/src/net/http/server.go?s=97511:97566#L3108
	 */
	api := http.Server{
		Addr:         cfg.Web.Address,
		Handler:      http.HandlerFunc(stationTypeHandler.List),
		ReadTimeout:  cfg.Web.ReadTimeout,
		WriteTimeout: cfg.Web.WriteTimeout,
	}

	// Make a channel to listen for errors coming from the listener. Use a buffered channel so the goroutine can exit
	// if we don't collect this error.
	serverErrors := make(chan error, 1)

	// Start a server listening ("bind") on port 8000 and responding using handler called ListStationTypes()
	// https://golang.org/pkg/net/http/#Server.ListenAndServe
	go func() {
		log.Printf("main : API listening on %s", api.Addr)
		serverErrors <- api.ListenAndServe()
	}()

	/**
	 * Make a channel to listen for an interrupt or terminate signal from the OS. Use a buffered channel because the
	 * signal package requires it. Note the value of "1", this limits the capacity to one to prevent processes from
	 * staying alive if more than one thread is added to the channel which would prevent termination.
	 *
	 * Listening for os.Interrupt and syscall.SIGTERM events.
	 */
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// =========================================================================
	// Shutdown

	// Blocking main and waiting for shutdown.
	// Note two active channels: serverErrors and shutdown as defined above
	select {
	case err := <-serverErrors:
		return errors.Wrap(err, "starting server")

	case <-shutdown:
		log.Println("main : Start shutdown")

		/**
		 * Give outstanding requests a deadline for completion. Context defines the Context type, which carries
		 * deadlines, cancellation signals, and other request-scoped values across API boundaries and between
		 * processes.
		 * https://golang.org/pkg/context/
		 */
		ctx, cancel := context.WithTimeout(context.Background(), cfg.Web.ShutdownTimeout)
		defer cancel()

		// Asking listener to shutdown and load shed.
		err := api.Shutdown(ctx)
		if err != nil {
			log.Printf("main : Graceful shutdown did not complete in %v : %v", cfg.Web.ShutdownTimeout, err)
			err = api.Close()
		}

		if err != nil {
			return errors.Wrap(err, "could not stop server gracefully")
		}
	}

	return nil
}
