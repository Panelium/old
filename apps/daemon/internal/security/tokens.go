package security

import (
	"panelium/common/id"
	"panelium/common/jwt"
	"panelium/daemon/internal/config"
	"time"
)

const nodeTokenDuration = 100 * 365 * 24 * time.Hour // 100 years
// TODO: this needs to be changed to a more reasonable value once token rotation is implemented

func CreateNodeToken(issuedAt time.Time) (token string, JTI string, expiration time.Time, err error) {
	return CreateToken(issuedAt, nodeTokenDuration, jwt.DaemonIssuer, jwt.NodeTokenType, nil, nil)
}

func CreateToken(issuedAt time.Time, duration time.Duration, issuer jwt.Issuer, tokenType jwt.TokenType, sessionId *string, uid *string) (token string, JTI string, expiration time.Time, err error) {
	JTI, err = id.New()
	if err != nil {
		return "", "", time.Time{}, err
	}

	tokenExpiration := issuedAt.Add(duration)

	tokenClaims := jwt.Claims{
		IssuedAt:   issuedAt.Unix(),
		Expiration: tokenExpiration.Unix(),
		Subject:    uid,
		Audience:   sessionId,
		Issuer:     issuer,
		TokenType:  tokenType,
		JTI:        JTI,
	}

	token, err = jwt.CreateJWT(tokenClaims, config.JWTPrivateKeyInstance)
	if err != nil {
		return "", "", time.Time{}, err
	}

	return token, JTI, tokenExpiration, nil
}
