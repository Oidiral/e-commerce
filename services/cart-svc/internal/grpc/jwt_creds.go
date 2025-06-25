package grpcx

import (
	"context"
	"github.com/oidiral/e-commerce/services/cart-svc/internal/authclient"
)

type JwtCreds struct {
	au *authclient.Client
}

func NewJWTPerRPCCreds(ac *authclient.Client) *JwtCreds {
	return &JwtCreds{au: ac}
}

func (c *JwtCreds) GetRequestMetadata(ctx context.Context, _ ...string) (map[string]string, error) {
	tok, err := c.au.Token(ctx)
	if err != nil {
		return nil, err
	}
	return map[string]string{
		"authorization": "Bearer " + tok,
	}, nil
}

func (*JwtCreds) RequireTransportSecurity() bool { return false }
