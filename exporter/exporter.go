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
}

// NewExporter return woocommece instance
func NewExporter(storage storage.Store) *woocommece {
	return &woocommece{
		storage: storage,
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
	var startTermID = 0
	var startTaxonomyID = 0
	terms, termsTaxonomy := Terms(&startTermID, &startTaxonomyID, 0, groups)

	b := builder{
		squirrel.Insert("terms").Columns("term_id", "name", "slug", "parent"),
	}

	b.Terms(terms)
	fmt.Println(squirrel.DebugSqlizer(b))

	btt := builderTermTaxonomy("")
	btt.TermsTaxonomy(0, termsTaxonomy)
	fmt.Println(squirrel.DebugSqlizer(btt))

	return
}

func Terms(startTermID *int, startTaxonomyID *int, parentID int, groups []commerceml.Group) ([]term, []termTaxonomy) {
	var terms []term
	var termsTaxonomy []termTaxonomy
	for _, i := range groups {
		*startTermID++
		*startTaxonomyID++
		t := term{
			ID:    termID(*startTermID),
			Name:  i.Name,
			Slug:  slug(i.Name),
			Group: termID(parentID),
		}
		terms = append(terms, t)
		termTaxonomyParent := startTermID
		if *startTaxonomyID > 0 {
			termTaxonomyParent = startTaxonomyID
		}
		tt := termTaxonomy{
			ID: taxonomyID(*startTaxonomyID),
			TermID: termID(*startTermID),
			Taxonomy: "product_cat",
			Description: i.Name, //group.description
			Parent: taxonomyID(*termTaxonomyParent),
		}
		termsTaxonomy = append(termsTaxonomy, tt)
		if len(i.Groups) > 0 {
			childsTerms, childsTermsTaxonomy := Terms(startTermID, startTaxonomyID, *startTermID, i.Groups)
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

func builderTermTaxonomy(prefix string) builder {
	return builder{
		squirrel.Insert(prefix + "term_taxonomy").Columns("term_taxonomy_id", "term_id", "taxonomy", "description", "parent", "count"),
	}
}