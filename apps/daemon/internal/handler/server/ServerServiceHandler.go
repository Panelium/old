package server

import "panelium/proto_gen_go/daemon/daemonconnect"

type ServerServiceHandler struct {
	daemonconnect.ServerServiceHandler
}

// ctx.Value("server_id")
