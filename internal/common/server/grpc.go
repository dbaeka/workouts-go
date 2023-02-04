package server

import (
	"fmt"
	"net"
	"os"

	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpcctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func RunGRPCServer(registerServer func(server *grpc.Server)) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	grpcEndpoint := fmt.Sprintf(":%s", port)

	logrusEntry := logrus.NewEntry(logrus.StandardLogger())
	grpc_logrus.ReplaceGrpcLogger(logrusEntry)

	grpcServer := grpc.NewServer(
		grpcmiddleware.WithUnaryServerChain(
			grpcctxtags.UnaryServerInterceptor(grpcctxtags.WithFieldExtractor(grpcctxtags.CodeGenRequestFieldExtractor)),
			grpc_logrus.UnaryServerInterceptor(logrusEntry),
		),
		grpcmiddleware.WithStreamServerChain(
			grpcctxtags.StreamServerInterceptor(grpcctxtags.WithFieldExtractor(grpcctxtags.CodeGenRequestFieldExtractor)),
			grpc_logrus.StreamServerInterceptor(logrusEntry),
		),
	)
	registerServer(grpcServer)
	listen, err := net.Listen("tcp", grpcEndpoint)
	if err != nil {
		logrus.Fatal(err)
	}
	logrus.WithField("grpc_endpoint", grpcEndpoint).Info("Starting: gRPC Listener")
	logrus.Fatal(grpcServer.Serve(listen))
}
