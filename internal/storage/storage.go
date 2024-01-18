package storage

type Storage interface {
	GetCounter(key string) (int, error)
	IncrementCounter(key string, ttl int) error
	RegisterBlock(key string, cooldown int) error
	IsBlocked(key string) (bool, error)
}
