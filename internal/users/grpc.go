package main

import (
	"context"
	"fmt"

	"github.com/golang/protobuf/ptypes/empty"

	"github.com/dbaeka/workouts-go/internal/common/genproto/users"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GrpcServer defines interface for the Porobuf GRPC generated code. Must also embed the Unimplemented type as well as
// the RPC functions
type GrpcServer struct {
	db db
}

func (g GrpcServer) GetTrainingBalance(
	ctx context.Context,
	r *users.GetTrainingBalanceRequest,
) (*users.GetTrainingBalanceResponse, error) {
	user, err := g.db.GetUser(ctx, r.UserId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &users.GetTrainingBalanceResponse{Amount: int64(user.Balance)}, nil
}

func (g GrpcServer) UpdateTrainingBalance(
	ctx context.Context,
	r *users.UpdateTrainingBalanceRequest,
) (*empty.Empty, error) {
	err := g.db.UpdateBalance(ctx, r.UserId, int(r.AmountChange))
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to update balance: %s", err))
	}

	return &empty.Empty{}, nil
}
