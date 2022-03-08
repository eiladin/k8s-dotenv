package v1beta1

import (
	"fmt"
)

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
