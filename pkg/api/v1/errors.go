package v1

import "fmt"

// NewResourceLoadError wraps API errors when a resource is not found.
func NewResourceLoadError(err error) error {
	return fmt.Errorf("error loading resource: %w", err)
}
