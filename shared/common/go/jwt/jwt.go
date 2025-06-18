package jwt

type Claims struct {
	IssuedAt   int64  `json:"iat"`
	Expiration int64  `json:"exp"`
	Subject    string `json:"sub"` // User ID
	Audience   string `json:"aud"` // Session ID
	Issuer     string `json:"iss"` // Issuer (backend/daemon)
	TokenType  string `json:"typ"` // Token type (e.g., "access", "refresh", "mfa")
}

func CreateJWT(claims Claims, secret string) string {
	return ""
}

func VerifyJWT(token string, secret string) (*Claims, error) {
	return nil, nil
}
