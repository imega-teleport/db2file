package storage

import "database/sql"

type ProductsPrices struct {
	ProductID string
	PriceType string
	Value     string
	Currency  string
	Unit      string
}

func (s *storage) GetProductsPrices(out chan<- interface{}, e chan<- error) {
	s.getRecords(out, e, "select parent_id, price_type_id, unit_price, currency, unit from bundling_offers_prices", func(rows *sql.Rows) (interface{}, error) {
		i := ProductsPrices{}
		err := rows.Scan(&i.ProductID, &i.PriceType, &i.Value, &i.Currency, &i.Unit)
		return i, err
	})
}
