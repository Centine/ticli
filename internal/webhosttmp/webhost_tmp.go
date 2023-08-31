package webhosttmp

// Temporary measure, to be replaced with a docker-based testing setup and eventually a proper web server hosted elsewhere

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	port       = ":12845"
	staticPath = "./scriptbundles"
)

func StartServer() {
	fs := http.FileServer(http.Dir(staticPath))
	http.Handle("/", fs)

	server := &http.Server{
		Addr: port,
	}

	go func() {
		log.Printf("Starting server on %s...", port)
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("Error running server: %v", err)
		}
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	<-signalChan

	log.Println("Executing exit routines...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}

	log.Println("Server has shut down gracefully")
}
