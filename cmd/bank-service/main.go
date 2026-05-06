package main

import (
	"net/http"

	"github.com/evbgsl/bank-service-evbgsl/internal/config"
	"github.com/evbgsl/bank-service-evbgsl/internal/db"
	"github.com/evbgsl/bank-service-evbgsl/internal/handlers"
	"github.com/evbgsl/bank-service-evbgsl/internal/middleware"
	"github.com/evbgsl/bank-service-evbgsl/internal/repositories"
	"github.com/evbgsl/bank-service-evbgsl/internal/services"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func main() {
	log := logrus.New()
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	log.SetLevel(logrus.InfoLevel)

	cfg := config.Load()

	database, err := db.NewPostgresConnection(cfg.DB)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer database.Close()

	log.Info("database connection established")

	userRepository := repositories.NewUserRepository(database)
	authService := services.NewAuthService(userRepository, cfg.JWTSecret)
	authHandler := handlers.NewAuthHandler(authService)

	router := mux.NewRouter()

	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	}).Methods(http.MethodGet)

	router.HandleFunc("/health/db", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if err := database.Ping(); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			_, _ = w.Write([]byte(`{"status":"error","database":"unavailable"}`))
			return
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok","database":"available"}`))
	}).Methods(http.MethodGet)

	router.HandleFunc("/register", authHandler.Register).Methods(http.MethodPost)
	router.HandleFunc("/login", authHandler.Login).Methods(http.MethodPost)

	authRouter := router.PathPrefix("/").Subrouter()
	authRouter.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	authRouter.HandleFunc("/me", authHandler.Me).Methods(http.MethodGet)

	serverAddr := ":" + cfg.AppPort

	log.Infof("bank service started on %s", serverAddr)

	if err := http.ListenAndServe(serverAddr, router); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
