package session

import (
	"panelium/backend/internal/config"
	"panelium/common/id"
	"panelium/common/jwt"
	"time"
)

func CreateRefreshToken(issuedAt time.Time, sessionId string, uid string) (token string, JTI string, expiration time.Time, err error) {
	return CreateToken(issuedAt, config.RefreshTokenDuration, jwt.BackendIssuer, jwt.RefreshTokenType, sessionId, uid)
}

func CreateAccessToken(issuedAt time.Time, sessionId string, uid string) (token string, JTI string, expiration time.Time, err error) {
	return CreateToken(issuedAt, config.AccessTokenDuration, jwt.BackendIssuer, jwt.AccessTokenType, sessionId, uid)
}

func CreateToken(issuedAt time.Time, duration time.Duration, issuer jwt.Issuer, tokenType jwt.TokenType, sessionId string, uid string) (token string, JTI string, expiration time.Time, err error) {
	JTI, err = id.New()
	if err != nil {
		return "", "", time.Time{}, err
	}

	tokenExpiration := issuedAt.Add(duration)

	tokenClaims := jwt.Claims{
		IssuedAt:   issuedAt.Unix(),
		Expiration: tokenExpiration.Unix(),
		Subject:    &uid,
		Audience:   sessionId,
		Issuer:     issuer,
		TokenType:  tokenType,
		JTI:        JTI,
	}

	token, err = jwt.CreateJWT(tokenClaims, config.JWTSecret)
	if err != nil {
		return "", "", time.Time{}, err
	}

	return token, JTI, tokenExpiration, nil
}
