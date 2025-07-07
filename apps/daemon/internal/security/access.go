package security

import (
	"context"
)

// TODO

func HasAccessToServer(userId string, serverId string) bool {
	//tx := db.Instance().First(&model.ServerUser{}, "user_id = ? AND sid = ?", userId, serverId)
	//if tx.Error != nil || tx.RowsAffected == 0 {
	//	tx = db.Instance().First(&model.Server{}, "sid = ? AND owner_id = ?", serverId, userId)
	//	if tx.Error != nil || tx.RowsAffected == 0 {
	//		return false
	//	}
	//}
	return true
}

func CheckServerAccess(ctx context.Context, serverId string) error {
	//userIdAny := ctx.Value("user_id")
	//userId, ok := userIdAny.(*string)
	//if !ok || userId == nil || *userId == "" {
	//	return errors.New("user not authenticated")
	//}
	//
	//if !HasAccessToServer(*userId, serverId) {
	//	return errors.New("user does not have access to this server")
	//}

	return nil
}
