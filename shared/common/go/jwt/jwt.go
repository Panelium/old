package jwt

import (
	"crypto/rsa"
	stdErrors "errors"
	"github.com/golang-jwt/jwt/v5"
	"panelium/common/errors"
	"time"
)

type TokenType string

const (
	AccessTokenType        TokenType = "access"
	RefreshTokenType       TokenType = "refresh"
	PasswordResetTokenType TokenType = "reset"
	MFATokenType           TokenType = "mfa"
)

type Issuer string // TODO: this might be changed to a url

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

func VerifyJWT(token string, key *rsa.PublicKey, expectedIssuer Issuer, expectedTokenType TokenType) (*Claims, error) {
	parser := jwt.NewParser(jwt.WithIssuedAt(), jwt.WithIssuer(string(expectedIssuer)), jwt.WithExpirationRequired(), jwt.WithValidMethods([]string{jwt.SigningMethodRS256.Alg()}))

	mapClaims := jwt.MapClaims{}

	parsedToken, err := parser.ParseWithClaims(token, &mapClaims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.InvalidCredentials
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

	if claims.Issuer != expectedIssuer {
		return nil, errors.InvalidCredentials
	}

	if claims.TokenType != expectedTokenType {
		return nil, errors.InvalidCredentials
	}

	if claims.Audience == "" {
		return nil, errors.InvalidCredentials
	}

	if (claims.Subject == nil || *claims.Subject == "") && expectedTokenType != MFATokenType {
		return nil, errors.InvalidCredentials
	}

	if claims.IssuedAt <= 0 || claims.Expiration <= 0 || claims.Expiration <= claims.IssuedAt {
		return nil, errors.InvalidCredentials
	}

	if claims.NotBefore != nil && time.Now().Unix() < *claims.NotBefore {
		return nil, stdErrors.New("token not valid yet (nbf)") //TODO - refactor
	}

	if time.Now().Unix() < claims.IssuedAt {
		return nil, stdErrors.New("token issued in the future (iat)") //TODO - refactor
	}

	if time.Now().Unix() > claims.Expiration {
		return nil, stdErrors.New("token expired (exp)") //TODO - refactor
	}

	return claims, nil
}
