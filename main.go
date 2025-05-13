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

    // @trigger: Will be implemented later. 
	//"github.com/software-architecture-proj/nova-backend-api-gateway/internal/clients"
	//"github.com/software-architecture-proj/nova-backend-api-gateway/internal/handlers"
)

func main() {
	cfg := config.LoadConfig()

	// Set up HTTP router
	router := mux.NewRouter()
    
    // @trigger: Also will be implemented later.
    // Initialize HTTP handlers
	//userHandler := handlers.NewUserHandler(userProductClient)
	// ... initialize other handlers (e.g., productHandler, accountHandler)

	// Define API Gateway endpoints
	//router.HandleFunc("/users", userHandler.CreateUser).Methods("POST")
	// router.HandleFunc("/users/{user_id}", userHandler.GetUser).Methods("GET")
	// router.HandleFunc("/products", productHandler.CreateProduct).Methods("POST")
	// ... add more routes for other services/methods

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
