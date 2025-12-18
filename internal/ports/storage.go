package ports

type Storage interface {
	Set(key string, value string) error
	Get(key string) (value string, ok bool, err error)
}
