package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gosimple/slug"
	"github.com/imega-teleport/db2file/exporter"
	"github.com/imega-teleport/db2file/mysql"
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
	woo := exporter.NewExporter(storage, "")
	woo.Export()
}

func check(e error) {
	if e != nil {
		fmt.Printf("error: %v", e)
		os.Exit(1)
	}
}
