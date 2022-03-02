package cmd

import "errors"

// ErrResourceNameRequired is raised when no resource name is provided.
var ErrResourceNameRequired = errors.New("resource name required")
