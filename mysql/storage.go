package mysql

import (
	"database/sql"
	"github.com/imega-teleport/xml2db/commerceml"
)

type storage struct {
	db *sql.DB
}

func NewStorage(sqlDB *sql.DB) *storage {
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
	return groups, err
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
	return groups, err
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
		err := rows.Scan(&g.ID, &g.ParentID, &g.Name)
		if err != nil {
			return nil, err
		}
		groups = append(groups, g)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return groups, nil
}
