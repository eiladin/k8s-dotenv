package client

import (
	"errors"
	"fmt"
)

// ErrResourceNameRequired is returned when no resource name is provided.
var ErrResourceNameRequired = errors.New("resource name required")

// ErrMissingResource is returned when the resource is not found.
var ErrMissingResource = errors.New("resource not found")

// ErrCreatingClient is returned when the client cannot be created.
var ErrCreatingClient = errors.New("client create error")

// ErrAPIGroup is returned when a kubernetes api call fails.
var ErrAPIGroup = errors.New("api group error")

// ErrNoFilename is returned when no filename is provided.
var ErrNoFilename = errors.New("no filename provided")

func newSecretErr(err error) error {
	return fmt.Errorf("secret error: %w", err)
}

func newConfigMapError(err error) error {
	return fmt.Errorf("configmap error: %w", err)
}

func newWriteError(err error) error {
	return fmt.Errorf("write error: %w", err)
}

func newMissingKubeClientError(client string) error {
	return fmt.Errorf("could not create %s client, missing call to WithKubeClient?", client)
}

// ResourceLoadError wraps API errors when a resource is not found.
type ResourceLoadError struct {
	Err      error
	Resource string
}

// NewResourceLoadError creates a `ResourceLoadError`.
func NewResourceLoadError(resource string, err error) error {
	return &ResourceLoadError{
		Err:      err,
		Resource: resource,
	}
}

// Error returns the message on the internal error (if there is one).
func (e *ResourceLoadError) Error() string {
	if e.Err != nil {
		return fmt.Errorf("error loading %s: %w", e.Resource, e.Err).Error()
	}

	return fmt.Sprintf("error loading %s", e.Resource)
}

// Unwrap returns the internal error.
func (e *ResourceLoadError) Unwrap() error {
	return e.Err
}
