package token_validation

import "context"

type ExternalTokenValidation interface {
	ValidateToken(ctx context.Context, token string) (*GetTokekenValidationResponse, error)
}
