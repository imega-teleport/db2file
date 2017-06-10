package paginator

import "math"

// Paginator is the interface that wraps the basic Processing method
type Paginator interface {
	Processing(items []interface{}, forItem func(interface{}) interface{}, forPage func(interface{}) interface{})
}

type paginator struct {
	int
}

// NewPaginator returns instance paginator
func NewPaginator(perPage int) Paginator {
	return &paginator{perPage}
}

func (p *paginator) Processing(items []interface{}, forItem func(interface{}) interface{}, forPage func(interface{}) interface{}) {
	piece := p.int
	pages := math.Ceil(float64(len(items)) / float64(piece))
	for n := 0; n < int(pages); n++ {
		var page []interface{}
		start, end := 0+n*piece, n*piece
		end = end + piece
		if end > len(items) {
			end = len(items)
		}
		items = items[start:end]
		for _, i := range items {
			page = append(page, forItem(i))
		}
		forPage(page)
	}
}
