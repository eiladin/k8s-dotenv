package client

import (
	"errors"
	"fmt"
)

// ErrMissingResource is returned when the resource is not found.
var ErrMissingResource = errors.New("resource not found")

// ErrAPIGroup is returned when a kubernetes api call fails.
var ErrAPIGroup = errors.New("api group error")

func newMissingKubeClientError(client string) error {
	//nolint
	return fmt.Errorf("could not create %s client, missing call to WithKubeClient?", client)
}
