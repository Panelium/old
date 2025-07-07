package client

import (
	"connectrpc.com/connect"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"panelium/backend/internal/db"
	"panelium/backend/internal/global"
	"panelium/backend/internal/middleware"
	"panelium/backend/internal/model"
	"panelium/common/errors"
	"panelium/common/id"
	"panelium/proto_gen_go"
	"panelium/proto_gen_go/backend"
	"panelium/proto_gen_go/backend/backendconnect"
	"panelium/proto_gen_go/daemon"
	"panelium/proto_gen_go/daemon/daemonconnect"
)

type ClientServiceHandler struct {
	backendconnect.ClientServiceHandler
}

func (s *ClientServiceHandler) GetAvailableBlueprints(
	ctx context.Context,
	req *connect.Request[proto_gen_go.Empty],
) (*connect.Response[backend.AvailableBlueprints], error) {
	var blueprints []model.Blueprint
	tx := db.Instance().Find(&blueprints)
	if tx.Error != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to fetch available blueprints"))
	}

	var blueprintsProto []*backend.AvailableBlueprint
	for _, blueprint := range blueprints {
		blueprintProto := &backend.AvailableBlueprint{
			Bid:  blueprint.BID,
			Name: blueprint.Name,
		}
		blueprintsProto = append(blueprintsProto, blueprintProto)
	}

	availableBlueprints := &backend.AvailableBlueprints{
		Blueprints: blueprintsProto,
	}

	return connect.NewResponse(availableBlueprints), nil
}

func (s *ClientServiceHandler) GetAvailableLocations(
	ctx context.Context,
	req *connect.Request[proto_gen_go.Empty],
) (*connect.Response[backend.AvailableLocations], error) {
	var locations []model.Location
	tx := db.Instance().Find(&locations)
	if tx.Error != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to fetch available locations"))
	}

	var locationsProto []*backend.AvailableLocation
	for _, location := range locations {
		locationProto := &backend.AvailableLocation{
			Lid:  location.LID,
			Name: location.Name,
		}
		locationsProto = append(locationsProto, locationProto)
	}

	availableLocations := &backend.AvailableLocations{
		Locations: locationsProto,
	}

	return connect.NewResponse(availableLocations), nil
}

func (s *ClientServiceHandler) GetAvailableNodes(
	ctx context.Context,
	req *connect.Request[proto_gen_go.Empty],
) (*connect.Response[backend.AvailableNodes], error) {
	var nodes []model.Node
	tx := db.Instance().Find(&nodes)
	if tx.Error != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to fetch available nodes"))
	}

	var nodesProto []*backend.AvailableNode
	for _, node := range nodes {
		nodeProto := &backend.AvailableNode{
			Nid:  node.NID,
			Lid:  node.Location.LID,
			Name: node.Name,
		}
		nodesProto = append(nodesProto, nodeProto)
	}

	availableNodes := &backend.AvailableNodes{
		Nodes: nodesProto,
	}

	return connect.NewResponse(availableNodes), nil
}

