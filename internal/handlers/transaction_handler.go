package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/software-architecture-proj/nova-backend-api-gateway/internal/clients"
	"github.com/software-architecture-proj/nova-backend-api-gateway/internal/common"
	"github.com/software-architecture-proj/nova-backend-api-gateway/internal/transformers"

	// Import from common-protos
	pb "github.com/software-architecture-proj/nova-backend-common-protos/gen/go/transaction_service"
)

type TransactionHandler struct {
	TransactionClient *clients.TransactionServiceClient
}

func NewTransactionHandler(TransactionClient *clients.TransactionServiceClient) *TransactionHandler {
	return &TransactionHandler{TransactionClient: TransactionClient}
}

// GetMovements handles GET /movements
func (h *TransactionHandler) GetMovements(w http.ResponseWriter, r *http.Request) {
	var reqBody struct {
		User_id   string `json:"user_id"`
		From_time string `json:"from_time"`
		To_time   string `json:"to_time"`
		Limit     bool   `json:"limit"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		common.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if reqBody.User_id == "" {
		common.RespondWithError(w, http.StatusBadRequest, "Missing user_id")
		return
	}

	fromTime, err := strconv.ParseUint(reqBody.From_time, 10, 64)
	if err != nil {
		common.RespondWithError(w, http.StatusBadRequest, "Invalid 'from' time format")
		return
	}
	toTime, err := strconv.ParseUint(reqBody.To_time, 10, 64)
	if err != nil {
		common.RespondWithError(w, http.StatusBadRequest, "Invalid 'to' time format")
		return
	}

	grpcReq := &pb.GetMovementsRequest{UserId: reqBody.User_id, FromTime: fromTime, ToTime: toTime, Limit: reqBody.Limit}
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	grpcResp, err := h.TransactionClient.Client.Movements(ctx, grpcReq)
	if err != nil {
		common.RespondGrpcError(w, err)
		return
	}

	httpResp := transformers.GetMovementsRespJSON(grpcResp)
	common.RespondWithJSON(w, http.StatusOK, httpResp)
}

// PostAccount CreateUser handles POST /account
func (h *TransactionHandler) PostAccount(w http.ResponseWriter, r *http.Request) {
	var reqBody struct {
		Username string `json:"username"`
		Bank     bool   `json:"bank"`
		UserId   string `json:"user_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		common.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	grpcReq := &pb.CreateAccountRequest{
		UserId:   reqBody.UserId,
		Bank:     reqBody.Bank,
		Username: reqBody.Username,
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	grpcResp, err := h.TransactionClient.Client.Account(ctx, grpcReq)
	if err != nil {
		common.RespondGrpcError(w, err)
		return
	}

	// Assuming you have a transformer function to convert the gRPC response to a suitable HTTP response
	httpResp := transformers.CreateAccountRespJSON(grpcResp) //  Create this function in transformers.go
	common.RespondWithJSON(w, http.StatusCreated, httpResp)
}

// GetBalance handles GET /balance
func (h *TransactionHandler) GetBalance(w http.ResponseWriter, r *http.Request) {
	var reqBody struct {
		User_id   string `json:"user_id"`
		From_time string `json:"from_time"`
		To_time   string `json:"to_time"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		common.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if reqBody.User_id == "" {
		common.RespondWithError(w, http.StatusBadRequest, "Missing user_id")
		return
	}

	fromTime, err := strconv.ParseUint(reqBody.From_time, 10, 64)
	if err != nil {
		common.RespondWithError(w, http.StatusBadRequest, "Invalid 'from' time format")
		return
	}
	toTime, err := strconv.ParseUint(reqBody.To_time, 10, 64)
	if err != nil {
		common.RespondWithError(w, http.StatusBadRequest, "Invalid 'to' time format")
		return
	}

	grpcReq := &pb.GetBalanceRequest{UserId: reqBody.User_id, FromTime: fromTime, ToTime: toTime}
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	grpcResp, err := h.TransactionClient.Client.Balance(ctx, grpcReq)
	if err != nil {
		common.RespondGrpcError(w, err)
		return
	}

	httpResp := transformers.GetBalanceRespJSON(grpcResp) // Create this function in transformers.go
	common.RespondWithJSON(w, http.StatusOK, httpResp)
}

// PostTransfer handles POST /transfer
func (h *TransactionHandler) PostTransfer(w http.ResponseWriter, r *http.Request) {
	var reqBody struct {
		FromUser string `json:"from_user"`
		ToUser   string `json:"to_user"`
		Amount   uint64 `json:"amount"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		common.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	grpcReq := &pb.TransferFundsRequest{
		FromUserId: reqBody.FromUser,
		ToUserId:   reqBody.ToUser,
		Amount:     reqBody.Amount,
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	grpcResp, err := h.TransactionClient.Client.Transfer(ctx, grpcReq)
	if err != nil {
		common.RespondGrpcError(w, err)
		return
	}

	// Assuming you have a transformer function to convert the gRPC response to a suitable HTTP response
	httpResp := transformers.TransferFundsRespJSON(grpcResp) //  Create this function in transformers.go
	common.RespondWithJSON(w, http.StatusCreated, httpResp)
}
