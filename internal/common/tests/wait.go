package tests

import (
	"time"

	net "github.com/dbaeka/workouts-go/internal/common/client"
)

func WaitForPort(address string) bool {
	return net.WaitForPort(address, 5*time.Second)
}
