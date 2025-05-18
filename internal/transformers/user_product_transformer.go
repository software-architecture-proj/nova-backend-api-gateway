
package transformers

import (
	pb "github.com/software-architecture-proj/nova-backend-common-protos/gen/go/user_product_service" 
)

// ToUserJSON transforms a gRPC User object to a JSON-friendly structure.
func ToUserJSON(user *pb.User) map[string]interface{} {
	return map[string]interface{}{
		"id":         user.GetId(),
		"email":      user.GetEmail(),
		"username":   user.GetUsername(),
		"phone":      user.GetPhone(),
		"firstName":  user.GetFirstName(),
		"lastName":   user.GetLastName(),
		"birthdate":  user.GetBirthdate(),
	}
}

// ToFavoriteJSON transforms a gRPC Favorite object to a JSON-friendly structure
func ToFavoriteJSON(favorite *pb.Favorite) map[string]interface{} {
	return map[string]interface{}{
		"id":               favorite.GetId(),
		"user_id":          favorite.GetUserId(),
		"favorite_user_id": favorite.GetFavoriteUserId(),
		"alias":            favorite.GetAlias(),
	}
}

// ToFavoriteListJSON transforms a slice of gRPC Favorite objects to a JSON-friendly structure.
func ToFavoriteListJSON(favorites []*pb.Favorite) []map[string]interface{} {
	result := make([]map[string]interface{}, len(favorites))
	for i, fav := range favorites {
		result[i] = ToFavoriteJSON(fav)
	}
	return result
}

// ToPocketJSON transforms a gRPC Pocket object to a JSON-friendly structure.
func ToPocketJSON(pocket *pb.Pocket) map[string]interface{} {
	return map[string]interface{}{
		"id":          pocket.GetId(),
		"user_id":     pocket.GetUserId(),
        "name":        pocket.GetName(),
		"category":    pocket.GetCategory(),
		"max_amount":  pocket.GetMaxAmount(),
	}
}

// ToPocketListJSON transforms a slice of gRPC Pocket objects to a JSON-friendly structure.
func ToPocketListJSON(pockets []*pb.Pocket) []map[string]interface{} {
	result := make([]map[string]interface{}, len(pockets))
	for i, pocket := range pockets {
		result[i] = ToPocketJSON(pocket)
	}
	return result
}
