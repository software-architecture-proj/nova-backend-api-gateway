package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/software-architecture-proj/nova-backend-api-gateway/internal/clients"
	"github.com/software-architecture-proj/nova-backend-api-gateway/internal/common"
	"github.com/software-architecture-proj/nova-backend-api-gateway/internal/transformers"

	// Import from common-protos
	ab "github.com/software-architecture-proj/nova-backend-common-protos/gen/go/auth_service"
	tb "github.com/software-architecture-proj/nova-backend-common-protos/gen/go/transaction_service"
	pb "github.com/software-architecture-proj/nova-backend-common-protos/gen/go/user_product_service"
)

type UserProductHandler struct {
	UserProductClient *clients.UserProductServiceClient
	TransactionClient *clients.TransactionServiceClient
	AuthClient        *clients.AuthServiceClient
}

func NewUserProductHandler(userClient *clients.UserProductServiceClient, transactionClient *clients.TransactionServiceClient, authClient *clients.AuthServiceClient) *UserProductHandler {
	return &UserProductHandler{
		UserProductClient: userClient,
		TransactionClient: transactionClient,
		AuthClient:        authClient,
	}
}

// GetCountryCodes handles GET /country-codes
func (h *UserProductHandler) GetCountryCodes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		common.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
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
	if r.Method != http.MethodPost {
		common.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var reqBody struct {
		Email     string `json:"email"`
		Username  string `json:"username"`
		Password  string `json:"password"`
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

	grpcReqAuth := &ab.CreateUserRequest{
		Username: reqBody.Username,
		Phone:    reqBody.Phone,
		Password: reqBody.Password,
		Email:    reqBody.Email,
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	authResp, authErr := h.AuthClient.Client.CreateUser(ctx, grpcReqAuth)
	if authErr != nil {
		defer cancel()
		log.Println("Error creating user in auth service:", authErr)
		common.RespondGrpcError(w, authErr)
		return
	}
	userID := authResp.Data
	log.Println("User created in auth service with ID:", userID)
	grpcReqTB := &tb.CreateAccountRequest{
		UserId:   userID,
		Username: reqBody.Username,
		Bank:     false,
	}

	grpcReqUS := &pb.CreateUserRequest{
		UserId:    userID,
		Email:     reqBody.Email,
		Username:  reqBody.Username,
		CodeId:    reqBody.CodeId,
		Phone:     reqBody.Phone,
		FirstName: reqBody.FirstName,
		LastName:  reqBody.LastName,
		Birthdate: reqBody.Birthdate,
	}

	var wg sync.WaitGroup
	wg.Add(2)

	var userResp *pb.CreateUserResponse
	var userErr error
	var tbErr error

	go func() {
		defer wg.Done()
		userResp, userErr = h.UserProductClient.Client.CreateUser(ctx, grpcReqUS)
	}()

	go func() {
		defer wg.Done()
		_, tbErr = h.TransactionClient.Client.Account(ctx, grpcReqTB)
		if tbErr != nil {
			log.Println("Error creating account: ", tbErr)
		}
	}()

	wg.Wait()

	if userErr != nil {
		defer cancel()
		log.Println("Error creating user:", userErr)
		common.RespondGrpcError(w, userErr)
		return
	}
	if tbErr != nil {
		defer cancel()
		log.Println("Error creating account v2:", tbErr)
		common.RespondGrpcError(w, tbErr)
		return
	}

	httpResp := transformers.CreateUserRespJSON(userResp)
	common.RespondWithJSON(w, http.StatusCreated, httpResp)
	defer cancel()
}

// GetUser handles GET /users/{user_id}
func (h *UserProductHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		common.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	vars := mux.Vars(r)
	userID := vars["user_id"]
	if userID == "" {
		common.RespondWithError(w, http.StatusBadRequest, "Missing user_id")
		return
	}

	grpcReq := &pb.GetUserByIdRequest{UserId: userID}
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)

	grpcResp, err := h.UserProductClient.Client.GetUserById(ctx, grpcReq)
	if err != nil {
		defer cancel()
		common.RespondGrpcError(w, err)
		return
	}

	httpResp := transformers.GetUserRespJSON(grpcResp) // Create this function in transformers.go
	common.RespondWithJSON(w, http.StatusOK, httpResp)
	defer cancel()
}

// GetUser handles GET /users/name/{username}
func (h *UserProductHandler) GetUsername(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		common.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	vars := mux.Vars(r)
	uname := vars["username"]
	if uname == "" {
		common.RespondWithError(w, http.StatusBadRequest, "Missing username")
		return
	}

	grpcReq := &pb.GetUserByUsernameRequest{Username: uname}
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)

	grpcResp, err := h.UserProductClient.Client.GetUserByUsername(ctx, grpcReq)
	if err != nil {
		common.RespondGrpcError(w, err)
		defer cancel()
		return
	}

	httpResp := transformers.GetUsernameRespJSON(grpcResp)
	common.RespondWithJSON(w, http.StatusOK, httpResp)
	defer cancel()
}

