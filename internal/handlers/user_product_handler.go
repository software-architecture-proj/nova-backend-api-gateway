package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/software-architecture-proj/nova-backend-api-gateway/internal/clients"
	"github.com/software-architecture-proj/nova-backend-api-gateway/internal/transformers"
	"github.com/software-architecture-proj/nova-backend-api-gateway/internal/common"

	// Import from common-protos
	pb "github.com/software-architecture-proj/nova-backend-common-protos/gen/go/user_product_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserHandler struct {
	UserProductClient *clients.UserProductServiceClient
}

func NewUserHandler(userProductClient *clients.UserProductServiceClient) *UserHandler {
	return &UserHandler{UserProductClient: userProductClient}
}

// CreateUser handles POST /users
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var reqBody struct {
		Email     string `json:"email"`
		Username  string `json:"username"`
		Phone     string `json:"phone"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Birthdate string `json:"birthdate"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	grpcReq := &pb.CreateUserRequest{
		Email:     reqBody.Email,
		Username:  reqBody.Username,
		Phone:     reqBody.Phone,
		FirstName: reqBody.FirstName,
		LastName:  reqBody.LastName,
		Birthdate: reqBody.Birthdate,
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	grpcResp, err := h.UserProductClient.Client.CreateUser(ctx, grpcReq)
	if err != nil {
		RespondGrpcError(w, err)
		return
	}

	// Assuming you have a transformer function to convert the gRPC response to a suitable HTTP response
	httpResp := transformers.ToUserJSON(grpcResp) //  Create this function in transformers.go
	RespondWithJSON(w, http.StatusCreated, httpResp)
}

// GetUser handles GET /users/{user_id}
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["user_id"]
	if userID == "" {
		RespondWithError(w, http.StatusBadRequest, "Missing user_id")
		return
	}

	grpcReq := &pb.GetUserByIdRequest{UserId: userID}
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	grpcResp, err := h.UserProductClient.Client.GetUserById(ctx, grpcReq)
	if err != nil {
		RespondGrpcError(w, err)
		return
	}

	httpResp := transformers.ToUserJSON(grpcResp) // Create this function in transformers.go
	RespondWithJSON(w, http.StatusOK, httpResp)
}

// UpdateUser handles PUT /users/{user_id}  
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["user_id"]
	if userID == "" {
		RespondWithError(w, http.StatusBadRequest, "Missing user_id")
		return
	}

	var reqBody struct {
		Email     string `json:"email"`
		Username  string `json:"username"`
		Phone     string `json:"phone"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Birthdate string `json:"birthdate"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	grpcReq := &pb.UpdateUserByIdRequest{ // Corrected struct name.
		Id:        userID, // Use the user_id from the path.
		Email:     reqBody.Email,
		Username:  reqBody.Username,
		Phone:     reqBody.Phone,
		FirstName: reqBody.FirstName,
		LastName:  reqBody.LastName,
		Birthdate: reqBody.Birthdate,
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	grpcResp, err := h.UserProductClient.Client.UpdateUserById(ctx, grpcReq) // Corrected method name
	if err != nil {
		RespondGrpcError(w, err)
		return
	}

	httpResp := transformers.ToUserJSON(grpcResp) // Create this function in transformers.go
	RespondWithJSON(w, http.StatusOK, httpResp)
}

// DeleteUser handles DELETE /users/{user_id}
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["user_id"]
	if userID == "" {
		RespondWithError(w, http.StatusBadRequest, "Missing user_id")
		return
	}

	grpcReq := &pb.DeleteUserByIdRequest{Id: userID} // Corrected struct.
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	_, err := h.UserProductClient.Client.DeleteUserById(ctx, grpcReq) // Corrected method.  Check the return.
	if err != nil {
		RespondGrpcError(w, err)
		return
	}

	RespondWithJSON(w, http.StatusOK, map[string]bool{"success": true}) //  Successful deletion
}

// GetFavoritesByUserId handles GET /users/{user_id}/favorites
func (h *UserHandler) GetFavoritesByUserId(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["user_id"]
	if userID == "" {
		RespondWithError(w, http.StatusBadRequest, "Missing user_id")
		return
	}

	grpcReq := &pb.GetFavoritesByUserIdRequest{UserId: userID}
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	grpcResp, err := h.UserProductClient.Client.GetFavoritesByUserId(ctx, grpcReq)
	if err != nil {
		RespondGrpcError(w, err)
		return
	}

	httpResp := transformers.ToFavoriteListJSON(grpcResp.GetFavorites()) //  Create this in transformers.go
	RespondWithJSON(w, http.StatusOK, httpResp)
}

// CreateFavorite handles POST /users/{user_id}/favorites
func (h *UserHandler) CreateFavorite(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    userID := vars["user_id"]
    if userID == "" {
        RespondWithError(w, http.StatusBadRequest, "Missing user_id")
        return
    }

    var reqBody struct {
        FavoriteUserId string `json:"favorite_user_id"`
        Alias          string `json:"alias"`
    }

    if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
        RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
        return
    }

    grpcReq := &pb.CreateFavoriteRequest{
        UserId:         userID,  // Use user_id from path
        FavoriteUserId: reqBody.FavoriteUserId,
        Alias:          reqBody.Alias,
    }

    ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
    defer cancel()

    grpcResp, err := h.UserProductClient.Client.CreateFavorite(ctx, grpcReq)
    if err != nil {
        RespondGrpcError(w, err)
        return
    }
    httpResp := transformers.ToFavoriteJSON(grpcResp)  // Create this in transformers
    RespondWithJSON(w, http.StatusCreated, httpResp)
}

// UpdateFavorite handles PUT /users/{user_id}/favorites/{favorite_id}
func (h *UserHandler) UpdateFavorite(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    userID := vars["user_id"]  // You might need this for context
    favoriteID := vars["favorite_id"]
    if userID == "" || favoriteID == "" {
        RespondWithError(w, http.StatusBadRequest, "Missing user_id or favorite_id")
        return
    }

    var reqBody struct {
        Alias string `json:"alias"`
    }
    if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
        RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
        return
    }

    grpcReq := &pb.UpdateFavoriteByIdRequest{ // Corrected name
        Id:    favoriteID,
        Alias: reqBody.Alias,
    }

    ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
    defer cancel()

    grpcResp, err := h.UserProductClient.Client.UpdateFavoriteById(ctx, grpcReq)  //Corrected name
    if err != nil {
        RespondGrpcError(w, err)
        return
    }

    httpResp := transformers.ToFavoriteJSON(grpcResp) // Create this
    RespondWithJSON(w, http.StatusOK, httpResp)
}

// DeleteFavorite handles DELETE /users/{user_id}/favorites/{favorite_id}
func (h *UserHandler) DeleteFavorite(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    userID := vars["user_id"]  //  May need this
    favoriteID := vars["favorite_id"]
    if userID == "" || favoriteID == ""{
        RespondWithError(w, http.StatusBadRequest, "Missing user_id or favorite_id")
        return
    }
    grpcReq := &pb.DeleteFavoriteByIdRequest{Id: favoriteID}  // Corrected
    ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
    defer cancel()

    _, err := h.UserProductClient.Client.DeleteFavoriteById(ctx, grpcReq) // Corrected
    if err != nil{
        RespondGrpcError(w, err)
        return
    }
    RespondWithJSON(w, http.StatusOK, map[string]bool{"success": true})
}
