package main

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/mdma-backend/mdma-backend/internal/api/service_account"
	"github.com/mdma-backend/mdma-backend/internal/api/user_account"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
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
	docsPath := "/docs"
	openAPIPath := docsPath + "/swagger.yaml"
	loginPath := "/login"
	tokenService := auth.JWTService{
		Secret:        []byte("change_me"),
		SigningMethod: jwt.SigningMethodHS256,
		Leeway:        5 * time.Second,
	}

	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(auth.Middleware(tokenService, "/", docsPath, openAPIPath, loginPath))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"*"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: false,
	}))

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, docsPath, http.StatusFound)
	})

	r.Handle(docsPath, api.SwaggerUIHandler(api.SwaggerUIOpts{
		Path:    docsPath,
		SpecURL: openAPIPath,
		Title:   "FoREST API Docs",
	}))
	r.Get(openAPIPath, api.SwaggerSpecsHandlerFunc())

	db, err := postgres.New(databaseDSN)
	if err != nil {
		return fmt.Errorf("connecting to postgres: %w", err)
	}

	r.Mount("/data", data.NewService(db))
	r.Mount("/mesh-node", mesh_node.NewService())
	r.Route("/accounts", func(r chi.Router) {
		r.Mount("/users", user_account.NewService(db))
		r.Mount("/services", service_account.NewService(db))
	})
	r.Mount("/roles", role.NewService(db))

	hashService := auth.Argon2IDService{
		SaltLen: 32,
		Time:    1,
		Memory:  64 * 1024, // 64 MB
		Threads: 4,
		KeyLen:  32,
	}

	hash, salt, err := hashService.Hash("password123")
	fmt.Printf("hash=%s salt=%s err=%v", base64.StdEncoding.EncodeToString(hash), base64.StdEncoding.EncodeToString(salt), err)

	r.Post(loginPath, auth.LoginHandler(db, tokenService, hashService))
	r.Delete("/logout", auth.LogoutHandler())

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
