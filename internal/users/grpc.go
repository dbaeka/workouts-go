package main

import (
	"context"
	"fmt"
	"github.com/dbaeka/workouts-go/internal/common/genproto/users"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GrpcServer struct {
	db db
	users.UnimplementedUsersServiceServer
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
) (*users.EmptyResponse, error) {
	err := g.db.UpdateBalance(ctx, r.UserId, int(r.AmountChange))
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to update balance: %s", err))
	}

	return &users.EmptyResponse{}, nil
}
