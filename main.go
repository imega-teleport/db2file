package main

import (
	"fmt"
	"os"

	"github.com/gosimple/slug"
	squirrel "gopkg.in/Masterminds/squirrel.v1"
)

type Term struct {
	ID    ID
	Name  string
	Slug  Slug
	Group int
}

type Slug string

func (s Slug) String() string {
	return slug.Make(string(s))
}

type ID int

func (i ID) String() string {
	return fmt.Sprintf("@max_term_id+%d", i)
}

func main() {
	t := Term{
		ID:   1,
		Name: "name",
		Slug: "Имя Фамилия",
	}

	builder := squirrel.Insert("terms")
	builder = builder.Columns("term_id", "name", "slug")
	builder = builder.Values(squirrel.Expr(t.ID.String()), t.Name, t.Slug)

	query := squirrel.DebugSqlizer(builder)
	fmt.Println(query)
}

func check(e error) {
	if e != nil {
		fmt.Printf("error: %v", e)
		os.Exit(1)
	}
}
