package main

import (
	"context"
	"errors"
	"fmt"
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
)

var commitHash string

func main() {
	log.Printf("Starting Backend %s\n", commitHash)

	if err := run(); err != nil {
		log.Println(err)
	}
}

func run() error {
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

	r.Mount("/data", data.NewService())
	r.Mount("/mesh-node", mesh_node.NewService())

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

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	select {
	case err := <-srvErrChan:
		return err
	case sig := <-signalChan:
		log.Printf("received %s\n", sig)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("shutting down server: %w", err)
	}

	return nil
}
