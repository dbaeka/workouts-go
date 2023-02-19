package client

import (
	"crypto/tls"
	"crypto/x509"
	"os"
	"strconv"
	"time"

	"github.com/dbaeka/workouts-go/internal/common/genproto/trainer"
	"github.com/dbaeka/workouts-go/internal/common/genproto/users"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

func NewTrainerClient() (client trainer.TrainerServiceClient, close func() error, err error) {
	grpcAddr := os.Getenv("TRAINER_GRPC_ADDR")

	client = nil
	close = func() error { return nil }

	if grpcAddr == "" {
		err = errors.New("empty env TRAINER_GRPC_ADDR")
		return
	}

	opts, err := grpcDialOpts(grpcAddr)
	if err != nil {
		return
	}

	conn, err := grpc.Dial(grpcAddr, opts...)
	if err != nil {
		return
	}

	client = trainer.NewTrainerServiceClient(conn)
	close = conn.Close
	return
}

func WaitForTrainerService(timeout time.Duration) bool {
	return waitForPort(os.Getenv("TRAINER_GRPC_ADDR"), timeout)
}

func NewUsersClient() (client users.UsersServiceClient, close func() error, err error) {
	grpcAddr := os.Getenv("USERS_GRPC_ADDR")

	client = nil
	close = func() error { return nil }

	if grpcAddr == "" {
		err = errors.New("empty env USERS_GRPC_ADDR")
		return
	}

	opts, err := grpcDialOpts(grpcAddr)
	if err != nil {
		return
	}

	conn, err := grpc.Dial(grpcAddr, opts...)
	if err != nil {
		return
	}

	client = users.NewUsersServiceClient(conn)
	close = conn.Close
	return
}

func WaitForUsersService(timeout time.Duration) bool {
	return waitForPort(os.Getenv("USERS_GRPC_ADDR"), timeout)
}

func grpcDialOpts(grpcAddr string) ([]grpc.DialOption, error) {
	if noTLS, _ := strconv.ParseBool(os.Getenv("GRPC_NO_TLS")); noTLS {
		return []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}, nil
	}

	systemRoots, err := x509.SystemCertPool()
	if err != nil {
		return nil, errors.Wrap(err, "cannot load root CA cert")
	}
	creds := credentials.NewTLS(&tls.Config{
		RootCAs:    systemRoots,
		MinVersion: tls.VersionTLS12,
	})

	return []grpc.DialOption{
		grpc.WithTransportCredentials(creds),
		grpc.WithPerRPCCredentials(newMetadataServerToken(grpcAddr)),
	}, nil
}
