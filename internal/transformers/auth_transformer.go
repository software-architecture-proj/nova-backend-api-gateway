package transformers

import pb "github.com/software-architecture-proj/nova-backend-common-protos/gen/go/auth_service"

func LoginRespJSON(resp *pb.LoginResponse) map[string]interface{} {
	return map[string]interface{}{
		"success": resp.GetSuccess(),
		"message": resp.GetMessage(),
	}
}
