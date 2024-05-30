package handler

import (
	"net/http"
	"super-payment-kun/internal/model"
	"super-payment-kun/internal/pkg"
)

type testLogin struct {
	jwter *pkg.JWTer
}

func NewTestLogin(jwter *pkg.JWTer) *testLogin {
	return &testLogin{
		jwter: jwter,
	}
}

func (h *testLogin) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// NOTE: Normally, we should check if the user exists in the company and if the user has the role to login.
	// but for the sake of testing login, we will skip this step.

	fixedCompanyID := "123e4567-e89b-12d3-a456-426614174000" // 会社A
	fixedUserID := "223e4567-e89b-12d3-a456-426614174000"    // 佐藤太郎
	var fixedRole int8 = int8(model.Admin)                   // admin (he can create invoice and do anything!)

	jwt, err := h.jwter.GenerateToken(ctx, fixedCompanyID, fixedUserID, fixedRole)
	if err != nil {
		DeriveStatusAndRespond(ctx, w, err)
		return
	}
	type data struct {
		// WARN: tricky field name for the sake of testing login!
		AccessToken string `json:"会社A_佐藤太郎_admin_access_token"`
	}

	resp := NormalResp{
		APICode: pkg.Success,
		Data: &data{
			AccessToken: string(jwt),
		},
	}

	RespondWithStatus(ctx, w, resp, http.StatusOK)
}
