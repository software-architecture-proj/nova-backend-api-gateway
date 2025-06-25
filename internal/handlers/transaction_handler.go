package handlers

import (
	"context"
	"encoding/json"
	"log"
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
	UserProductClient *clients.UserProductServiceClient
}

func NewTransactionHandler(TransactionClient *clients.TransactionServiceClient, userClient *clients.UserProductServiceClient) *TransactionHandler {
	return &TransactionHandler{
		TransactionClient: TransactionClient,
		UserProductClient: userClient,
	}
}

// GetMovements handles GET /movements
func (h *TransactionHandler) GetMovements(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		common.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	query := r.URL.Query()
	userId := query.Get("id")
	fromTimeStr := query.Get("from")
	toTimeStr := query.Get("to")
	limitStr := query.Get("lim")

	if userId == "" {
		common.RespondWithError(w, http.StatusBadRequest, "Missing user_id")
		return
	}

	var fromTime, toTime uint64
	var err error

	if fromTimeStr == "" && toTimeStr == "" {
		return
	} else {
		fromTime, err = strconv.ParseUint(fromTimeStr, 10, 64)
		if err != nil {
			common.RespondWithError(w, http.StatusBadRequest, "Invalid 'from' time format")
			return
		}
		toTime, err = strconv.ParseUint(toTimeStr, 10, 64)
		if err != nil {
			common.RespondWithError(w, http.StatusBadRequest, "Invalid 'to' time format")
			return
		}
	}

	limit := false
	if limitStr != "" {
		limit, err = strconv.ParseBool(limitStr)
		if err != nil {
			common.RespondWithError(w, http.StatusBadRequest, "Invalid limit format")
			return
		}
	}

	grpcReq := &pb.GetMovementsRequest{UserId: userId, FromTime: fromTime, ToTime: toTime, Limit: limit}
	log.Println("GetMovements called with userId:", userId, "fromTime:", fromTime, "toTime:", toTime, "limit:", limit)
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)

	grpcResp, err := h.TransactionClient.Client.Movements(ctx, grpcReq)
	if err != nil {
		defer cancel()
		common.RespondGrpcError(w, err)
		return
	}

	httpResp := transformers.GetMovementsRespJSON(grpcResp)
	common.RespondWithJSON(w, http.StatusOK, httpResp)
	defer cancel()
}

// PostAccount handles POST /account
func (h *TransactionHandler) PostAccount(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		common.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

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
	grpcResp, err := h.TransactionClient.Client.Account(ctx, grpcReq)
	if err != nil {
		common.RespondGrpcError(w, err)
		defer cancel()
		return
	}

	httpResp := transformers.CreateAccountRespJSON(grpcResp)
	common.RespondWithJSON(w, http.StatusCreated, httpResp)
	defer cancel()
}

// GetBalance handles GET /balance
func (h *TransactionHandler) GetBalance(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		common.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	query := r.URL.Query()
	userId := query.Get("user_id")
	fromTimeStr := query.Get("from_time")
	toTimeStr := query.Get("to_time")

	if userId == "" {
		common.RespondWithError(w, http.StatusBadRequest, "Missing user_id")
		return
	}

	fromTime, err := strconv.ParseUint(fromTimeStr, 10, 64)
	if err != nil {
		common.RespondWithError(w, http.StatusBadRequest, "Invalid 'from' time format")
		return
	}
	toTime, err := strconv.ParseUint(toTimeStr, 10, 64)
	if err != nil {
		common.RespondWithError(w, http.StatusBadRequest, "Invalid 'to' time format")
		return
	}

	grpcReq := &pb.GetBalanceRequest{UserId: userId, FromTime: fromTime, ToTime: toTime}
	log.Println("GetBalance called with userId:", userId, "fromTime:", fromTime, "toTime:", toTime)
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	grpcResp, err := h.TransactionClient.Client.Balance(ctx, grpcReq)
	if err != nil {
		defer cancel()
		log.Println("Error calling Balance method:", err)
		common.RespondGrpcError(w, err)
		return
	}
	httpResp := transformers.GetBalanceRespJSON(grpcResp)
	common.RespondWithJSON(w, http.StatusOK, httpResp)
	defer cancel()
}

// PostTransfer handles POST /transfer
func (h *TransactionHandler) PostTransfer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		common.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var reqBody struct {
		FromUser string `json:"from_user"`
		ToUser   string `json:"to_user"`
		Amount   uint64 `json:"amount"`
		Email    string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		common.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Get user email from user service
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)

	grpcReq := &pb.TransferFundsRequest{
		FromUserId:    reqBody.FromUser,
		ToUserId:      reqBody.ToUser,
		Amount:        reqBody.Amount,
		FromUserEmail: reqBody.Email,
	}

	grpcResp, err := h.TransactionClient.Client.Transfer(ctx, grpcReq)
	if err != nil {
		defer cancel()
		common.RespondGrpcError(w, err)
		return
	}

	httpResp := transformers.TransferFundsRespJSON(grpcResp)
	common.RespondWithJSON(w, http.StatusCreated, httpResp)
	defer cancel()
}
