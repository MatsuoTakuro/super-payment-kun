package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"super-payment-kun/internal/pkg"
)

type NormalResp struct {
	pkg.APICode `json:"api_code"`
	Data        any `json:"data"`
}

func RespondWithStatus(ctx context.Context, w http.ResponseWriter, resp any, status int) {

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	respJSON, err := json.Marshal(resp)
	if err != nil {
		log.Printf("encode response error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)

		rsp := pkg.NewAPIError(pkg.Unknown, err, "failed to encode response")
		if err := json.NewEncoder(w).Encode(rsp); err != nil {
			log.Printf("write error response error: %v", err)
		}
		return
	}

	w.WriteHeader(status)

	if _, err := fmt.Fprintf(w, "%s", respJSON); err != nil {
		log.Printf("write response error: %v", err)
	}

	err, ok := resp.(error)
	if !ok {
		log.Printf("status: %d, resp_json: %s\n", status, respJSON)
	} else {
		log.Printf("status: %d, resp_json: %s, error: %v\n", status, respJSON, err)
	}
}

// DeriveStatusAndRespond derives status from app code if error is AppError, and responds with the status
func DeriveStatusAndRespond(ctx context.Context, w http.ResponseWriter, err error) {

	var appErr *pkg.APIError
	if !errors.As(err, &appErr) {
		log.Printf("not found app_error: %v", err)
		RespondWithStatus(ctx, w, err, http.StatusInternalServerError)
		return
	}

	// Derive status from app code and pass it to Respond
	switch appErr.APICode {
	case pkg.DueDateIsPassed, pkg.DueDateExceedMaxDeferDate, pkg.PaymentAmountTooSmall:
		RespondWithStatus(ctx, w, appErr, http.StatusBadRequest)
		return

	case pkg.PermissionDenied:
		RespondWithStatus(ctx, w, appErr, http.StatusForbidden)

		// TODO: Add more cases here

	default:
		RespondWithStatus(ctx, w, appErr, http.StatusInternalServerError)
		return
	}
}
