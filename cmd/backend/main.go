package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/mdma-backend/mdma-backend/internal/api/account"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/mdma-backend/mdma-backend/api"
	"github.com/mdma-backend/mdma-backend/internal/api/data"
	"github.com/mdma-backend/mdma-backend/internal/api/mesh_node"
	"github.com/mdma-backend/mdma-backend/internal/pkg/storage/postgres"
)

const (
	envVarPrefix = "MDMA_"
)

var (
	commitHash  = "dev-build"
	databaseDSN = "postgres://postgres:postgres@localhost/postgres?sslmode=disable&connect_timeout=3"
)

func envString(name, value string) string {
	if v := os.Getenv(envVarPrefix + name); v != "" {
		return v
	}
	return value
}

func initEnvVars() {
	databaseDSN = envString("DATABASE_DSN", databaseDSN)
}

func init() {
	initEnvVars()
}

func main() {
	log.Printf("starting backend %s\n", commitHash)

	if err := run(); err != nil {
		log.Println(err)
	}

	log.Println("shutdown complete\nbye bye <3")
}

func run() error {
	db, err := postgres.New(databaseDSN)
	if err != nil {
		return fmt.Errorf("connecting to postgres: %w", err)
	}

	r := chi.NewRouter()
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"*"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: false,
	}))

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/docs", http.StatusFound)
	})

	docsPath := "/docs"
	openAPIPath := docsPath + "/swagger.yaml"
	r.Handle(docsPath, api.SwaggerUIHandler(api.SwaggerUIOpts{
		Path:    docsPath,
		SpecURL: openAPIPath,
		Title:   "FoREST API Docs",
	}))
	r.Get(openAPIPath, api.SwaggerSpecsHandlerFunc())

	r.Mount("/data", data.NewService(db))
	r.Mount("/mesh-node", mesh_node.NewService())
	r.Mount("/accounts", account.NewService(db))

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	srvErrChan := make(chan error, 1)
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			srvErrChan <- err
		}
	}()

	log.Println("backend is online")

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	select {
	case err := <-srvErrChan:
		return err
	case sig := <-signalChan:
		log.Printf("shutting down; received %s\n", sig)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("shutting down server: %w", err)
	}

	return nil
}