// UpdateUser handles PUT /users/{user_id}
func (h *UserProductHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		common.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
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

	grpcResp, err := h.UserProductClient.Client.UpdateUserById(ctx, grpcReq) // Corrected method name
	if err != nil {
		defer cancel()
		common.RespondGrpcError(w, err)
		return
	}

	httpResp := transformers.UpdateUserRespJSON(grpcResp) // Create this function in transformers.go
	common.RespondWithJSON(w, http.StatusOK, httpResp)
	defer cancel()
}

// DeleteUser handles DELETE /users/{user_id}
func (h *UserProductHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		common.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	vars := mux.Vars(r)
	userID := vars["user_id"]
	if userID == "" {
		common.RespondWithError(w, http.StatusBadRequest, "Missing user_id")
		return
	}

	grpcReq := &pb.DeleteUserByIdRequest{Id: userID} // Corrected struct.
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)

	grpcResp, err := h.UserProductClient.Client.DeleteUserById(ctx, grpcReq) // Corrected method.  Check the return.
	if err != nil {
		defer cancel()
		common.RespondGrpcError(w, err)
		return
	}

	httpResp := transformers.DeleteUserRespJSON(grpcResp) // Create this function in transformers.go
	common.RespondWithJSON(w, http.StatusOK, httpResp)
	defer cancel()
}

// GetFavoritesByUserId handles GET /users/{user_id}/favorites
func (h *UserProductHandler) GetFavoritesByUserId(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		common.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	vars := mux.Vars(r)
	userID := vars["user_id"]
	if userID == "" {
		common.RespondWithError(w, http.StatusBadRequest, "Missing user_id")
		return
	}

	grpcReq := &pb.GetFavoritesByUserIdRequest{UserId: userID}
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)

	grpcResp, err := h.UserProductClient.Client.GetFavoritesByUserId(ctx, grpcReq)
	if err != nil {
		defer cancel()
		common.RespondGrpcError(w, err)
		return
	}

	httpResp := transformers.GetFavoritesRespJSON(grpcResp) //  Create this in transformers.go
	common.RespondWithJSON(w, http.StatusOK, httpResp)
	defer cancel()
}

// CreateFavorite handles POST /users/{user_id}/favorites
func (h *UserProductHandler) CreateFavorite(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		common.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
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
		UserId:         userID, // Use user_id from path
		FavoriteUserId: reqBody.FavoriteUserId,
		Alias:          reqBody.Alias,
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)

	grpcResp, err := h.UserProductClient.Client.CreateFavorite(ctx, grpcReq)
	if err != nil {
		defer cancel()
		common.RespondGrpcError(w, err)
		return
	}

	httpResp := transformers.CreateFavoriteRespJSON(grpcResp) // Create this function in transformers.go
	common.RespondWithJSON(w, http.StatusCreated, httpResp)
	defer cancel()
}

