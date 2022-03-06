package client

import "io"

func WithExport(shouldExport bool) ConfigureFunc {
	return func(client *Client) {
		client.shouldExport = shouldExport
	}
}

func WithWriter(writer io.Writer) ConfigureFunc {
	return func(client *Client) {
		client.writer = writer
	}
}

func WithNamespace(namespace string) ConfigureFunc {
	return func(client *Client) {
		client.namespace = namespace
	}
}

func WithFilename(filename string) ConfigureFunc {
	return func(client *Client) {
		client.filename = filename
	}
}
