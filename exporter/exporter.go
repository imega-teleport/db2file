package exporter

import (
	"fmt"

	"github.com/gosimple/slug"
	"github.com/imega-teleport/db2file/storage"
	"github.com/imega-teleport/xml2db/commerceml"
	"gopkg.in/Masterminds/squirrel.v1"
)

type woocommece struct {
	storage storage.Store
}

// NewExporter return new exporter instance
func NewExporter(storage storage.Store) *woocommece {
	return &woocommece{
		storage: storage,
	}
}

type term struct {
	ID    ID
	Name  string
	Slug  Slug
	Group ID
}

type Slug string

func (s Slug) String() string {
	return slug.Make(string(s))
}

type ID int

func (i ID) String() string {
	return fmt.Sprintf("@max_term_id+%d", i)
}

func (w *woocommece) Export() (err error) {
	//_, err = writer.Write([]byte("dsf"))
	groups, err := w.storage.Groups("")
	if err != nil {
		return err
	}
	var id = 0
	terms := Terms(&id, 0, groups)

	b := builder{
		squirrel.Insert("terms").Columns("term_id", "name", "slug", "parent"),
	}

	b.Terms(terms)
	fmt.Println(squirrel.DebugSqlizer(b))

	return
}

func Terms(startID *int, parentID int, groups []commerceml.Group) []term {
	var terms []term
	for _, i := range groups {
		*startID++
		t := term{
			ID:    ID(*startID),
			Name:  i.Name,
			Slug:  Slug(i.Name),
			Group: ID(parentID),
		}
		terms = append(terms, t)
		if len(i.Groups) > 0 {
			childs := Terms(startID, *startID, i.Groups)
			terms = append(terms, childs...)
		}
	}
	return terms
}

type builder struct {
	squirrel.InsertBuilder
}

func (b *builder) Terms(terms []term) {
	for _, i := range terms {
		*b = builder{
			b.Values(squirrel.Expr(i.ID.String()), i.Name, i.Slug, squirrel.Expr(i.Group.String())),
		}
	}
}
