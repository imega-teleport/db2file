package storage

import "github.com/imega-teleport/xml2db/commerceml"

type Store interface {
	Groups(parentID string) (groups []commerceml.Group, err error)
}
