package grpc

import (
	"context"
	"net"
	"time"

	"github.com/google/uuid"
	"github.com/horiondreher/go-parking-lot/parkings/internal/adapters/grpc/proto"
	"github.com/horiondreher/go-parking-lot/parkings/internal/adapters/queue"
	"github.com/horiondreher/go-parking-lot/parkings/internal/utils"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type GRPCServer struct {
	proto.UnimplementedReservationServiceServer
	queueAdapter *queue.QueueAdapter
}

func (s *GRPCServer) publishMessage(ctx context.Context) error {
	return nil
}

func (s *GRPCServer) GetReservation(ctx context.Context, req *proto.GetReservationRequest) (*proto.GetReservationResponse, error) {
	timeRemaining := time.Now().Add(2 * time.Hour)

	err := s.queueAdapter.PublishOnUserUpdated("new parking reservation created")
	if err != nil {
		log.Err(err).Msg("error publishing message")
	}

	return &proto.GetReservationResponse{
		Id:            uuid.New().String(),
		Type:          "car",
		RemainingTime: timestamppb.New(timeRemaining),
	}, nil
}

type GRPCAdapter struct {
	config *utils.Config
	server *grpc.Server
}

func NewAdapter(queueAdapter *queue.QueueAdapter) *GRPCAdapter {
	config := utils.GetConfig()

	gRPCServer := grpc.NewServer()
	proto.RegisterReservationServiceServer(gRPCServer, &GRPCServer{queueAdapter: queueAdapter})
	reflection.Register(gRPCServer)

	gRPCAdapter := &GRPCAdapter{
		config: config,
		server: gRPCServer,
	}

	return gRPCAdapter
}

func (adapter *GRPCAdapter) Start() error {
	listener, err := net.Listen("tcp", adapter.config.GRPCServerAddress)
	if err != nil {
		return err
	}

	log.Info().Str("address", adapter.config.GRPCServerAddress).Msg("starting gRPC server")
	return adapter.server.Serve(listener)
}
