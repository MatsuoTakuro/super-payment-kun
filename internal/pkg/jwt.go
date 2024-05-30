package pkg

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"super-payment-kun/internal/config"

	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

const (
	companyIDKey = "company_id"
	userIDKey    = "user_id"
	roleKey      = "role"
)

// JWTer handles JWT operations with symmetric key (HS256)
type JWTer struct {
	Secret  []byte
	Clocker Clocker
}

func NewJWTer(c Clocker, cfg *config.Config) (*JWTer, error) {

	if c == nil || cfg == nil || cfg.JWTSecret == "" {
		return nil, fmt.Errorf("failed to create jwter")
	}
	return &JWTer{
		Secret:  []byte(cfg.JWTSecret),
		Clocker: c,
	}, nil
}

func (j *JWTer) GenerateToken(ctx context.Context, companyID, userID string, role int8) ([]byte, error) {
	token, err := jwt.NewBuilder().
		JwtID(uuid.New().String()).
		Issuer(`github.com/MatsuoTakuro/super-payment-kun`).
		Subject("access_token").
		IssuedAt(j.Clocker.Now()).
		Expiration(j.Clocker.Now().Add(24*time.Hour)). // 1 day long for testing login
		Claim(companyIDKey, companyID).
		Claim(userIDKey, userID).
		Claim(roleKey, role).
		Build()
	if err != nil {
		return nil, NewAPIError(Unknown, err, "failed to build token")
	}
	// TODO: save token to memory db like redis

	signed, err := jwt.Sign(token, jwt.WithKey(jwa.HS256, j.Secret))
	if err != nil {
		return nil, NewAPIError(Unknown, err, "failed to sign token")
	}
	return signed, nil
}

func (j *JWTer) Validate(r *http.Request) (*http.Request, error) {
	tok, err := j.getAndValidateToken(r)
	if err != nil {
		return nil, err
	}
	companyID, ok := tok.Get(companyIDKey)
	if !ok {
		err := fmt.Errorf("failed to get company_id from token")
		return nil, NewAPIError(Unknown, err, err.Error())
	}
	ctx := SetCompanyID(r.Context(), companyID.(string))

	userID, ok := tok.Get(userIDKey)
	if !ok {
		err := fmt.Errorf("failed to get user_id from token")
		return nil, NewAPIError(Unknown, err, err.Error())
	}
	ctx = SetUserID(ctx, userID.(string))

	role, ok := tok.Get(roleKey)
	if !ok {
		err := fmt.Errorf("failed to get role from token")
		return nil, NewAPIError(Unknown, err, err.Error())
	}
	ctx = SetRole(ctx, int8(role.(float64)))

	return r.Clone(ctx), nil
}

func (j *JWTer) getAndValidateToken(r *http.Request) (jwt.Token, error) {
	token, err := jwt.ParseRequest(
		r,
		jwt.WithKey(jwa.HS256, j.Secret),
		jwt.WithValidate(false),
	)
	if err != nil {
		return nil, NewAPIError(Unknown, err, "failed to get token from request")
	}
	if err := jwt.Validate(token, jwt.WithClock(j.Clocker)); err != nil {
		return nil, NewAPIError(Unauthorized, err, "failed to validate token")
	}
	// TODO: check token from memory db like redis to make sure it doesn't expire

	return token, nil
}
