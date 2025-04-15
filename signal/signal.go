package signal

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

func GracefulHttpServerWithShutdown() {
	// Start http server in separate goroutine, so we can handle shutdown gracefully
	router := mux.NewRouter()

	// Endpoint for Cluster Alive / Readiness Probes
	for _, p := range []string{"/health", "/health.json"} {
		router.HandleFunc(p, HealthHandler).Methods("GET")
	}

	server := http.Server{Addr: fmt.Sprintf(":%d", 0000), Handler: router}
	go func() {
		log.Info().Msg("Dispatcher and router initialized, launching http server")
		log.Fatal().Err(server.ListenAndServe())
	}()

	// Listen for INT and TERM signals, shutdown workers and http server gracefully
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	sig := <-signalChan // wait for shutdown signal
	log.Info().Msgf("Received signal %v, prepare to say goodbye", sig)
	// shutdown http server first, so we avoid further incoming requests
	if err := server.Shutdown(context.Background()); err != nil {
		log.Warn().Err(err)
	}
	// Initiate dispatcher shutdown, which will wait for all workers to finish work and close
	os.Exit(0)
}

func HealthHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	sResponse, err := json.Marshal(map[string]interface{}{
		"sResponse":  "up",
		"info":       fmt.Sprintf("%s is healthy", "application"),
		"time":       time.Now().Format(time.RFC3339),
		"goroutines": runtime.NumGoroutine(),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if _, err := w.Write(sResponse); err != nil {
		log.Error().Msg(err.Error())
	}
}
