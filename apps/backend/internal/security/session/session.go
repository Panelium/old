package session

import (
	stdErrors "errors"
	"panelium/backend/internal/errors"
	"panelium/backend/internal/global"
	"panelium/backend/internal/model"
	"panelium/common/id"
	"panelium/common/jwt"
	"time"
)

func CreateSession(uid string) (sessionId string, refreshToken string, err error) {
	sessionId, err = id.New()
	if err != nil {
		return "", "", stdErrors.Join(err, errors.SessionCreationFailed)
	}
	jti, err := id.New()
	if err != nil {
		return "", "", stdErrors.Join(err, errors.SessionCreationFailed)
	}

	claims := jwt.Claims{
		IssuedAt:   time.Now().Unix(),
		NotBefore:  nil,
		Expiration: time.Now().Add(time.Hour * 24).Unix(), // TODO: this needs more thought, perhaps config?
		Subject:    &uid,
		Audience:   sessionId,
		Issuer:     "backend", // TODO: we might want to make this shorter
		TokenType:  "refresh", // TODO: we might want to make this shorter
		JTI:        &jti,
	}

	refreshToken, err = jwt.CreateJWT(claims, global.JWTSecret)
	if err != nil {
		return "", "", stdErrors.Join(err, errors.SessionCreationFailed)
	}

	result := global.DB.Model(model.UserSession{}).Create(&model.UserSession{
		SessionID:  sessionId,
		UserID:     uid,
		RefreshJTI: jti,
	})

	if result.Error != nil {
		return "", "", stdErrors.Join(result.Error, errors.SessionCreationFailed)
	}
	if result.RowsAffected == 0 {
		return "", "", errors.SessionCreationFailed
	}

	return sessionId, refreshToken, nil
}
