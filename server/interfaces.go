package server

import "context"

type Server interface {
	ListenAndServe() error
	Shutdown(ctx context.Context) error
}