func (s *ClientServiceHandler) NewServer(
	ctx context.Context,
	req *connect.Request[backend.NewServerRequest],
) (*connect.Response[backend.NewServerResponse], error) {
	sessionInfoData := ctx.Value("panelium_session_info")
	sessionInfo, ok := sessionInfoData.(*middleware.SessionInfo)
	if !ok || sessionInfo == nil || sessionInfo.SessionID == "" || sessionInfo.UserID == "" {
		return nil, errors.ConnectInvalidCredentials
	}

	dbInstance := db.Instance()
	tx := dbInstance.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var user model.User
	if err := tx.First(&user, "uid = ?", sessionInfo.UserID).Error; err != nil {
		tx.Rollback()
		return nil, connect.NewError(connect.CodeNotFound, errors.UserNotFound)
	}

	var node model.Node

	if req.Msg.Nid != nil && *req.Msg.Nid != "" {
		if err := tx.First(&node, "nid = ?", *req.Msg.Nid).Error; err != nil {
			tx.Rollback()
			return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("node with ID %s not found", *req.Msg.Nid))
		}
	} else if req.Msg.Lid != nil && *req.Msg.Lid != "" {
		var location model.Location
		if err := tx.First(&location, "lid = ?", *req.Msg.Lid).Error; err != nil {
			tx.Rollback()
			return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("location with ID %s not found", *req.Msg.Lid))
		}

		var nodes []model.Node
		if err := tx.Find(&nodes, "location_id = ?", location.ID).Error; err != nil || len(nodes) == 0 {
			tx.Rollback()
			return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("no nodes found in location with ID %s", *req.Msg.Lid))
		}
		node = nodes[0]
	} else {
		if err := tx.First(&node).Error; err != nil {
			tx.Rollback()
			return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("no nodes available"))
		}
	}

	var blueprint model.Blueprint
	if err := tx.First(&blueprint, "bid = ?", req.Msg.Bid).Error; err != nil {
		tx.Rollback()
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("blueprint with ID %s not found", req.Msg.Bid))
	}

	var availableDockerImages []struct {
		Name  string `json:"name"`
		Image string `json:"image"`
	}
	if err := json.Unmarshal(blueprint.DockerImages, &availableDockerImages); err != nil {
		tx.Rollback()
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to parse docker images for blueprint"))
	}

	if len(availableDockerImages) == 0 {
		tx.Rollback()
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("no docker images available for blueprint %s", req.Msg.Bid))
	}

	dockerImage := availableDockerImages[0].Image // For now, just use the first available image

	resourceLimit := model.ResourceLimit{
		CPU:     50,   // 50% CPU
		RAM:     1024, // 1 GB
		SWAP:    0,
		Storage: 4096, // 4 GB
	}

	sid, err := id.New()
	if err != nil {
		tx.Rollback()
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to create server (ID)"))
	}

	server := &model.Server{
		SID:           sid,
		Name:          req.Msg.Name,
		Description:   req.Msg.Description,
		OwnerID:       user.ID,
		NodeID:        node.ID,
		ResourceLimit: resourceLimit,
		DockerImage:   dockerImage,
		BID:           req.Msg.Bid,
	}
	if err := tx.Create(server).Error; err != nil {
		tx.Rollback()
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to create server"))
	}

	var allocation model.NodeAllocation
	if err := tx.Where("node_id = ? AND server_id IS NULL", node.ID).First(&allocation).Error; err != nil {
		tx.Rollback()
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to find available node allocation for node %s", node.NID))
	}

	allocation.ServerID = &server.ID
	if err := tx.Save(&allocation).Error; err != nil {
		tx.Rollback()
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to update node allocation with server ID"))
	}

	// Use dbInstance for read-only operations after this point if needed

	daemonClient := daemonconnect.NewBackendServiceClient(http.DefaultClient, fmt.Sprintf("https://%s:%d", node.FQDN, node.DaemonPort), connect.WithGRPC())
	if daemonClient == nil {
		tx.Rollback()
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to create daemon client"))
	}

	createServerReq := connect.NewRequest(&daemon.Server{
		Sid:     server.SID,
		OwnerId: user.UID,
		Allocations: []*proto_gen_go.IPAllocation{
			{
				Ip:   allocation.IP,
				Port: uint32(allocation.Port),
			},
		},
		ResourceLimit: &proto_gen_go.ResourceLimit{
			Cpu:     uint32(resourceLimit.CPU),
			Ram:     uint32(resourceLimit.RAM),
			Swap:    uint32(resourceLimit.SWAP),
			Storage: uint32(resourceLimit.Storage),
		},
		DockerImage: dockerImage,
		Bid:         req.Msg.Bid,
	})

	if node.EncryptedNodeTokenBase64 == nil || *node.EncryptedNodeTokenBase64 == "" {
		tx.Rollback()
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("node %s not properly set up", node.NID))
	}

	encryption := *global.EncryptionInstance()
	encryptedNodeTokenBytes, err := base64.StdEncoding.DecodeString(*node.EncryptedNodeTokenBase64)
	if err != nil {
		tx.Rollback()
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to decode encrypted node token"))
	}
	decryptedNodeTokenBytes := make([]byte, len(encryptedNodeTokenBytes))
	encryption.Decrypt(decryptedNodeTokenBytes, encryptedNodeTokenBytes)
	decryptedNodeToken := string(decryptedNodeTokenBytes)

	createServerReq.Header().Add("Authorization", decryptedNodeToken)

	createServerRes, err := daemonClient.CreateServer(context.Background(), createServerReq)
	if err != nil {
		tx.Rollback()
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to create server on daemon"))
	}

	if createServerRes.Msg == nil || !createServerRes.Msg.Success {
		tx.Rollback()
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to create server on daemon"))
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to commit transaction"))
	}

	res := &backend.NewServerResponse{
		Sid: server.SID,
	}

	return connect.NewResponse(res), nil
}
