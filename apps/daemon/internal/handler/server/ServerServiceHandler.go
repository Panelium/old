package server

import "panelium/proto_gen_go/daemon/daemonconnect"

type ServerServiceHandler struct {
	daemonconnect.ServerServiceHandler
}

// TODO: check if user has access to the server and if server id is even valid
