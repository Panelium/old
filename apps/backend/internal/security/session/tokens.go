package session

import (
	"panelium/backend/internal/config"
	"panelium/common/id"
	"panelium/common/jwt"
	"time"
)

func CreateRefreshToken(issuedAt time.Time, sessionId string, uid string) (token string, JTI string, expiration time.Time, err error) {
	return CreateToken(issuedAt, config.ConfigInstance.GetRefreshTokenDuration(), jwt.BackendIssuer, jwt.RefreshTokenType, &sessionId, &uid)
}

func CreateAccessToken(issuedAt time.Time, sessionId string, uid string) (token string, JTI string, expiration time.Time, err error) {
	return CreateToken(issuedAt, config.ConfigInstance.GetAccessTokenDuration(), jwt.BackendIssuer, jwt.AccessTokenType, &sessionId, &uid)
}

func CreatePasswordResetToken(issuedAt time.Time, sessionId string) (token string, JTI string, expiration time.Time, err error) {
	return CreateToken(issuedAt, config.ConfigInstance.GetPasswordResetTokenDuration(), jwt.BackendIssuer, jwt.PasswordResetTokenType, &sessionId, nil)
}

func CreateMFAToken(issuedAt time.Time, sessionId string) (token string, JTI string, expiration time.Time, err error) {
	return CreateToken(issuedAt, config.ConfigInstance.GetMFATokenDuration(), jwt.BackendIssuer, jwt.MFATokenType, &sessionId, nil)
}

const backendTokenDuration = 100 * 365 * 24 * time.Hour // 100 years
// TODO: this needs to be changed to a more reasonable value once token rotation is implemented

func CreateBackendToken(issuedAt time.Time) (token string, JTI string, expiration time.Time, err error) {
	return CreateToken(issuedAt, backendTokenDuration, jwt.BackendIssuer, jwt.BackendTokenType, nil, nil)
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
