package client

import "io"

// WithExport flags the client to include `export` statements in the output.
func WithExport(shouldExport bool) ConfigureFunc {
	return func(client *Client) {
		client.shouldExport = shouldExport
	}
}

// WithWriter sets the `io.Writer` to use for output.
func WithWriter(writer io.Writer) ConfigureFunc {
	return func(client *Client) {
		client.writer = writer
	}
}

// WithNamespace sets the namespace to use when interacting with the Kubernetes API.
func WithNamespace(namespace string) ConfigureFunc {
	return func(client *Client) {
		client.namespace = namespace
	}
}

// WithFilename sets the name of the file to output into.
func WithFilename(filename string) ConfigureFunc {
	return func(client *Client) {
		client.filename = filename
	}
}
