package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/birmsi/test-signer/internal/helpers"
	"github.com/birmsi/test-signer/internal/signatures/api"
	"github.com/birmsi/test-signer/internal/signatures/repository"
	"github.com/birmsi/test-signer/internal/signatures/service"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type DatabaseConfig struct {
	connectionString string
	username         string
	password         string
	name             string
}

type ApplicationConfiguration struct {
	db          DatabaseConfig
	serverPort  string
	environment string
	version     string
	logger      *slog.Logger
}

func main() {
	fmt.Println("Starting Signatures API")

	app := handleEnvVariables()
	app.logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	db, err := app.openDB()
	if err != nil {
		fmt.Println(err)
		return
	}

	err = app.serve(db)
	if err != nil {
		app.logger.Error(err.Error())
		os.Exit(1)
	}

}

func handleEnvVariables() ApplicationConfiguration {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		log.Fatal("Missing PORT from .env")
	} else {
		port = fmt.Sprintf(":%s", port)
	}

	env := os.Getenv("ENVIRONMENT")
	if port == "" {
		log.Fatal("Missing environment from .env")
	}

	version := os.Getenv("VERSION")
	if version == "" {
		log.Fatal("Missing VERSION from .env")
	}

	dbAddress := os.Getenv("DB_ADDRESS")
	if dbAddress == "" {
		log.Fatal("Missing DB_ADDRESS from .env")
	}
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		log.Fatal("Missing DB_NAME from .env")
	}
	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		log.Fatal("Missing DB_USER from .env")
	}
	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		log.Fatal("Missing DB_PASSWORD from .env")
	}

	return ApplicationConfiguration{
		db: DatabaseConfig{
			connectionString: dbAddress,
			username:         dbUser,
			password:         dbPassword,
			name:             dbName,
		},
		serverPort:  port,
		environment: env,
		version:     version,
	}
}

func (app ApplicationConfiguration) openDB() (*sql.DB, error) {

	dsn := fmt.Sprintf("postgres://%s:%s@%s/%s", app.db.username, app.db.password, app.db.connectionString, app.db.name)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (app ApplicationConfiguration) serve(db *sql.DB) error {
	srv := http.Server{
		Addr:         app.serverPort,
		Handler:      app.loadHandlers(db),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	shutdownError := make(chan error)
	go func() {

		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit

		app.logger.Info("shutting down server", "signal", s.String())

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		shutdownError <- srv.Shutdown(ctx)
	}()

	app.logger.Info(fmt.Sprintf("starting %s server on %s", app.environment, srv.Addr))

	if err := srv.ListenAndServe(); err != nil {
		fmt.Println(err.Error())
	}

	err := <-shutdownError
	if err != nil {
		return err
	}

	app.logger.Info("stopped server", "addr", srv.Addr)
	return nil
}

func (app ApplicationConfiguration) loadHandlers(db *sql.DB) http.Handler {

	signaturesRepository := repository.NewSignaturesRepository(*app.logger, db)
	signaturesService := service.NewSignaturesService(*app.logger, signaturesRepository)
	signaturesAPI := api.NewSignaturesAPI(*app.logger, signaturesService)

	mux := http.NewServeMux()
	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		env := helpers.JsonEnvelope{
			"status": "available",
			"system_info": map[string]string{
				"environment": app.environment,
				"version":     app.version,
			},
		}
		err := helpers.WriteJSON(w, http.StatusOK, env, nil)
		if err != nil {
			app.logger.Error(err.Error())
			http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
		}

	})

	signaturesAPI.Handlers(mux)

	return app.loggingMiddleware(mux)
}

var requestIDCounter int
var mu sync.Mutex

func (app ApplicationConfiguration) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		mu.Lock()
		requestIDCounter++
		requestID := requestIDCounter
		mu.Unlock()

		log.Printf(
			"Received %s request %d for %s from %s",
			r.Method,
			requestID,
			r.URL.Path,
			r.RemoteAddr,
		)

		for name, values := range r.Header {
			for _, value := range values {
				log.Printf("Request Header: %s: %s", name, value)
			}
		}

		startTime := time.Now()

		next.ServeHTTP(w, r)

		log.Printf("Request %d responded in %s", requestID, time.Since(startTime))
	})
}