// UpdateFavorite handles PUT /users/{user_id}/favorites/{favorite_id}
func (h *UserProductHandler) UpdateFavorite(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		common.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
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

	grpcResp, err := h.UserProductClient.Client.UpdateFavoriteById(ctx, grpcReq) //Corrected name
	if err != nil {
		defer cancel()
		common.RespondGrpcError(w, err)
		return
	}

	httpResp := transformers.UpdateFavoriteRespJSON(grpcResp) // Create this
	common.RespondWithJSON(w, http.StatusOK, httpResp)
	defer cancel()
}

// DeleteFavorite handles DELETE /users/{user_id}/favorites/{favorite_id}
func (h *UserProductHandler) DeleteFavorite(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		common.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	vars := mux.Vars(r)
	userID := vars["user_id"]
	favoriteID := vars["favorite_id"]
	if userID == "" || favoriteID == "" {
		common.RespondWithError(w, http.StatusBadRequest, "Missing user_id or favorite_id")
		return
	}

	grpcReq := &pb.DeleteFavoriteByIdRequest{Id: favoriteID}
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)

	grpcResp, err := h.UserProductClient.Client.DeleteFavoriteById(ctx, grpcReq)
	if err != nil {
		defer cancel()
		common.RespondGrpcError(w, err)
		return
	}

	httpResp := transformers.DeleteFavoriteRespJSON(grpcResp)
	common.RespondWithJSON(w, http.StatusOK, httpResp)
	defer cancel()
}

// GetPocketsByUserId handles GET /users/{user_id}/pockets
func (h *UserProductHandler) GetPocketsByUserId(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		common.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	vars := mux.Vars(r)
	userID := vars["user_id"]
	if userID == "" {
		common.RespondWithError(w, http.StatusBadRequest, "Missing user_id")
		return
	}

	grpcReq := &pb.GetPocketsByUserIdRequest{UserId: userID}
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	grpcResp, err := h.UserProductClient.Client.GetPocketsByUserId(ctx, grpcReq)
	if err != nil {
		defer cancel()

		common.RespondGrpcError(w, err)
		return
	}

	httpResp := transformers.GetPocketsRespJSON(grpcResp)
	common.RespondWithJSON(w, http.StatusOK, httpResp)
	defer cancel()

}

