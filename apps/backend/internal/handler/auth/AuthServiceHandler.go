package auth

import "panelium/proto-gen-go/proto_gen_goconnect"

type AuthServiceHandler struct {
	proto_gen_goconnect.AuthServiceHandler
}

func noop(input ...any) {} // TODO: remove this, just for avoiding unused import error
