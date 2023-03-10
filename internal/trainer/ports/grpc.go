package ports

import (
	"context"
	"errors"
	"time"

	"github.com/dbaeka/workouts-go/internal/trainer/app"
	"google.golang.org/grpc/codes"

	"github.com/dbaeka/workouts-go/internal/common/genproto/trainer"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/golang/protobuf/ptypes/timestamp"
	"google.golang.org/grpc/status"
)

type GrpcServer struct {
	app app.Application
}

func NewGrpcServer(app app.Application) GrpcServer {
	return GrpcServer{app: app}
}

func (g GrpcServer) MakeHourAvailable(ctx context.Context, req *trainer.UpdateHourRequest) (*empty.Empty, error) {
	trainingTime, err := protoTimestampToTime(req.Time)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "unable to parse time")
	}

	if err := g.app.Commands.MakeHoursAvailable.Handle(ctx, []time.Time{trainingTime}); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &empty.Empty{}, nil
}

func (g GrpcServer) ScheduleTraining(ctx context.Context, request *trainer.UpdateHourRequest) (*empty.Empty, error) {
	trainingTime, err := protoTimestampToTime(request.Time)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "unable to parse time")
	}

	if err := g.app.Commands.ScheduleTraining.Handle(ctx, trainingTime); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &empty.Empty{}, nil
}

func (g GrpcServer) CancelTraining(ctx context.Context, request *trainer.UpdateHourRequest) (*empty.Empty, error) {
	trainingTime, err := protoTimestampToTime(request.Time)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "unable to parse time")
	}

	if err := g.app.Commands.CancelTraining.Handle(ctx, trainingTime); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &empty.Empty{}, nil
}

func (g GrpcServer) IsHourAvailable(ctx context.Context, req *trainer.IsHourAvailableRequest) (*trainer.IsHourAvailableResponse, error) {
	trainingTime, err := protoTimestampToTime(req.Time)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "unable to parse time")
	}

	isAvailable, err := g.app.Queries.HourAvailability.Handle(ctx, trainingTime)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &trainer.IsHourAvailableResponse{IsAvailable: isAvailable}, nil
}

func protoTimestampToTime(timestamp *timestamp.Timestamp) (time.Time, error) {
	err := timestamp.CheckValid()
	if err != nil {
		return time.Time{}, errors.New("unable to parse time")
	}

	t := timestamp.AsTime()
	t = t.UTC().Truncate(time.Hour)

	return t, nil
}
