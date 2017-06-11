package mysql

import (
	"database/sql"

	storageI "github.com/imega-teleport/db2file/storage"
	"github.com/imega-teleport/xml2db/commerceml"
)

type storage struct {
	db *sql.DB
}

// NewStorage return storage instance
func NewStorage(sqlDB *sql.DB) storageI.Store {
	return &storage{
		db: sqlDB,
	}
}

func (s *storage) Groups(parentID string) (groups []commerceml.Group, err error) {
	groups = []commerceml.Group{}
	items, err := s.groups("")
	if err != nil {
		return nil, err
	}
	for _, i := range items {
		g := commerceml.Group{
			IdName: commerceml.IdName{
				Id:   i.ID,
				Name: i.Name,
			},
		}
		childs, err := s.childsGroup(g)
		if err != nil {
			return nil, err
		}
		g.Groups = childs
		groups = append(groups, g)
	}
	return
}

type group struct {
	ID       string
	ParentID string
	Name     string
}

func (s *storage) childsGroup(group commerceml.Group) (groups []commerceml.Group, err error) {
	groups = []commerceml.Group{}
	items, err := s.groups(group.Id)
	for _, i := range items {
		g := commerceml.Group{
			IdName: commerceml.IdName{
				Id:   i.ID,
				Name: i.Name,
			},
		}
		childs, err := s.childsGroup(g)
		if err != nil {
			return nil, err
		}
		g.Groups = childs
		groups = append(groups, g)
	}
	return
}

func (s *storage) groups(parentID string) (groups []group, err error) {
	groups = []group{}
	rows, err := s.db.Query("select id,parent_id,name from groups where parent_id = ?", parentID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		g := struct {
			ID       string
			ParentID string
			Name     string
		}{}
		if err = rows.Scan(&g.ID, &g.ParentID, &g.Name); err != nil {
			return nil, err
		}
		groups = append(groups, g)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return
}

func (s *storage) Products() (products []commerceml.Product, err error) {
	rows, err := s.db.Query("select id, name, description, barcode, article, full_name, country, brand from products")
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		item := struct {
			ID          string
			Name        string
			Description string
			Barcode     string
			Article     string
			FullName    string
			Country     string
			Brand       string
		}{}
		err = rows.Scan(&item.ID, &item.Name, &item.Description, &item.Barcode, &item.Article, &item.FullName, &item.Country, &item.Brand)
		if err != nil {
			return nil, err
		}
		product := commerceml.Product{
			IdName: commerceml.IdName{
				Id:   item.ID,
				Name: item.Name,
			},
			Description: commerceml.Description{
				Value: item.Description,
			},
			BarCode:  item.Barcode,
			Article:  item.Article,
			FullName: item.FullName,
			Country:  item.Country,
			Brand:    item.Brand,
			Groups:   s.productGroup(item.ID),
		}
		products = append(products, product)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return
}

func (s *storage) productGroup(parentID string) (groups []commerceml.Group, err error) {
	rows, err := s.db.Query("select id from products_groups where = ?", parentID)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		item := struct {
			ID string
		}{}
		err = rows.Scan(&item.ID)
		if err != nil {
			return nil, err
		}
		group := commerceml.Group{
			IdName: commerceml.IdName{
				Id: item.ID,
			},
		}
		groups = append(groups, group)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return
}

func (s *storage) ProductGroup() (products []commerceml.Product, err error) {
	rows, err := s.db.Query("select parent_id,id from products_groups")
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		item := struct {
			ParentID string
			ID       string
		}{}
		err = rows.Scan(&item.ParentID, &item.ID)
		if err != nil {
			return nil, err
		}
		product := commerceml.Product{
			IdName: commerceml.IdName{
				Id: item.ParentID,
			},
			Groups: commerceml.Group{
				IdName: commerceml.IdName{
					Id: item.ID,
				},
			},
		}
		products = append(products, product)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return
}
