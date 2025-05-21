package transformers

import pb "github.com/software-architecture-proj/nova-backend-common-protos/gen/go/transaction_service"

func TransferFundsRespJSON(resp *pb.TransferFundsResponse) map[string]interface{} {
	return map[string]interface{}{
		"success":     resp.GetSuccess(),
		"message":     resp.GetMessage(),
		"transfer_id": resp.GetTransferId(),
		"timestamp":   resp.GetTimestamp(),
	}
}

func CreateAccountRespJSON(resp *pb.CreateAccountResponse) map[string]interface{} {
	return map[string]interface{}{
		"success":   resp.GetSuccess(),
		"message":   resp.GetMessage(),
		"user_id":   resp.GetUserId(),
		"timestamp": resp.GetTimestamp(),
	}
}

func GetBalanceRespJSON(resp *pb.GetBalanceResponse) map[string]interface{} {
	balances := make([]map[string]string, len(resp.GetBalances()))
	for i, gbResult := range resp.GetBalances() {
		balances[i] = map[string]string{
			"income":  gbResult.GetIncome(),
			"outcome": gbResult.GetOutcome(),
		}
	}

	return map[string]interface{}{
		"success":   resp.GetSuccess(),
		"message":   resp.GetMessage(),
		"timestamp": resp.GetTimestamp(),
		"current":   resp.GetCurrent(),
		"balances":  balances,
	}
}

func GetMovementsRespJSON(resp *pb.GetMovementsResponse) map[string]interface{} {
	movements := make([]map[string]string, len(resp.GetMovements()))
	for i, gtResult := range resp.GetMovements() {
		movements[i] = map[string]string{
			"transferId":   gtResult.GetTransferId(),
			"fromUsername": gtResult.GetFromUsername(),
			"toUsername":   gtResult.GetToUsername(),
			"amount":       gtResult.GetAmount(),
			"timestamp":    gtResult.GetTimestamp(),
		}
	}

	return map[string]interface{}{
		"success":   resp.GetSuccess(),
		"message":   resp.GetMessage(),
		"movements": movements,
	}
}
