package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/slaven/go-calculator/calcserver"
)

func main() {
	log.Println("[Server] Starting http server")

	// Create context and listen for the interrupt signal
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// Create handler
	calcHandler := calcserver.Create()

	srv := http.Server{
		Addr:         ":8080",
		Handler:      calcHandler,
		ReadTimeout:  time.Duration(30 * time.Second),
		WriteTimeout: time.Duration(30 * time.Second),
		IdleTimeout:  time.Duration(30 * time.Second),
	}

	// Listen on a different Goroutine to avoid blocking
	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("[Server] ListenAndServe: %v", err)
		}
	}()

	// Listen for the interrupt signal
	<-ctx.Done()

	// Cancel context for the interrupt signal
	stop()
	log.Println("[Server] Shutting down gracefully, press Ctrl+C again to force")

	// Perform server shutdown with timeout
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(timeoutCtx); err != nil {
		log.Printf("[Server] Shutdown error: %v", err)
	}
}
