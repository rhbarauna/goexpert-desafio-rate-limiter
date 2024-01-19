package storage

type Storage interface {
	GetCounter(key string) (int, error)
	IncrementCounter(key string, ttl int) (int64, error)
	RegisterBlock(key string, cooldown int) error
	IsBlocked(key string) (bool, error)

	// Increment(ctx context.Context, key string) error
	// Get(ctx context.Context, key string) (int, error)
	// Expire(ctx context.Context, key string, duration time.Duration) error
}
