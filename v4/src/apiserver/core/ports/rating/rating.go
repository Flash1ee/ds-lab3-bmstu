package rating

import "context"

type Config struct {
	Host            string
	Port            string
	MaxErrorsTrying int64 `mapstructure:"max_errors_trying"`
}

type Client interface {
	GetUserRating(ctx context.Context, username string) (Rating, error)
	UpdateUserRating(ctx context.Context, username string, diff int) error
}
