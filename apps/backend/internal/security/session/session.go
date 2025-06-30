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

func CreateSession(uid string) (sessionId string, refreshToken string, accessToken string, refreshTokenExpiration time.Time, accessTokenExpiration time.Time, err error) {
	sessionId, err = id.New()
	if err != nil {
		return "", "", "", time.Time{}, time.Time{}, stdErrors.Join(err, errors.SessionCreationFailed)
	}

	timeNow := time.Now()

	refreshToken, refreshJTI, refreshTokenExpiration, err := CreateRefreshToken(timeNow, sessionId, uid)
	if err != nil {
		return "", "", "", time.Time{}, time.Time{}, stdErrors.Join(err, errors.SessionCreationFailed)
	}
	accessToken, accessJTI, accessTokenExpiration, err := CreateAccessToken(timeNow, sessionId, uid)
	if err != nil {
		return "", "", "", time.Time{}, time.Time{}, stdErrors.Join(err, errors.SessionCreationFailed)
	}

	tx := db.Instance().Model(model.UserSession{}).Create(&model.UserSession{
		SessionID:  sessionId,
		UserID:     uid,
		AccessJTI:  accessJTI,
		RefreshJTI: refreshJTI,
		Expiration: refreshTokenExpiration,
	})

	if tx.Error != nil || tx.RowsAffected == 0 {
		return "", "", "", time.Time{}, time.Time{}, stdErrors.Join(tx.Error, errors.SessionCreationFailed)
	}

	return sessionId, refreshToken, accessToken, refreshTokenExpiration, accessTokenExpiration, nil
}

func RefreshSession(sessionId string) (refreshToken string, accessToken string, refreshTokenExpiration time.Time, accessTokenExpiration time.Time, err error) {
	session := &model.UserSession{}
	tx := db.Instance().Model(model.UserSession{}).First(session, "session_id = ?", sessionId)
	if tx.Error != nil || tx.RowsAffected == 0 {
		return "", "", time.Time{}, time.Time{}, stdErrors.Join(tx.Error, errors.SessionNotFound)
	}

	uid := session.UserID

	timeNow := time.Now()

	refreshToken, refreshJTI, refreshTokenExpiration, err := CreateRefreshToken(timeNow, sessionId, uid)
	if err != nil {
		return "", "", time.Time{}, time.Time{}, stdErrors.Join(err, errors.SessionCreationFailed)
	}
	accessToken, accessJTI, accessTokenExpiration, err := CreateAccessToken(timeNow, sessionId, uid)
	if err != nil {
		return "", "", time.Time{}, time.Time{}, stdErrors.Join(err, errors.SessionCreationFailed)
	}

	tx = db.Instance().Model(model.UserSession{}).Where("session_id = ?", sessionId).Updates(&model.UserSession{
		AccessJTI:  accessJTI,
		RefreshJTI: refreshJTI,
		Expiration: refreshTokenExpiration,
	})
	if tx.Error != nil || tx.RowsAffected == 0 {
		return "", "", time.Time{}, time.Time{}, stdErrors.Join(tx.Error, errors.SessionCreationFailed)
	}

	return refreshToken, accessToken, refreshTokenExpiration, accessTokenExpiration, nil
}

func DeleteSession(sessionId string) error {
	tx := db.Instance().Model(model.UserSession{}).Where("session_id = ?", sessionId).Delete(&model.UserSession{})
	if tx.Error != nil || tx.RowsAffected == 0 {
		return stdErrors.Join(tx.Error, errors.SessionNotFound)
	}

	return nil
}
