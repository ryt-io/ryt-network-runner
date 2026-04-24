package api

import (
	"context"

	"github.com/ryt-io/ryt-v2/api/health"
	"github.com/ryt-io/ryt-v2/utils/rpc"
)

// HealthClient defines the interface for health client operations
type HealthClient interface {
	Health(ctx context.Context, tags []string, options ...rpc.Option) (*health.APIReply, error)
	Readiness(ctx context.Context, tags []string, options ...rpc.Option) (*health.APIReply, error)
	Liveness(ctx context.Context, tags []string, options ...rpc.Option) (*health.APIReply, error)
}