// CreatePocket handles POST /users/{user_id}/pockets
func (h *UserProductHandler) CreatePocket(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		common.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	vars := mux.Vars(r)
	userID := vars["user_id"]
	if userID == "" {
		common.RespondWithError(w, http.StatusBadRequest, "Missing user_id")
		return
	}

	var reqBody struct {
		Username  string `json:"username"`
		Name      string `json:"name"`
		Category  string `json:"category"`
		MaxAmount int32  `json:"max_amount"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		common.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	grpcReqUS := &pb.CreatePocketRequest{
		UserId:    userID,
		Name:      reqBody.Name,
		Category:  reqBody.Category,
		MaxAmount: reqBody.MaxAmount,
	}

	grpcReqTB := &tb.CreateAccountRequest{
		UserId:   uuid.New().String(), // Generate a new UUID for the account
		Username: reqBody.Username,
		Bank:     false,
	}
	log.Println("User: ", grpcReqUS.UserId, grpcReqUS.Category)
	log.Println("Trans: ", grpcReqTB.UserId, grpcReqTB.Username)

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)

	var wg sync.WaitGroup
	wg.Add(2)

	var pocketResp *pb.CreatePocketResponse
	var pocketErr error
	var tbErr error
	go func() {
		defer wg.Done()
		pocketResp, pocketErr = h.UserProductClient.Client.CreatePocket(ctx, grpcReqUS)
	}()
	go func() {
		defer wg.Done()
		_, tbErr = h.TransactionClient.Client.Account(ctx, grpcReqTB)
		if tbErr != nil {
			log.Println("Error creating pocket account: ", tbErr)
			defer cancel()
			return
		}
	}()
	wg.Wait()

	if pocketErr != nil {
		log.Println("Error creating pocket: ", pocketErr)
		common.RespondGrpcError(w, pocketErr)
		defer cancel()
		return
	}
	if tbErr != nil {
		log.Println("Error creating pocket account v2:", tbErr)
		common.RespondGrpcError(w, tbErr)
		defer cancel()
		return
	}

	httpResp := transformers.CreatePocketRespJSON(pocketResp) // Create this in transformers
	common.RespondWithJSON(w, http.StatusCreated, httpResp)
	defer cancel()

}

// UpdatePocket handles PUT /users/{user_id}/pockets/{pocket_id}
func (h *UserProductHandler) UpdatePocket(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		common.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	vars := mux.Vars(r)
	userID := vars["user_id"]
	pocketID := vars["pocket_id"]
	if userID == "" || pocketID == "" {
		common.RespondWithError(w, http.StatusBadRequest, "Missing user_id or pocket_id")
		return
	}

	var reqBody struct {
		Name      string `json:"name"`
		Category  string `json:"category"`
		MaxAmount int32  `json:"max_amount"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		common.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	grpcReq := &pb.UpdatePocketByIdRequest{
		Id:        pocketID,
		Name:      reqBody.Name,
		Category:  reqBody.Category,
		MaxAmount: reqBody.MaxAmount,
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)

	grpcResp, err := h.UserProductClient.Client.UpdatePocketById(ctx, grpcReq)
	if err != nil {
		defer cancel()

		common.RespondGrpcError(w, err)
		return
	}

	httpResp := transformers.UpdatePocketRespJSON(grpcResp)
	common.RespondWithJSON(w, http.StatusOK, httpResp)
	defer cancel()

}

// DeletePocket handles DELETE /users/{user_id}/pockets/{pocket_id}
func (h *UserProductHandler) DeletePocket(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		common.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	vars := mux.Vars(r)
	userID := vars["user_id"]
	pocketID := vars["pocket_id"]
	if userID == "" || pocketID == "" {
		common.RespondWithError(w, http.StatusBadRequest, "Missing user_id or pocket_id")
		return
	}

	grpcReq := &pb.DeletePocketByIdRequest{Id: pocketID}
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)

	grpcResp, err := h.UserProductClient.Client.DeletePocketById(ctx, grpcReq)
	if err != nil {
		defer cancel()

		common.RespondGrpcError(w, err)
		return
	}

	httpResp := transformers.DeletePocketRespJSON(grpcResp)
	common.RespondWithJSON(w, http.StatusOK, httpResp)
	defer cancel()

}

// GetVerificationsByUserId handles GET /users/{user_id}/verifications
func (h *UserProductHandler) GetVerificationsByUserId(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		common.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	vars := mux.Vars(r)
	userID := vars["user_id"]
	if userID == "" {
		common.RespondWithError(w, http.StatusBadRequest, "Missing user_id")
		return
	}

	grpcReq := &pb.GetVerificationsByUserIdRequest{UserId: userID}
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	grpcResp, err := h.UserProductClient.Client.GetVerificationsByUserId(ctx, grpcReq)
	if err != nil {
		defer cancel()

		common.RespondGrpcError(w, err)
		return
	}

	httpResp := transformers.GetVerificationsRespJSON(grpcResp)
	common.RespondWithJSON(w, http.StatusOK, httpResp)
	defer cancel()

}

// UpdateVerificationByUserId handles PUT /users/{user_id}/verifications
func (h *UserProductHandler) UpdateVerificationByUserId(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		common.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	vars := mux.Vars(r)
	userID := vars["user_id"]
	if userID == "" {
		common.RespondWithError(w, http.StatusBadRequest, "Missing user_id")
		return
	}

	var reqBody struct {
		Type   string `json:"type"`
		Status string `json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		common.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	grpcReq := &pb.UpdateVerificationByUserIdRequest{
		UserId: userID,
		Type:   reqBody.Type,
		Status: reqBody.Status,
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)

	grpcResp, err := h.UserProductClient.Client.UpdateVerificationByUserId(ctx, grpcReq)
	if err != nil {
		defer cancel()

		common.RespondGrpcError(w, err)
		return
	}

	httpResp := transformers.UpdateVerificationRespJSON(grpcResp)
	common.RespondWithJSON(w, http.StatusOK, httpResp)
	defer cancel()

}
