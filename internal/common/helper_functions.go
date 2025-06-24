package common

import (
	"encoding/json"
	"log"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Helper function to handle gRPC errors and map them to HTTP responses
func RespondGrpcError(w http.ResponseWriter, err error) {
	st, ok := status.FromError(err)
	if ok {
		switch st.Code() {
		case codes.NotFound:
			RespondWithError(w, http.StatusNotFound, st.Message())
		case codes.InvalidArgument:
			RespondWithError(w, http.StatusBadRequest, st.Message())
		case codes.AlreadyExists:
			RespondWithError(w, http.StatusConflict, st.Message())
		case codes.Unavailable:
			RespondWithError(w, http.StatusServiceUnavailable, "Backend service unavailable: "+err.Error())
		case codes.DeadlineExceeded:
			RespondWithError(w, http.StatusGatewayTimeout, "Backend service timeout: "+err.Error())
		default:
			log.Printf("gRPC error: code=%v, message=%s", st.Code(), st.Message())
			RespondWithError(w, http.StatusBadGateway, "Backend service error: "+err.Error())
		}
	} else {
		log.Printf("Non-gRPC error: %v", err)
		RespondWithError(w, http.StatusInternalServerError, "Internal server error: "+err.Error())
	}
}

// RespondWithJSON writes a JSON response to the HTTP response writer.
func RespondWithJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshaling JSON: %v", err)
		RespondWithError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_, err = w.Write(response)
	if err != nil {
		log.Printf("Error writing response: %v", err)
	}
}

// RespondWithError writes an error JSON response to the HTTP response writer.
func RespondWithError(w http.ResponseWriter, statusCode int, message string) {
	RespondWithJSON(w, statusCode, map[string]string{"error": message})
}
