package main

import (
	// Core packages
	"context"
	"crypto/rsa"
	_ "expvar" // Register the expvar handlers
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	_ "net/http/pprof" // Register the pprof handlers
	"os"
	"os/signal"
	"syscall"
	"time"

	// Internal packages
	"github.com/deezone/HydroBytes-BaseStation/internal/platform/auth"
	"github.com/deezone/HydroBytes-BaseStation/cmd/api/internal/handlers"
	"github.com/deezone/HydroBytes-BaseStation/internal/platform/conf"
	"github.com/deezone/HydroBytes-BaseStation/internal/platform/database"

	// Third-party packages
	"contrib.go.opencensus.io/exporter/zipkin"
	"github.com/dgrijalva/jwt-go"
	openzipkin "github.com/openzipkin/zipkin-go"
	zipkinHTTP "github.com/openzipkin/zipkin-go/reporter/http"
	"github.com/pkg/errors"
	"go.opencensus.io/trace"
)

// Main entry point for program.
func main() {

	// Only make termination calls (log.Fatalf()) in main() to allow all defers to complete before shutdown in the
	// case of an error.
	if err := run(); err != nil {
		log.Println("shutting down", "error:", err)
		os.Exit(1)
	}
}

// Main application logic.
func run() error {

	// Define constants that includes their type to ensure a "kind" decarlation is not used resulting in much
	// larger memory space use. See https://education.ardanlabs.com/courses/take/ultimate-syntax/lessons/13570526-constants-pt-2
	// ex: readTimeout time.Duration   = 5 * time.Second

	// =========================================================================
	// Logging

	// Use "shadowing" to override the global log package value. See https://golang.org/src/log/log.go#L37 for possible
	// bit values to manage output
	log := log.New(os.Stdout, "STATIONS API : ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

	// =========================================================================
	// Configuration

	var cfg struct {
		Web struct {
			Address         string        `conf:"default:localhost:8000"`
			Debug           string        `conf:"default:localhost:6060"`
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
		Auth struct {
			KeyID          string `conf:"default:1"`
			PrivateKeyFile string `conf:"default:private.pem"`
			Algorithm      string `conf:"default:RS256"`
		}
		Trace struct {
			URL         string  `conf:"default:http://localhost:9411/api/v2/spans"`
			Service     string  `conf:"default:station-api"`
			Probability float64 `conf:"default:1"` // reduce this value to increase sampling - 1 = 100% of requests
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
	// Initialize authentication support

	authenticator, err := createAuth(
		cfg.Auth.PrivateKeyFile,
		cfg.Auth.KeyID,
		cfg.Auth.Algorithm,
	)
	if err != nil {
		return errors.Wrap(err, "constructing authenticator")
	}

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
	// Start Tracing Support

	closer, err := registerTracer(
		cfg.Trace.Service,
		cfg.Web.Address,
		cfg.Trace.URL,
		cfg.Trace.Probability,
	)
	if err != nil {
		return err
	}
	defer closer()

	// =========================================================================
	// Start Debug Service

	/**
	 * GET /debug/pprof - Added to the default mux by importing the net/http/pprof package.
	 * GET /debug/vars - Added to the default mux by importing the expvar package.
	 *
	 * Not concerned with shutting this down when the application is shutdown.
	 */
	go func() {
		log.Println("debug service listening on", cfg.Web.Debug)
		err := http.ListenAndServe(cfg.Web.Debug, http.DefaultServeMux)
		log.Println("debug service closed", err)
	}()

	// =========================================================================
	// Start API Service

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
		Handler:      handlers.API(db, log, authenticator),
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

func createAuth(privateKeyFile, keyID, algorithm string) (*auth.Authenticator, error) {

	keyContents, err := ioutil.ReadFile(privateKeyFile)
	if err != nil {
		return nil, errors.Wrap(err, "reading auth private key")
	}

	key, err := jwt.ParseRSAPrivateKeyFromPEM(keyContents)
	if err != nil {
		return nil, errors.Wrap(err, "parsing auth private key")
	}

	public := auth.NewSimpleKeyLookupFunc(keyID, key.Public().(*rsa.PublicKey))

	return auth.NewAuthenticator(key, keyID, algorithm, public)
}

func registerTracer(service, httpAddr, traceURL string, probability float64) (func() error, error) {
	localEndpoint, err := openzipkin.NewEndpoint(service, httpAddr)
	if err != nil {
		return nil, errors.Wrap(err, "creating the local zipkinEndpoint")
	}
	reporter := zipkinHTTP.NewReporter(traceURL)

	trace.RegisterExporter(zipkin.NewExporter(reporter, localEndpoint))
	trace.ApplyConfig(trace.Config{
		DefaultSampler: trace.ProbabilitySampler(probability),
	})

	return reporter.Close, nil
}
