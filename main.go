package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux" // For HTTP routing.

	"github.com/software-architecture-proj/nova-backend-api-gateway/config"

	"github.com/software-architecture-proj/nova-backend-api-gateway/internal/clients"
	"github.com/software-architecture-proj/nova-backend-api-gateway/internal/handlers"
	"github.com/software-architecture-proj/nova-backend-api-gateway/internal/middleware"
)

func main() {
	cfg := config.LoadConfig()

	// Create a gRPC client for each Service.
	userProductClient, err := clients.NewUserProductServiceClient(cfg.UserProductServiceGRPCHost)
	if err != nil {
		log.Fatalf("Failed to create UserProductServiceClient: %v", err) //  Critical
	}
	defer userProductClient.CloseConnection() // Ensure connection is closed when main exits.

	AuthClient, err := clients.NewAuthServiceClient(cfg.AuthServiceGRPCHost)
	if err != nil {
		log.Fatalf("Failed to create AuthServiceClient: %v", err) //  Critical
	}
	defer AuthClient.CloseConnection() // Ensure connection is closed when main exits.

	TransactionClient, err := clients.NewTransactionServiceClient(cfg.TransactionServiceGRPCHost)
	if err != nil {
		log.Fatalf("Failed to create TransactionServiceClient: %v", err) //  Critical
	}
	defer TransactionClient.CloseConnection() // Ensure connection is closed when main exits.

	// Set up HTTP router
	router := mux.NewRouter()
	apiRouter := router.PathPrefix("/api").Subrouter()

	// Initialize HTTP handlers
	userProductHandler := handlers.NewUserProductHandler(userProductClient, TransactionClient, AuthClient)
	AuthHandler := handlers.NewAuthHandler(AuthClient)
	TransactionHandler := handlers.NewTransactionHandler(TransactionClient)
	Middle := middleware.NewMiddleware()
	// ... initialize other handlers (e.g., productHandler, accountHandler)
	apiRouter.HandleFunc("/country-codes", userProductHandler.GetCountryCodes).Methods(http.MethodGet)
	apiRouter.HandleFunc("/users", userProductHandler.CreateUser).Methods(http.MethodPost)
	apiRouter.HandleFunc("/users/{user_id}", userProductHandler.GetUser).Methods(http.MethodGet)
	apiRouter.HandleFunc("/login", AuthHandler.PostLogin).Methods(http.MethodPost)
	apiRouter.HandleFunc("/balance", TransactionHandler.GetBalance).Methods(http.MethodGet)
	apiRouter.HandleFunc("/movements", TransactionHandler.GetMovements).Methods(http.MethodGet)

	apiRouter.Use(Middle.AuthToken) // Middleware for logging requests
	// User and Products routes
	apiRouter.HandleFunc("/users/{user_id}", userProductHandler.UpdateUser).Methods(http.MethodPut)
	apiRouter.HandleFunc("/users/{user_id}", userProductHandler.DeleteUser).Methods(http.MethodDelete)

	// Favorites routes
	apiRouter.HandleFunc("/users/{user_id}/favorites", userProductHandler.GetFavoritesByUserId).Methods(http.MethodGet)
	apiRouter.HandleFunc("/users/{user_id}/favorites", userProductHandler.CreateFavorite).Methods(http.MethodPost)
	apiRouter.HandleFunc("/users/{user_id}/favorites/{favorite_id}", userProductHandler.UpdateFavorite).Methods(http.MethodPut)
	apiRouter.HandleFunc("/users/{user_id}/favorites/{favorite_id}", userProductHandler.DeleteFavorite).Methods(http.MethodDelete)

	// Pockets routes
	apiRouter.HandleFunc("/users/{user_id}/pockets", userProductHandler.GetPocketsByUserId).Methods(http.MethodGet)
	apiRouter.HandleFunc("/users/{user_id}/pockets", userProductHandler.CreatePocket).Methods(http.MethodPost)
	apiRouter.HandleFunc("/users/{user_id}/pockets/{pocket_id}", userProductHandler.UpdatePocket).Methods(http.MethodPut)
	apiRouter.HandleFunc("/users/{user_id}/pockets/{pocket_id}", userProductHandler.DeletePocket).Methods(http.MethodDelete)

	// Verification routes
	apiRouter.HandleFunc("/users/{user_id}/verifications", userProductHandler.GetVerificationsByUserId).Methods(http.MethodGet)
	apiRouter.HandleFunc("/users/{user_id}/verifications", userProductHandler.UpdateVerificationByUserId).Methods(http.MethodPut)

	// Auth routes

	// Transaction routes
	apiRouter.HandleFunc("/accounts", TransactionHandler.PostAccount).Methods(http.MethodPost) //deprecated
	apiRouter.HandleFunc("/transfers", TransactionHandler.PostTransfer).Methods(http.MethodPost)

	// Create HTTP server
	server := &http.Server{
		Addr:         ":" + cfg.APIGatewayPort,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Start HTTP server in a goroutine.
	go func() {
		log.Printf("API Gateway listening on port %s", cfg.APIGatewayPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("API Gateway server failed: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down API Gateway...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("API Gateway forced to shutdown: %v", err)
	}
	log.Println("API Gateway exited gracefully.")
}
