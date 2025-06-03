package cache

type Cache interface {
	Get(key string) (any, error)
	Set(key string, value any, expire int64)
}
