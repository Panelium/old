package admin

import (
	"panelium/backend/internal/model"
	"panelium/proto_gen_go/backend/admin"
)

func UserModelToProto(u *model.User) *admin.User {
	if u == nil {
		return nil
	}
	return &admin.User{
		Uid:       u.UID,
		Username:  u.Username,
		Email:     u.Email,
		Admin:     u.Admin,
		MfaNeeded: u.MFANeeded,
	}
}

func UserProtoToModel(u *admin.User) *model.User {
	if u == nil {
		return nil
	}
	return &model.User{
		UID:       u.Uid,
		Username:  u.Username,
		Email:     u.Email,
		Admin:     u.Admin,
		MFANeeded: u.MfaNeeded,
	}
}
