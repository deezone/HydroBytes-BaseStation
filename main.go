package main

import (
	// Standard packages
	"context"   // https://golang.org/pkg/context/
	"fmt"       // https://golang.org/pkg/fmt/
	"log"       // https://golang.org/pkg/log/
	"math/rand"
	"net/http"  // https://golang.org/pkg/net/http/
	"os"
	"os/signal" // https://golang.org/src/os/signal/doc.go
	"syscall"
	"time"
)

// Main entry point for program.
func main() {

	// Define constants that includes their type to ensure a "kind" decarlation is not used resulting in much
	// larger memory space use. See https://education.ardanlabs.com/courses/take/ultimate-syntax/lessons/13570526-constants-pt-2
	const (
		readTimeout time.Duration   = 5 * time.Second
		writeTimeout time.Duration  = 5 * time.Second
		cancelTimeout time.Duration = 5 * time.Second
	)

	// =========================================================================
	// App Starting

	log.Printf("main : Started")
	defer log.Println("main : Completed")

	// =========================================================================
	// Start API Service

	/**
	 * Convert the Echo function to a type that implements http.Handler
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
		Addr:         "localhost:8000",
		Handler:      http.HandlerFunc(Echo),
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	}

	// Make a channel to listen for errors coming from the listener. Use a buffered channel so the goroutine can exit
	// if we don't collect this error.
	//
	serverErrors := make(chan error, 1)

	// Start a server listening ("bind") on port 8000 and responding using handler called Echo() as defined below
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
		log.Fatalf("error: listening and serving: %s", err)

	case <-shutdown:
		log.Println("main : Start shutdown")

		/**
		 * Give outstanding requests a deadline for completion. Context defines the Context type, which carries
		 * deadlines, cancellation signals, and other request-scoped values across API boundaries and between
		 * processes.
		 * https://golang.org/pkg/context/
		 */
		ctx, cancel := context.WithTimeout(context.Background(), cancelTimeout)
		defer cancel()

		// Asking listener to shutdown and load shed.
		err := api.Shutdown(ctx)
		if err != nil {
			log.Printf("main : Graceful shutdown did not complete in %v : %v", cancelTimeout, err)
			err = api.Close()
		}

		if err != nil {
			log.Fatalf("main : could not stop server gracefully : %v", err)
		}
	}
}

/**
 * Echo is a basic HTTP Handler.
 * Also know as a "controller" in the MVC pattern
 *
 * The parameters must follow the signature defined in the http.HandlerFunc() adapter
 * (https://golang.org/src/net/http/server.go?s=97511:97566#L2034) to convert this method into a HTTP handler type.
 *
 * If you open localhost:8000 in your browser, you may notice double requests being made. This happens because the
 * browser sends a request in the background for a website favicon. More the reason to use Postman to test!
 */
func Echo(w http.ResponseWriter, r *http.Request) {

	// Print a random number at the beginning and end of each request.
	n := rand.Intn(1000)
	log.Println("start", n)
	defer log.Println("end", n)

	// Simulate a long-running request.
	time.Sleep(3 * time.Second)

	fmt.Fprintf(w, "You asked to %s %s\n", r.Method, r.URL.Path)
}
