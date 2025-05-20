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
)

func main() {
	cfg := config.LoadConfig()

    // Create a gRPC client for the UserProductService.
	userProductClient, err := clients.NewUserProductServiceClient(cfg.UserProductServiceGRPCHost)
	if err != nil {
		log.Fatalf("Failed to create UserProductServiceClient: %v", err) //  Critical
	}
	defer userProductClient.CloseConnection() // Ensure connection is closed when main exits.

	// Set up HTTP router
	router := mux.NewRouter()
    
    // Initialize HTTP handlers
	userProductHandler := handlers.NewUserProductHandler(userProductClient)
	// ... initialize other handlers (e.g., productHandler, accountHandler)

    // User and Products routes
	router.HandleFunc("", userProductHandler.CreateUser).Methods(http.MethodPost)
	router.HandleFunc("/{user_id}", userProductHandler.GetUser).Methods(http.MethodGet)
	router.HandleFunc("/{user_id}", userProductHandler.UpdateUser).Methods(http.MethodPut) // Corrected to PUT
	router.HandleFunc("/{user_id}", userProductHandler.DeleteUser).Methods(http.MethodDelete)
	router.HandleFunc("/{user_id}/favorites", userProductHandler.GetFavoritesByUserId).Methods(http.MethodGet)
	router.HandleFunc("/{user_id}/favorites", userProductHandler.CreateFavorite).Methods(http.MethodPost)
	router.HandleFunc("/{user_id}/favorites/{favorite_id}", userProductHandler.UpdateFavorite).Methods(http.MethodPut) // corrected
	router.HandleFunc("/{user_id}/favorites/{favorite_id}", userProductHandler.DeleteFavorite).Methods(http.MethodDelete)
	router.HandleFunc("/{user_id}/pockets", userProductHandler.GetPocketsByUserId).Methods(http.MethodGet)
	router.HandleFunc("/{user_id}/pockets", userProductHandler.CreatePocket).Methods(http.MethodPost)
	router.HandleFunc("/{user_id}/pockets/{pocket_id}", userProductHandler.UpdatePocketById).Methods(http.MethodPut) // Corrected
	router.HandleFunc("/{user_id}/pockets/{pocket_id}", userProductHandler.DeletePocketById).Methods(http.MethodDelete)

	// Create HTTP server
	server := &http.Server{
		Addr:         ":" + cfg_str.APIGatewayPort,
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
