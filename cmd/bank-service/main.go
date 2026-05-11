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
	accountRepository := repositories.NewAccountRepository(database)
	transactionRepository := repositories.NewTransactionRepository(database)
	cardRepository := repositories.NewCardRepository(database)

	authService := services.NewAuthService(userRepository, cfg.JWTSecret)
	accountService := services.NewAccountService(accountRepository)
	transactionService := services.NewTransactionService(transactionRepository)
	cardService := services.NewCardService(cardRepository, accountRepository)

	authHandler := handlers.NewAuthHandler(authService)
	accountHandler := handlers.NewAccountHandler(accountService)
	transactionHandler := handlers.NewTransactionHandler(transactionService)
	cardHandler := handlers.NewCardHandler(cardService)

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

	authRouter.HandleFunc("/accounts", accountHandler.CreateAccount).Methods(http.MethodPost)
	authRouter.HandleFunc("/accounts", accountHandler.GetUserAccounts).Methods(http.MethodGet)
	authRouter.HandleFunc("/accounts/{accountId}/deposit", accountHandler.Deposit).Methods(http.MethodPost)

	authRouter.HandleFunc("/transfer", transactionHandler.Transfer).Methods(http.MethodPost)
	authRouter.HandleFunc("/transactions", transactionHandler.GetUserTransactions).Methods(http.MethodGet)

	authRouter.HandleFunc("/cards", cardHandler.CreateCard).Methods(http.MethodPost)
	authRouter.HandleFunc("/cards", cardHandler.GetUserCards).Methods(http.MethodGet)

	serverAddr := ":" + cfg.AppPort

	log.Infof("bank service started on %s", serverAddr)

	if err := http.ListenAndServe(serverAddr, router); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
