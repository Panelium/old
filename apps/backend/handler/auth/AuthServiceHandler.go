package auth

import "panelium/proto-gen-go/proto_gen_goconnect"

type AuthServiceHandler struct {
	proto_gen_goconnect.AuthServiceHandler
}

func NewAuthServiceHandler() *AuthServiceHandler {
	return &AuthServiceHandler{}
}
