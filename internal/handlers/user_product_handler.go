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

type UserProductHandler struct {
	UserProductClient *clients.UserProductServiceClient
}

func NewUserProductHandler(userProductClient *clients.UserProductServiceClient) *UserProductHandler {
	return &UserProductHandler{UserProductClient: userProductClient}
}

// GetCountryCodes handles GET /country-codes
func (h *UserProductHandler) GetCountryCodes(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	grpcResp, err := h.UserProductClient.Client.GetCountryCodes(ctx, &pb.GetCountryCodesRequest{})
	if err != nil {
		common.RespondGrpcError(w, err)
		return
	}

	httpResp := transformers.GetCountryCodesRespJSON(grpcResp)
	common.RespondWithJSON(w, http.StatusOK, httpResp)
}

// CreateUser handles POST /users
func (h *UserProductHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var reqBody struct {
		Email     string `json:"email"`
		Username  string `json:"username"`
        CodeId    string `json:"code_id"`
		Phone     string `json:"phone"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Birthdate string `json:"birthdate"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		common.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	grpcReq := &pb.CreateUserRequest{
		Email:     reqBody.Email,
		Username:  reqBody.Username,
        CodeId:    reqBody.CodeId,
		Phone:     reqBody.Phone,
		FirstName: reqBody.FirstName,
		LastName:  reqBody.LastName,
		Birthdate: reqBody.Birthdate,
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	grpcResp, err := h.UserProductClient.Client.CreateUser(ctx, grpcReq)
	if err != nil {
		common.RespondGrpcError(w, err)
		return
	}

	// Assuming you have a transformer function to convert the gRPC response to a suitable HTTP response
	httpResp := transformers.CreateUserRespJSON(grpcResp) //  Create this function in transformers.go
	common.RespondWithJSON(w, http.StatusCreated, httpResp)
}

// GetUser handles GET /users/{user_id}
func (h *UserProductHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["user_id"]
	if userID == "" {
		common.RespondWithError(w, http.StatusBadRequest, "Missing user_id")
		return
	}

	grpcReq := &pb.GetUserByIdRequest{UserId: userID}
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	grpcResp, err := h.UserProductClient.Client.GetUserById(ctx, grpcReq)
	if err != nil {
		common.RespondGrpcError(w, err)
		return
	}

	httpResp := transformers.GetUserRespJSON(grpcResp) // Create this function in transformers.go
	common.RespondWithJSON(w, http.StatusOK, httpResp)
}

// UpdateUser handles PUT /users/{user_id}  
func (h *UserProductHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["user_id"]
	if userID == "" {
		common.RespondWithError(w, http.StatusBadRequest, "Missing user_id")
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
		common.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
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
		common.RespondGrpcError(w, err)
		return
	}

	httpResp := transformers.UpdateUserRespJSON(grpcResp) // Create this function in transformers.go
	common.RespondWithJSON(w, http.StatusOK, httpResp)
}

// DeleteUser handles DELETE /users/{user_id}
func (h *UserProductHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["user_id"]
	if userID == "" {
		common.RespondWithError(w, http.StatusBadRequest, "Missing user_id")
		return
	}

	grpcReq := &pb.DeleteUserByIdRequest{Id: userID} // Corrected struct.
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	grcpResp, err := h.UserProductClient.Client.DeleteUserById(ctx, grpcReq) // Corrected method.  Check the return.
	if err != nil {
		common.RespondGrpcError(w, err)
		return
	}
    
	httpResp := transformers.DeleteUserRespJSON(grpcResp) // Create this function in transformers.go
	common.RespondWithJSON(w, http.StatusOK, httpResp) 
}

// GetFavoritesByUserId handles GET /users/{user_id}/favorites
func (h *UserProductHandler) GetFavoritesByUserId(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["user_id"]
	if userID == "" {
		common.RespondWithError(w, http.StatusBadRequest, "Missing user_id")
		return
	}

	grpcReq := &pb.GetFavoritesByUserIdRequest{UserId: userID}
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	grpcResp, err := h.UserProductClient.Client.GetFavoritesByUserId(ctx, grpcReq)
	if err != nil {
		common.RespondGrpcError(w, err)
		return
	}

	httpResp := transformers.GetFavoritesRespJSON(grpcResp) //  Create this in transformers.go
	common.RespondWithJSON(w, http.StatusOK, httpResp)
}

// CreateFavorite handles POST /users/{user_id}/favorites
func (h *UserProductHandler) CreateFavorite(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    userID := vars["user_id"]
    if userID == "" {
        common.RespondWithError(w, http.StatusBadRequest, "Missing user_id")
        return
    }

    var reqBody struct {
        FavoriteUserId string `json:"favorite_user_id"`
        Alias          string `json:"alias"`
    }

    if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
        common.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
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
        common.RespondGrpcError(w, err)
        return
    }
    httpResp := transformers.CreateFavoriteRespJSON(grpcResp)  // Create this in transformers
    common.RespondWithJSON(w, http.StatusCreated, httpResp)
}

// UpdateFavorite handles PUT /users/{user_id}/favorites/{favorite_id}
func (h *UserProductHandler) UpdateFavorite(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    userID := vars["user_id"]   
    favoriteID := vars["favorite_id"]
    if userID == "" || favoriteID == "" {
        common.RespondWithError(w, http.StatusBadRequest, "Missing user_id or favorite_id")
        return
    }

    var reqBody struct {
        Alias string `json:"alias"`
    }
    if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
        common.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
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
        common.RespondGrpcError(w, err)
        return
    }

    httpResp := transformers.UpdateFavoriteRespJSON(grpcResp) // Create this
    common.RespondWithJSON(w, http.StatusOK, httpResp)
}

