package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Version is set at build time via ldflags
var Version = "dev"

func main() {
	// Load .env file if present (ignored in production where env vars are set directly)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// CLI flags
	showVersion := flag.Bool("version", false, "Print version and exit")
	port := flag.Int("port", 0, "Port to listen on (overrides DS2API_PORT env var)")
	flag.Parse()

	if *showVersion {
		fmt.Printf("ds2api version %s\n", Version)
		os.Exit(0)
	}

	// Determine listen port
	listenPort := resolvePort(*port)

	// Build and start the HTTP server
	router := newRouter()
	addr := fmt.Sprintf(":%d", listenPort)

	log.Printf("ds2api %s starting on %s", Version, addr)

	server := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server error: %v", err)
	}
}

// resolvePort determines the port to listen on.
// Priority: CLI flag > DS2API_PORT env var > default (3000)
func resolvePort(flagPort int) int {
	if flagPort > 0 {
		return flagPort
	}

	if envPort := os.Getenv("DS2API_PORT"); envPort != "" {
		p, err := strconv.Atoi(envPort)
		if err != nil {
			log.Printf("Invalid DS2API_PORT value %q, falling back to default 3000", envPort)
			return 3000
		}
		return p
	}

	return 3000
}
