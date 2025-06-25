package jwt

import (
	"crypto/rsa"
	stdErrors "errors"
	"github.com/golang-jwt/jwt/v5"
	"panelium/common/errors"
)

type TokenType string

const (
	AccessTokenType  TokenType = "access"
	RefreshTokenType TokenType = "refresh"
	MFATokenType     TokenType = "mfa"
)

type Issuer string

const (
	BackendIssuer Issuer = "backend"
	DaemonIssuer  Issuer = "daemon"
)

type Claims struct {
	IssuedAt   int64     `json:"iat"`           // Issued at time
	NotBefore  *int64    `json:"nbf,omitempty"` // Not before time (optional as not all tokens need it)
	Expiration int64     `json:"exp"`           // Expiration time
	Subject    *string   `json:"sub,omitempty"` // User ID (optional as MFA session tokens should not include this)
	Audience   string    `json:"aud"`           // Session ID
	Issuer     Issuer    `json:"iss"`           // Issuer (backend/daemon)
	TokenType  TokenType `json:"typ"`           // Token type (e.g., "access", "refresh", "mfa")
	JTI        string    `json:"jti"`           // JWT ID - unique identifier for the token
}

func CreateJWT(claims Claims, key *rsa.PrivateKey) (string, error) {
	mapClaims := jwt.MapClaims{
		"iat": claims.IssuedAt,
		"exp": claims.Expiration,
		"aud": claims.Audience,
		"iss": claims.Issuer,
		"typ": claims.TokenType,
		"jti": claims.JTI,
	}
	if claims.NotBefore != nil {
		mapClaims["nbf"] = *claims.NotBefore
	}
	if claims.Subject != nil {
		mapClaims["sub"] = *claims.Subject
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, mapClaims)
	signedToken, err := token.SignedString(key)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func VerifyJWT(token string, key *rsa.PublicKey) (*Claims, error) {
	// check JTI against database to prevent replay attacks (if not exist, delete session - logout)
	// check if not before nbf or iat
	// check if not after exp
	// check if audience is valid
	// check if issuer is valid
	// check if token type is valid
	// check if token is signed with the correct algorithm
	// check if token is signed with the correct secret

	mapClaims := jwt.MapClaims{}

	// https://pkg.go.dev/github.com/golang-jwt/jwt/v5#example-Parse-Hmac

	// TODO: this needs to be reviewed, I have no clue what any of this does nor whether it is correct
	parsedToken, err := jwt.ParseWithClaims(token, &mapClaims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, stdErrors.Join(jwt.ErrSignatureInvalid, errors.InvalidCredentials)
		}
		return key, nil
	})

	if err != nil {
		return nil, err
	}

	if !parsedToken.Valid {
		return nil, errors.InvalidCredentials
	}

	claims := &Claims{
		IssuedAt:   int64(mapClaims["iat"].(float64)),
		Expiration: int64(mapClaims["exp"].(float64)),
		Audience:   mapClaims["aud"].(string),
		Issuer:     mapClaims["iss"].(Issuer),
		TokenType:  mapClaims["typ"].(TokenType),
		JTI:        mapClaims["jti"].(string),
	}
	if nbf, ok := mapClaims["nbf"]; ok {
		nbfInt := int64(nbf.(float64))
		claims.NotBefore = &nbfInt
	}
	if sub, ok := mapClaims["sub"]; ok {
		subStr := sub.(string)
		claims.Subject = &subStr
	}

	return claims, nil
}
