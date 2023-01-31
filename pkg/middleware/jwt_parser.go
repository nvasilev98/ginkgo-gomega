package middleware

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

//go:generate mockgen --source=jwt_parser.go --destination mocks/jwt_parser.go --package mocks

const authorization = "Authorization"

type JWTValidator interface {
	ValidateJWT(ctx context.Context, rawToken string) error
}

type JWTParser struct {
	claims    []string
	validator JWTValidator
}

type Option func(*JWTParser)

// NewJWTParserWithOptions is a constructor function
func NewJWTParserWithOptions(options ...Option) *JWTParser {
	parser := &JWTParser{}
	for _, opt := range options {
		opt(parser)
	}
	return parser
}

// WithClaims is an option for the parser
func WithClaims(claims []string) Option {
	return func(p *JWTParser) {
		p.claims = claims
	}
}

// WithValidator is an option for the parser
func WithValidator(validator JWTValidator) Option {
	return func(p *JWTParser) {
		p.validator = validator
	}
}

// Handle is a middleware for parsing JWT token with provided options
func (p *JWTParser) Handle(c *gin.Context) {
	tokenB64, err := extractRawToken(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "Authentication failed")
		c.Abort()
		return
	}

	if err := p.validateJWT(c.Request.Context(), tokenB64); err != nil {
		c.JSON(http.StatusUnauthorized, "Invalid JWT token")
		c.Abort()
		return
	}

	payload, err := decodePayload(tokenB64)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "Parsing token failed")
		c.Abort()
		return
	}

	for _, claim := range p.claims {
		c.Set(claim, payload[claim])
	}

	c.Next()
}

func (p *JWTParser) validateJWT(ctx context.Context, token string) error {
	if p.validator != nil {
		if err := p.validator.ValidateJWT(ctx, token); err != nil {
			return fmt.Errorf("failed to validate jwt token: %w", err)
		}
	}
	return nil
}

func extractRawToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get(authorization)

	if authHeader != "" {
		splitAuthHeader := strings.Fields(strings.TrimSpace(authHeader))
		if strings.EqualFold(splitAuthHeader[0], "bearer") && len(splitAuthHeader) == 2 {
			return splitAuthHeader[1], nil
		}
	}

	return "", fmt.Errorf("extracting token from request header failed")
}

func decodePayload(tokenB64 string) (map[string]interface{}, error) {
	jwtParts := strings.Split(tokenB64, ".")
	if len(jwtParts) < 3 {
		return nil, fmt.Errorf("failed to split token")
	}

	payloadJSON, err := base64.RawURLEncoding.DecodeString(jwtParts[1])
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64 string: %w", err)
	}

	var payload map[string]interface{}
	if err = json.Unmarshal(payloadJSON, &payload); err != nil {
		return nil, fmt.Errorf("failed to unmarshal json: %w", err)
	}

	return payload, nil
}
