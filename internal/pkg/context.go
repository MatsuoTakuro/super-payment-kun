package pkg

import (
	"context"
	"errors"
)

type ctxKeyCompanyID struct{}

func SetCompanyID(ctx context.Context, companyID string) context.Context {
	return context.WithValue(ctx, ctxKeyCompanyID{}, companyID)
}

func GetCompanyID(ctx context.Context) (string, error) {
	if companyID, ok := ctx.Value(ctxKeyCompanyID{}).(string); ok {
		return companyID, nil
	}
	return "", NewAPIError(Unknown, errors.New("failed to get company_id from context"))
}

type ctxKeyUserID struct{}

func SetUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, ctxKeyUserID{}, userID)
}

func GetUserID(ctx context.Context) (string, error) {
	if userID, ok := ctx.Value(ctxKeyUserID{}).(string); ok {
		return userID, nil
	}
	return "", NewAPIError(Unknown, errors.New("failed to get user_id from context"))
}

type ctxKeyRole struct{}

func SetRole(ctx context.Context, role int8) context.Context {
	return context.WithValue(ctx, ctxKeyRole{}, role)
}

func GetRole(ctx context.Context) (int8, error) {
	role, ok := ctx.Value(ctxKeyRole{}).(int8)
	if ok {
		return role, nil
	}
	return -1, NewAPIError(Unknown, errors.New("failed to get role from context"))
}
