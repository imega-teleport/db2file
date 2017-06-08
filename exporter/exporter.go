package exporter

import (
	"fmt"

	slugmaker "github.com/gosimple/slug"
	"github.com/imega-teleport/db2file/storage"
	"github.com/imega-teleport/xml2db/commerceml"
	"gopkg.in/Masterminds/squirrel.v1"
)

type woocommece struct {
	storage storage.Store
	prefix  string
}

// NewExporter return woocommece instance
func NewExporter(storage storage.Store, prefix string) *woocommece {
	return &woocommece{
		storage: storage,
		prefix:  prefix,
	}
}

type term struct {
	ID    termID
	Name  string
	Slug  slug
	Group termID
}

type slug string

func (s slug) String() string {
	return slugmaker.Make(string(s))
}

type termID int

func (i termID) String() string {
	return fmt.Sprintf("@max_term_id+%d", i)
}

type termTaxonomy struct {
	ID          taxonomyID
	TermID      termID
	Taxonomy    string
	Description string
	Parent      taxonomyID
	Count       int
}

type taxonomyID int

func (i taxonomyID) String() string {
	return fmt.Sprintf("@max_term_taxonomy_id+%d", i)
}

func (w *woocommece) Export() (err error) {
	//_, err = writer.Write([]byte("dsf"))
	groups, err := w.storage.Groups("")
	if err != nil {
		return err
	}
	var startTermID, startTaxonomyID = 1, 0
	terms, termsTaxonomy := Terms(&startTermID, startTaxonomyID, groups)

	b := w.builderTerm()
	b.Terms(terms)
	fmt.Println(squirrel.DebugSqlizer(b))

	btt := w.builderTermTaxonomy()
	btt.TermsTaxonomy(0, termsTaxonomy)
	fmt.Println(squirrel.DebugSqlizer(btt))

	return
}

func Terms(startTermID *int, startTaxonomyID int, groups []commerceml.Group) ([]term, []termTaxonomy) {
	var terms []term
	var termsTaxonomy []termTaxonomy
	for _, i := range groups {
		parentTaxonomyID := startTaxonomyID
		if startTaxonomyID == 0 {
			parentTaxonomyID = *startTermID
		}

		t := term{
			ID:   termID(*startTermID),
			Name: i.Name,
			Slug: slug(i.Name),
		}
		terms = append(terms, t)

		tt := termTaxonomy{
			ID:          taxonomyID(*startTermID),
			TermID:      termID(*startTermID),
			Taxonomy:    "product_cat",
			Description: i.Name, //group.description
			Parent:      taxonomyID(parentTaxonomyID),
		}
		termsTaxonomy = append(termsTaxonomy, tt)
		*startTermID++
		if len(i.Groups) > 0 {
			childsTerms, childsTermsTaxonomy := Terms(startTermID, parentTaxonomyID, i.Groups)
			terms = append(terms, childsTerms...)
			termsTaxonomy = append(termsTaxonomy, childsTermsTaxonomy...)
		}

	}
	return terms, termsTaxonomy
}

type builder struct {
	squirrel.InsertBuilder
}

func (b *builder) Terms(terms []term) {
	for _, i := range terms {
		*b = builder{
			b.Values(squirrel.Expr(i.ID.String()), i.Name, i.Slug, 0),
		}
	}
}

func (b *builder) TermsTaxonomy(taxonomyID int, t []termTaxonomy) {
	for _, i := range t {
		*b = builder{
			b.Values(squirrel.Expr(i.ID.String()), squirrel.Expr(i.TermID.String()), i.Taxonomy, i.Description, squirrel.Expr(i.Parent.String()), 0),
		}
	}
}

func (w *woocommece) builderTerm() builder {
	return builder{
		squirrel.Insert(w.prefix + "terms").Columns("term_id", "name", "slug", "parent"),
	}
}

func (w *woocommece) builderTermTaxonomy() builder {
	return builder{
		squirrel.Insert(w.prefix + "term_taxonomy").Columns("term_taxonomy_id", "term_id", "taxonomy", "description", "parent", "count"),
	}
}
