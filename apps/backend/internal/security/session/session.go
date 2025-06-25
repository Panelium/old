package session

import (
	stdErrors "errors"
	"panelium/backend/internal/db"
	"panelium/backend/internal/model"
	"panelium/common/errors"
	"panelium/common/id"
	"time"
)

// TODO: refactor this file

func CreateSession(uid string) (sessionId string, refreshToken string, accessToken string, err error) {
	sessionId, err = id.New()
	if err != nil {
		return "", "", "", stdErrors.Join(err, errors.SessionCreationFailed)
	}

	timeNow := time.Now()

	refreshToken, refreshJTI, refreshExpiration, err := CreateRefreshToken(timeNow, sessionId, uid)
	if err != nil {
		return "", "", "", stdErrors.Join(err, errors.SessionCreationFailed)
	}
	accessToken, accessJTI, _, err := CreateAccessToken(timeNow, sessionId, uid)
	if err != nil {
		return "", "", "", stdErrors.Join(err, errors.SessionCreationFailed)
	}

	result := db.Instance().Model(model.UserSession{}).Create(&model.UserSession{
		SessionID:  sessionId,
		UserID:     uid,
		AccessJTI:  accessJTI,
		RefreshJTI: refreshJTI,
		Expiration: refreshExpiration,
	})

	if result.Error != nil {
		return "", "", "", stdErrors.Join(result.Error, errors.SessionCreationFailed)
	}
	if result.RowsAffected == 0 {
		return "", "", "", errors.SessionCreationFailed
	}

	return sessionId, refreshToken, accessToken, nil
}

func RefreshSession(sessionId string) (refreshToken string, accessToken string, err error) {
	// TODO: need to review this, might have to change the logic
	result := db.Instance().Model(model.UserSession{}).Where("session_id = ?", sessionId).First(&model.UserSession{})
	if result.Error != nil {
		return "", "", stdErrors.Join(result.Error, errors.SessionNotFound)
	}
	if result.RowsAffected == 0 {
		return "", "", errors.SessionNotFound
	}

	session := &model.UserSession{}
	if err := result.Scan(session); err.Error != nil {
		return "", "", stdErrors.Join(err.Error, errors.SessionNotFound)
	}
	uid := session.UserID

	timeNow := time.Now()

	refreshToken, refreshJTI, refreshExpiration, err := CreateRefreshToken(timeNow, sessionId, uid)
	if err != nil {
		return "", "", stdErrors.Join(err, errors.SessionCreationFailed)
	}
	accessToken, accessJTI, _, err := CreateAccessToken(timeNow, sessionId, uid)
	if err != nil {
		return "", "", stdErrors.Join(err, errors.SessionCreationFailed)
	}

	result = db.Instance().Model(model.UserSession{}).Where("session_id = ?", sessionId).Updates(&model.UserSession{
		AccessJTI:  accessJTI,
		RefreshJTI: refreshJTI,
		Expiration: refreshExpiration,
	})
	if result.Error != nil {
		return "", "", stdErrors.Join(result.Error, errors.SessionCreationFailed)
	}
	if result.RowsAffected == 0 {
		return "", "", errors.SessionNotFound
	}

	return refreshToken, accessToken, nil
}

func DeleteSession(sessionId string) error {
	result := db.Instance().Model(model.UserSession{}).Where("session_id = ?", sessionId).Delete(&model.UserSession{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.SessionNotFound
	}

	return nil
}
