package storage

type Store interface {
	Groups() ([]Group, err error)
}
