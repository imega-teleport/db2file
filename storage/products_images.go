package storage

import "database/sql"

// ProductImage запись картинки продукта
type ProductImage struct {
	ProductID string
	EntityID  string
	URL       string
}

func (s *storage) GetProductsImages(out chan<- interface{}, e chan<- error) {
	s.getRecords(out, e, "select parent_id, entity_id, url from products_images", func(rows *sql.Rows) (interface{}, error) {
		i := ProductImage{}
		err := rows.Scan(&i.ProductID, &i.EntityID, &i.URL)
		return i, err
	})
}
