package client

import (
	"fmt"

	uuid "github.com/satori/go.uuid"
)

type (
	ID = uuid.UUID
)

var (
	ErrClientWasBlocked         = fmt.Errorf("external server was blocked")
	ErrExternalServiceZeroLimit = fmt.Errorf("external service return zero limit")
	ErrExternalServiceZeroDelay = fmt.Errorf("external service return zero delay")
)

type Item struct {
	ID ID
	// other data
}
