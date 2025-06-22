package session

import (
	"crypto/rand"
	goErrors "github.com/pkg/errors"
	"panelium/backend/internal/errors"
	"panelium/backend/internal/global"
	"panelium/backend/internal/model"
	"panelium/common/jwt"
	"time"
)

func CreateSession(uid string) (sessionId string, refreshToken string, err error) {
	sessionId = rand.Text() // TODO: do the math to ensure this is sufficient enough not to run into collisions
	jti := rand.Text()      // TODO: generate this differently probably- maybe use a UUID?

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
		return "", "", goErrors.Wrap(err, "failed to create JWT for session") // TODO: move error message to errors package
	}

	result := global.DB.Model(model.UserSession{}).Create(&model.UserSession{
		SessionID:  sessionId,
		UserID:     uid,
		RefreshJTI: jti,
	})

	if result.Error != nil {
		return "", "", goErrors.Wrap(result.Error, "failed to create user session") // TODO: move error message to errors package
	}
	if result.RowsAffected == 0 {
		return "", "", errors.SessionCreationFailed
	}

	return sessionId, refreshToken, nil
}
