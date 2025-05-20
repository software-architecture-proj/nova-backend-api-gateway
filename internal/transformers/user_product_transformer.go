package transformers

import pb "github.com/software-architecture-proj/nova-backend-common-protos/gen/go/user_product_service"

func CreateUserRespJSON(resp *pb.CreateUserResponse) map[string]interface{} {
	return map[string]interface{}{
		"success": resp.GetSuccess(),
		"message": resp.GetMessage(),
		"user_id": resp.GetUserId(),
	}
}

func GetUserRespJSON(resp *pb.GetUserByIdResponse) map[string]interface{} {
	return map[string]interface{}{
		"success":     resp.GetSuccess(),
		"message":     resp.GetMessage(),
		"email":       resp.GetEmail(),
		"username":    resp.GetUsername(),
		"phone":       resp.GetPhone(),
		"first_name":  resp.GetFirstName(),
		"last_name":   resp.GetLastName(),
		"birthdate":   resp.GetBirthdate(),
	}
}

func UpdateUserRespJSON(resp *pb.UpdateUserByIdResponse) map[string]interface{} {
	return map[string]interface{}{
		"success":     resp.GetSuccess(),
		"message":     resp.GetMessage(),
		"email":       resp.GetEmail(),
		"username":    resp.GetUsername(),
		"phone":       resp.GetPhone(),
		"first_name":  resp.GetFirstName(),
		"last_name":   resp.GetLastName(),
		"birthdate":   resp.GetBirthdate(),
	}
}

func DeleteUserRespJSON(resp *pb.DeleteUserByIdResponse) map[string]interface{} {
	return map[string]interface{}{
		"success": resp.GetSuccess(),
		"message": resp.GetMessage(),
	}
}

func GetFavoritesRespJSON(resp *pb.GetFavoritesByUserIdResponse) map[string]interface{} {
	var favorites []map[string]interface{}
	for _, f := range resp.GetFavorites() {
		favorites = append(favorites, map[string]interface{}{
			"id":                f.GetId(),
			"user_id":           f.GetUserId(),
			"favorite_user_id":  f.GetFavoriteUserId(),
			"favorite_username": f.GetFavoriteUsername(),
			"alias":             f.GetAlias(),
		})
	}
	return map[string]interface{}{
		"success":   resp.GetSuccess(),
		"message":   resp.GetMessage(),
		"favorites": favorites,
	}
}

func CreateFavoriteRespJSON(resp *pb.CreateFavoriteResponse) map[string]interface{} {
	return map[string]interface{}{
		"success":     resp.GetSuccess(),
		"message":     resp.GetMessage(),
		"favorite_id": resp.GetFavoriteId(),
	}
}

func UpdateFavoriteRespJSON(resp *pb.UpdateFavoriteByIdResponse) map[string]interface{} {
	return map[string]interface{}{
		"success":    resp.GetSuccess(),
		"message":    resp.GetMessage(),
		"new_alias":  resp.GetNewAlias(),
	}
}

func DeleteFavoriteRespJSON(resp *pb.DeleteFavoriteByIdResponse) map[string]interface{} {
	return map[string]interface{}{
		"success": resp.GetSuccess(),
		"message": resp.GetMessage(),
	}
}

func GetPocketsRespJSON(resp *pb.GetPocketsByUserIdResponse) map[string]interface{} {
	var pockets []map[string]interface{}
	for _, p := range resp.GetPockets() {
		pockets = append(pockets, map[string]interface{}{
			"id":         p.GetId(),
			"user_id":    p.GetUserId(),
			"name":       p.GetName(),
			"category":   p.GetCategory(),
			"max_amount": p.GetMaxAmount(),
		})
	}
	return map[string]interface{}{
		"success": resp.GetSuccess(),
		"message": resp.GetMessage(),
		"pockets": pockets,
	}
}

func CreatePocketRespJSON(resp *pb.CreatePocketResponse) map[string]interface{} {
	return map[string]interface{}{
		"success":    resp.GetSuccess(),
		"message":    resp.GetMessage(),
		"pocket_id":  resp.GetPocketId(),
	}
}

func UpdatePocketRespJSON(resp *pb.UpdatePocketByIdResponse) map[string]interface{} {
	return map[string]interface{}{
		"success":    resp.GetSuccess(),
		"message":    resp.GetMessage(),
		"name":       resp.GetName(),
		"category":   resp.GetCategory(),
		"max_amount": resp.GetMaxAmount(),
	}
}

func DeletePocketRespJSON(resp *pb.DeletePocketByIdResponse) map[string]interface{} {
	return map[string]interface{}{
		"success": resp.GetSuccess(),
		"message": resp.GetMessage(),
	}
}

func GetVerificationsRespJSON(resp *pb.GetVerificationsByUserIdResponse) map[string]interface{} {
	var verifications []map[string]interface{}
	for _, v := range resp.GetVerifications() {
		verifications = append(verifications, map[string]interface{}{
			"id":       v.GetId(),
			"user_id":  v.GetUserId(),
			"type":     v.GetType(),
			"status":   v.GetStatus(),
		})
	}
	return map[string]interface{}{
		"success":       resp.GetSuccess(),
		"message":       resp.GetMessage(),
		"verifications": verifications,
	}
}

func UpdateVerificationRespJSON(resp *pb.UpdateVerificationByUserIdResponse) map[string]interface{} {
	return map[string]interface{}{
		"success": resp.GetSuccess(),
		"message": resp.GetMessage(),
		"type":    resp.GetType(),
		"status":  resp.GetStatus(),
	}
}

func GetCountryCodesRespJSON(resp *pb.GetCountryCodesResponse) map[string]interface{} {
	var codes []map[string]interface{}
	for _, c := range resp.GetCodes() {
		codes = append(codes, map[string]interface{}{
			"id":      c.GetId(),
			"name":    c.GetName(),
			"code":    c.GetCode(),
		})
	}
	return map[string]interface{}{
		"success":       resp.GetSuccess(),
		"message":       resp.GetMessage(),
		"country_codes": codes,
	}
}
