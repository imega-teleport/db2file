package exporter

import "github.com/imega-teleport/db2file/storage"

type woocommece struct {
	storage storage.Store
}

func NewExporter(storage storage.Store) *woocommece {
	return &woocommece{
		storage: storage,
	}
}
