package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
)

type Claims struct {
	IssuedAt   int64   `json:"iat"`           // Issued at time
	NotBefore  *int64  `json:"nbf,omitempty"` // Not before time (optional as not all tokens need it)
	Expiration int64   `json:"exp"`           // Expiration time
	Subject    *string `json:"sub,omitempty"` // User ID (optional as MFA session tokens should not include this)
	Audience   string  `json:"aud"`           // Session ID
	Issuer     string  `json:"iss"`           // Issuer (backend/daemon)
	TokenType  string  `json:"typ"`           // Token type (e.g., "access", "refresh", "mfa")
	JTI        *string `json:"jti,omitempty"` // JWT ID - unique identifier for the token (optional as not all tokens need it)
}

func CreateJWT(claims Claims, secret string) (string, error) {
	mapClaims := jwt.MapClaims{
		"iat": claims.IssuedAt,
		"exp": claims.Expiration,
		"aud": claims.Audience,
		"iss": claims.Issuer,
		"typ": claims.TokenType,
	}
	if claims.NotBefore != nil {
		mapClaims["nbf"] = *claims.NotBefore
	}
	if claims.Subject != nil {
		mapClaims["sub"] = *claims.Subject
	}
	if claims.JTI != nil {
		mapClaims["jti"] = *claims.JTI
	}

	// TODO: this needs to be reviewed
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, mapClaims)
	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func VerifyJWT(token string, secret string) (*Claims, error) {
	mapClaims := jwt.MapClaims{}

	// TODO: this needs to be reviewed, I have no clue what any of this does nor whether it is correct
	parsedToken, err := jwt.ParseWithClaims(token, &mapClaims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.Wrap(jwt.ErrSignatureInvalid, "unexpected signing method")
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if !parsedToken.Valid {
		return nil, errors.New("invalid token")
	}

	claims := &Claims{
		IssuedAt:   int64(mapClaims["iat"].(float64)),
		Expiration: int64(mapClaims["exp"].(float64)),
		Audience:   mapClaims["aud"].(string),
		Issuer:     mapClaims["iss"].(string),
		TokenType:  mapClaims["typ"].(string),
	}
	if nbf, ok := mapClaims["nbf"]; ok {
		nbfInt := int64(nbf.(float64))
		claims.NotBefore = &nbfInt
	}
	if sub, ok := mapClaims["sub"]; ok {
		subStr := sub.(string)
		claims.Subject = &subStr
	}
	if jti, ok := mapClaims["jti"]; ok {
		jtiStr := jti.(string)
		claims.JTI = &jtiStr
	}

	return claims, nil
}
