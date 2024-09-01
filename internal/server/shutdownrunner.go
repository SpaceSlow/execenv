package server

import "context"

type ShutdownRunner interface {
	Run() error
	Shutdown(ctx context.Context) error
}
