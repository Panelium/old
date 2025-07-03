package server

import (
	"google.golang.org/protobuf/types/known/timestamppb"
	"panelium/daemon/internal/model"
	"panelium/proto_gen_go"
)

func Status(s *model.Server) *proto_gen_go.ServerStatus {
	status := proto_gen_go.ServerStatus{
		Status:         s.Status,
		TimestampStart: timestamppb.New(s.TimestampStart),
		TimestampEnd:   timestamppb.New(s.TimestampEnd),
		OfflineReason:  &s.OfflineReason,
	}

	return &status
}
