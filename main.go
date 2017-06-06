package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/gosimple/slug"
	squirrel "gopkg.in/Masterminds/squirrel.v1"
	"github.com/imega-teleport/db2file/mysql"
	"github.com/imega-teleport/db2file/exporter"
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
	user, pass, host, dbname := "root", "", "10.0.3.94:3306", "teleport"
	dsn := fmt.Sprintf("mysql://%s:%s@tcp(%s)/%s", user, pass, host, dbname)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		fmt.Printf("error: %v", err)
		os.Exit(1)
	}

	err = db.Ping()
	if err != nil {
		fmt.Printf("error: %v", err)
		os.Exit(1)
	}
	defer func() {
		err := db.Close()
		if err != nil {
			fmt.Printf("error: %v", err)
			return
		}
		fmt.Println("Closed db connection")
	}()

	storage := mysql.NewStorage(db)
	exporter := exporter.NewExporter(storage)

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
