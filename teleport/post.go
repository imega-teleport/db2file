package teleport

import (
	"fmt"
	"time"

	"gopkg.in/Masterminds/squirrel.v1"
)

// Post is a record for WP
type Post struct {
	ID       UUID
	AuthorID int
	Date     time.Time
	Content  string
	Title    string
	Excerpt  string
	Name     string
	Modified time.Time
	Type     string
	ParentID UUID
}

// SizeOf return size a query
func (p Post) SizeOf() int {
	return len(p.ID) + (lengthDefineDate * 4) + len(p.Name) + len(p.Title) + len(p.Content) + len(p.Excerpt) + lengthDefineIndex
}

// BuilderPost return builder
func (w *Wpwc) BuilderPost() builder {
	return builder{
		squirrel.Insert(fmt.Sprintf("%sposts", w.Prefix)).Columns(
			"id",
			"post_author",
			"post_date",
			"post_date_gmt",
			"post_content",
			"post_title",
			"post_excerpt",
			"post_name",
			"post_modified",
			"post_modified_gmt",
			"post_parent",
			"post_type",
		),
	}
}

func (b *builder) AddPost(post Post) {
	*b = builder{
		b.Values(
			squirrel.Expr(fmt.Sprintf("@max_post_id+%s", post.ID.ToVar())),
			squirrel.Expr("1"),
			post.Date.Format("2006-01-02 15:04:05"),
			post.Date.UTC().Format("2006-01-02 15:04:05"),
			b.MysqlRealEscapeString(post.Content),
			b.MysqlRealEscapeString(post.Title),
			b.MysqlRealEscapeString(post.Excerpt),
			b.MysqlRealEscapeString(post.Name),
			post.Modified.Format("2006-01-02 15:04:05"),
			post.Modified.UTC().Format("2006-01-02 15:04:05"),
			squirrel.Expr(post.ParentID.ToVar()),
			post.Type,
		),
	}
}
