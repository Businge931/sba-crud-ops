package main

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/Businge931/api-gateway/internal/app/handlers"
	"github.com/Businge931/api-gateway/internal/core/domain"

	log "github.com/sirupsen/logrus"
)

type application struct {
	config      config
	oddsHandler *handlers.OddsHandler
}

type config struct {
	addr string
	db   dbConfig
	env  string
}

type dbConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

func (app *application) mount() http.Handler {
	r := mux.NewRouter()

	// Add middleware
	r.Use(mux.CORSMethodMiddleware(r))
	r.Use(loggingMiddleware)
	r.Use(recoveryMiddleware)

	// API routes
	v1 := r.PathPrefix("/v1").Subrouter()

	// Health check
	v1.HandleFunc("/health", app.healthCheckHandler).Methods("GET")

	// Odds routes
	odds := v1.PathPrefix("/odds").Subrouter()
	odds.HandleFunc("/create", app.oddsHandler.CreateOddsHandler).Methods("POST")
	odds.HandleFunc("/read", app.oddsHandler.ReadOddsHandler).Methods("GET")
	odds.HandleFunc("/update", app.oddsHandler.UpdateOddsHandler).Methods("PUT")
	odds.HandleFunc("/delete", app.oddsHandler.DeleteOddsHandler).Methods("DELETE")

	return r
}

func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status":  "ok",
		"env":     app.config.env,
		"version": version,
	}

	if err := domain.WriteJSON(w, http.StatusOK, data); err != nil {
		domain.InternalServerError(w, r, err)
	}

	// w.WriteHeader(http.StatusOK)
	// w.Write([]byte("OK"))
}

func (app *application) run(mux http.Handler) error {
	srv := http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	log.Printf("Server has started at %s", app.config.addr)

	return srv.ListenAndServe()
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		next.ServeHTTP(w, r)
	})
}

func recoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic: %+v", err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
