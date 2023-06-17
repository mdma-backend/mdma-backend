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

	"github.com/mdma-backend/mdma-backend/internal/api/area"
	"github.com/mdma-backend/mdma-backend/internal/api/mesh_node_update"
	"github.com/mdma-backend/mdma-backend/internal/api/service_account"
	"github.com/mdma-backend/mdma-backend/internal/api/user_account"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
	"github.com/golang-jwt/jwt/v5"
	"github.com/mdma-backend/mdma-backend/api"
	"github.com/mdma-backend/mdma-backend/internal/api/auth"
	"github.com/mdma-backend/mdma-backend/internal/api/data"
	"github.com/mdma-backend/mdma-backend/internal/api/mesh_node"
	"github.com/mdma-backend/mdma-backend/internal/api/role"
	"github.com/mdma-backend/mdma-backend/internal/pkg/storage/postgres"
)

const (
	envVarPrefix = "MDMA_"
)

var (
	commitHash  = "dev-build"
	databaseDSN = "postgres://postgres:postgres@localhost/postgres?sslmode=disable&connect_timeout=3"
	jwtSecret   = "change_me"
)

func envString(name, value string) string {
	if v := os.Getenv(envVarPrefix + name); v != "" {
		return v
	}
	return value
}

func initEnvVars() {
	databaseDSN = envString("DATABASE_DSN", databaseDSN)
	jwtSecret = envString("JWT_SECRET", jwtSecret)
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
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// Connect to Database
	db, err := postgres.New(databaseDSN)
	if err != nil {
		return fmt.Errorf("connecting to postgres: %w", err)
	}

	tokenService := auth.JWTService{
		Secret:        []byte(jwtSecret),
		SigningMethod: jwt.SigningMethodHS256,
		Leeway:        5 * time.Second,
	}

	hashService := auth.Argon2IDService{
		SaltLen: 32,
		Time:    1,
		Memory:  64 * 1024, // 64 MB
		Threads: 4,
		KeyLen:  32,
	}

	// Unprotected routes
	r.Group(func(r chi.Router) {
		r.Use(httprate.LimitByIP(100, 1*time.Minute))

		// Login
		r.Post("/login", auth.LoginHandler(db, tokenService, hashService))

		docsPath := "/docs"
		openAPIPath := docsPath + "/swagger.yaml"

		// Docs Redirect
		r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, docsPath, http.StatusFound)
		})

		// Swagger Docs
		r.Handle(docsPath, api.SwaggerUIHandler(api.SwaggerUIOpts{
			Path:    docsPath,
			SpecURL: openAPIPath,
			Title:   "FoREST API Docs",
		}))
		r.Get(openAPIPath, api.SwaggerSpecsHandlerFunc())
	})

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(auth.Middleware(tokenService))

		// Mount Features
		r.Mount("/data", data.NewService(db))
		r.Mount("/mesh-nodes", mesh_node.NewService(db))
		r.Route("/accounts", func(r chi.Router) {
			r.Mount("/users", user_account.NewService(db, hashService))
			r.Mount("/services", service_account.NewService(db, tokenService))
		})
		r.Mount("/roles", role.NewService(db))
		r.Mount("/mesh-node-updates", mesh_node_update.NewService(db))
		r.Mount("/areas", area.NewService())
		r.Delete("/logout", auth.LogoutHandler())
	})

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
