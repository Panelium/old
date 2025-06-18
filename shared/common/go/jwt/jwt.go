package jwt

type Claims struct {
	IssuedAt   int64   `json:"iat"`           // Issued at time
	NotBefore  *int64  `json:"nbf"`           // Not before time (optional as not all tokens need it)
	Expiration int64   `json:"exp"`           // Expiration time
	Subject    *string `json:"sub"`           // User ID (optional as MFA session tokens should not include this)
	Audience   string  `json:"aud"`           // Session ID
	Issuer     string  `json:"iss"`           // Issuer (backend/daemon)
	TokenType  string  `json:"typ"`           // Token type (e.g., "access", "refresh", "mfa")
	ID         *string `json:"jti,omitempty"` // JWT ID - unique identifier for the token (optional as not all tokens need it)
}

func CreateJWT(claims Claims, secret string) string {
	return ""
}

func VerifyJWT(token string, secret string) (*Claims, error) {
	return nil, nil
}
