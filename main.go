package main

import (
	// Standard packages
	"fmt"      // https://golang.org/pkg/fmt/
	"log"      // https://golang.org/pkg/log/
	"net/http" // https://golang.org/pkg/net/http/
)

// Main entry point for program
func main() {

	/**
	 * Convert the Echo function to a type that implements http.Handler
	 * See https://education.ardanlabs.com/courses/take/ultimate-syntax/lessons/13570357-type-conversions for details
	 * on "types" and conversions
     *
     * Details on using HandlerFunc which is an adapter
     * "to allow the use of ordinary functions as HTTP handlers"
     * https://golang.org/pkg/net/http/#HandlerFunc
	 */
	h := http.HandlerFunc(Echo)

	// Start a server listening ("bind") on port 8000 and responding using handler called Echo() as defined below
	// https://golang.org/pkg/net/http/#Server.ListenAndServe
	log.Println("Listening on localhost:8000")
	if err := http.ListenAndServe("localhost:8000", h); err != nil {
		log.Fatalf("error: listening and serving: %s", err)
	}
}

/**
 * Echo is a basic HTTP Handler.
 * Also know as a "controller" in the MVC pattern
 *
 * The parameters must follow the signature defined in ?? in order for http.HandlerFunc() to convert the method into
 * a HTTP handler type
 */
func Echo(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "You asked to %s %s\n", r.Method, r.URL.Path)
}