// DeleteFavorite handles DELETE /users/{user_id}/favorites/{favorite_id}
func (h *UserProductHandler) DeleteFavorite(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    userID := vars["user_id"]  //  May need this
    favoriteID := vars["favorite_id"]
    if userID == "" || favoriteID == ""{
        common.RespondWithError(w, http.StatusBadRequest, "Missing user_id or favorite_id")
        return
    }
    grpcReq := &pb.DeleteFavoriteByIdRequest{Id: favoriteID}  // Corrected
    ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
    defer cancel()

    grpcResp, err := h.UserProductClient.Client.DeleteFavoriteById(ctx, grpcReq) // Corrected
    if err != nil{
        common.RespondGrpcError(w, err)
        return
    }
	httpResp := transformers.DeleteFavoriteRespJSON(grpcResp) // Create this function in transformers.go
	common.RespondWithJSON(w, http.StatusOK, httpResp) 
}

// GetPocketsByUserId handles GET /users/{user_id}/pockets
func (h *UserProductHandler) GetPocketsByUserId(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["user_id"]
	if userID == "" {
		common.RespondWithError(w, http.StatusBadRequest, "Missing user_id")
		return
	}

	grpcReq := &pb.GetPocketsByUserIdRequest{UserId: userID}
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	grpcResp, err := h.UserProductClient.Client.GetPocketsByUserId(ctx, grpcReq)
	if err != nil {
		common.RespondGrpcError(w, err)
		return
	}

	httpResp := transformers.GetPocketsRespJSON(grpcResp)
	common.RespondWithJSON(w, http.StatusOK, httpResp)
}

// CreatePocket handles POST /users/{user_id}/pockets
func (h *UserProductHandler) CreatePocket(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    userID := vars["user_id"]
    if userID == "" {
        common.RespondWithError(w, http.StatusBadRequest, "Missing user_id")
        return
    }

    var reqBody struct {
        Name           string `json:"name"`
        Category       string `json:"category"`
        MaxAmount      int32  `json:"max_amount"`
    }

    if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
        common.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
        return
    }

    grpcReq := &pb.CreatePocketRequest{
        UserId:         userID,  // Use user_id from path
        Name:           reqBody.Name,
        Category:       reqBody.Alias,
        MaxAmount:      reqBody.MaxAmount,
    }

    ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
    defer cancel()

    grpcResp, err := h.UserProductClient.Client.CreatePocket(ctx, grpcReq)
    if err != nil {
        common.RespondGrpcError(w, err)
        return
    }
    httpResp := transformers.CreatePocketRespJSON(grpcResp)  // Create this in transformers
    common.RespondWithJSON(w, http.StatusCreated, httpResp)
}

// UpdatePocket handles PUT /users/{user_id}/pockets/{pocket_id}
func (h *UserProductHandler) UpdatePocket(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    userID := vars["user_id"]   
    pocketID := vars["pocket_id"]
    if userID == "" || pocketID == "" {
        common.RespondWithError(w, http.StatusBadRequest, "Missing user_id or pocket_id")
        return
    }

    var reqBody struct {
        Name           string `json:"name"`
        Category       string `json:"category"`
        MaxAmount      int32  `json:"max_amount"`
    }
    if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
        common.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
        return
    }

    grpcReq := &pb.UpdatePocketByIdRequest{
        Id:             pocketID,  
        Name:           reqBody.Name,
        Category:       reqBody.Alias,
        MaxAmount:      reqBody.MaxAmount,
    }

    ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
    defer cancel()

    grpcResp, err := h.UserProductClient.Client.UpdatePocketById(ctx, grpcReq)
    if err != nil {
        common.RespondGrpcError(w, err)
        return
    }

    httpResp := transformers.UpdatePocketRespJSON(grpcResp) // Create this
    common.RespondWithJSON(w, http.StatusOK, httpResp)
}

// DeletePocket handles DELETE /users/{user_id}/pockets/{pocket_id}
func (h *UserProductHandler) DeletePocket(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    userID := vars["user_id"]  //  May need this
    pocketID := vars["pocket_id"]
    if userID == "" || pocketID == ""{
        common.RespondWithError(w, http.StatusBadRequest, "Missing user_id or pocket_id")
        return
    }
    grpcReq := &pb.DeletePocketByIdRequest{Id: pocketID}  
    ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
    defer cancel()

    grpcResp, err := h.UserProductClient.Client.DeletePocketById(ctx, grpcReq) 
    if err != nil{
        common.RespondGrpcError(w, err)
        return
    }

	httpResp := transformers.DeletePocketRespJSON(grpcResp) // Create this function in transformers.go
	common.RespondWithJSON(w, http.StatusOK, httpResp) 
}
