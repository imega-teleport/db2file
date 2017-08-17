package storage

import "github.com/imega-teleport/xml2db/commerceml"

// Store is the interface the basic Storage
type Store interface {
	CheckCompleteAllTasks() (bool, error)
	Groups(parentID string) (groups []commerceml.Group, err error)
	Products() (products []commerceml.Product, err error)
}
