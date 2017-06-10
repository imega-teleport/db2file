package exporter

import (
	"fmt"
	"io"
	"reflect"
	"strings"
	"time"

	slugmaker "github.com/gosimple/slug"
	"github.com/imega-teleport/db2file/paginator"
	"github.com/imega-teleport/db2file/storage"
	"github.com/imega-teleport/xml2db/commerceml"
	squirrel "gopkg.in/Masterminds/squirrel.v1"
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

func (w *woocommece) Export(writer *io.PipeWriter) (err error) {
	groups, err := w.storage.Groups("")
	if err != nil {
		return err
	}
	var startTermID, startTaxonomyID = 1, 0
	terms, termsTaxonomy := Terms(&startTermID, startTaxonomyID, groups)

	posts, err := w.storage.Posts("")
	if err != nil {
		return err
	}

	go func() {
		writer.Write([]byte("start transaction;\n"))

		for k, v := range variables {
			key := strings.Replace(slugmaker.Make(k), "-", "", -1)
			writer.Write([]byte(fmt.Sprintf("set @%s=%d;\n", key, v)))
		}

		values := make([]interface{}, len(terms))
		for i, v := range terms {
			values[i] = v
		}

		b := w.builderTerm()

		p := paginator.NewPaginator(100000)

		p.Processing(
			values,
			func(t interface{}) interface{} {
				b.AddTerm(t.(term))
				return false
			},
			func(interface{}) interface{} {
				writer.Write([]byte(fmt.Sprintf("%s;\n", squirrel.DebugSqlizer(b))))
				b = w.builderTerm()
				return false
			},
		)

		bt := w.builderTermTaxonomy()
		values1 := make([]interface{}, len(termsTaxonomy))
		for i, v := range termsTaxonomy {
			values1[i] = v
		}

		p.Processing(
			values1,
			func(t interface{}) interface{} {
				bt.AddTermsTaxonomy(t.(termTaxonomy))
				return false
			},
			func(interface{}) interface{} {
				writer.Write([]byte(fmt.Sprintf("%s;\n", squirrel.DebugSqlizer(bt))))
				bt = w.builderTermTaxonomy()
				return false
			},
		)

		bp := w.builderPost()
		valP := make([]interface{}, len(posts))
		for i, v := range posts {
			valP[i] = v
		}
		p.Processing(
			valP,
			func(p interface{}) interface{} {
				bp.AddPost(p.(post))
				return false
			},
			func(interface{}) interface{} {
				writer.Write([]byte(fmt.Sprintf("%s;\n", squirrel.DebugSqlizer(bp))))
				bp = w.builderPost()
				return false
			},
		)

		writer.Write([]byte("commit;\n"))
		writer.Close()
	}()

	return
}

var variables = make(map[string]int)

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

		variables[i.Id] = *startTermID
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

type post struct {
	ID              postID
	AuthorID        authorID
	Date            time.Time
	Content         string
	Title           string
	Excerpt         string
	Status          string //publish
	CommentStatus   string //open
	PingStatus      string //open
	Password        string
	Name            string
	ToPing          string
	Pinged          string
	Modified        time.Time
	ContentFiltered string
	ParentID        postID
	GUID            string
	MenuOrder       int
	Type            string //post
	MimeType        string
	CommentCount    int
}

type postID int

func (i postID) String() string {
	return fmt.Sprintf("@max_post_id+%d", i)
}

type authorID int

func (i authorID) String() string {
	return "@author_id"
}

func Posts(products []commerceml.Product) ([]post, []termRelationship) {
	var posts []post
	var rels []termRelationship
	var startID = 1
	for _, i := range products {
		p := post{
			ID:      postID(startID),
			Content: i.Description.Value,
		}
		startID++

		for _, g := range i.Groups {
			tr := termRelationship{
				ObjectType:     &p,
				ObjectID:       i.Id,
				TermTaxonomyID: g.Id,
			}
			rels = append(rels, tr)
		}
	}
	return posts, rels
}

type termRelationship struct {
	ObjectType     interface{}
	ObjectID       string
	TermTaxonomyID string
	TermOrder      int
}

type builder struct {
	squirrel.InsertBuilder
}

func (b *builder) AddTerm(t term) {
	*b = builder{
		b.Values(squirrel.Expr(t.ID.String()), t.Name, t.Slug, 0),
	}
}

func (b *builder) AddTermsTaxonomy(t termTaxonomy) {
	*b = builder{
		b.Values(squirrel.Expr(t.ID.String()), squirrel.Expr(t.TermID.String()), t.Taxonomy, t.Description, squirrel.Expr(t.Parent.String()), 0),
	}
}

func (b *builder) AddPost(post post) {
	*b = builder{
		b.Values(
			squirrel.Expr(post.ID.String()),
			post.AuthorID.String(),
			post.Date.String(),
			post.Date.UTC().String(),
			post.Content,
			post.ContentFiltered,
			post.Title,
			post.Excerpt,
			post.Status,
			post.Type,
			post.CommentStatus,
			post.PingStatus,
			post.Password,
			post.Name,
			post.ToPing,
			post.Pinged,
			post.Modified.String(),
			post.Modified.UTC().String(),
			squirrel.Expr(post.ParentID.String()),
			post.MenuOrder,
			post.MimeType,
			post.GUID,
		),
	}
}

func (b *builder) AddTermRelationships(r termRelationship) {
	var prefix string
	switch reflect.TypeOf(r.ObjectType) {
	case post{}:
		prefix = "max_post_id"
	case term{}:
		prefix = "max_term_id"
	}

	*b = builder{
		b.Values(
			squirrel.Expr(fmt.Sprintf("@%s+@%s", prefix, r.ObjectID)),
			squirrel.Expr(fmt.Sprintf("@max_term_id+@%s", r.TermTaxonomyID)),
			r.TermOrder,
		),
	}
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
		squirrel.Insert(w.prefix+"terms").Columns("term_id", "name", "slug", "parent"),
	}
}

func (w *woocommece) builderTermTaxonomy() builder {
	return builder{
		squirrel.Insert(w.prefix+"term_taxonomy").Columns("term_taxonomy_id", "term_id", "taxonomy", "description", "parent", "count"),
	}
}

func (w *woocommece) builderPost() builder {
	return builder{
		squirrel.Insert(w.prefix+"posts").Columns(
			"post_author",
			"post_date",
			"post_date_gmt",
			"post_content",
			"post_content_filtered",
			"post_title",
			"post_excerpt",
			"post_status",
			"post_type",
			"comment_status",
			"ping_status",
			"post_password",
			"post_name",
			"to_ping",
			"pinged",
			"post_modified",
			"post_modified_gmt",
			"post_parent",
			"menu_order",
			"post_mime_type",
			"guid",
		),
	}
}

func (w *woocommece) builderTermRelationships() builder {
	return builder{
		squirrel.Insert(w.prefix+"term_relationships").Columns("object_id", "term_taxonomy_id", "term_order"),
	}
}
