package main // import "github.com/imega-teleport/db2file"

import (
	"database/sql"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/imega-teleport/db2file/exporter"
	"github.com/imega-teleport/db2file/mysql"
)

func main() {
	user, pass, host := os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_HOST")

	dbname := flag.String("db", "", "Database name")
	flag.Parse()

	dsn := fmt.Sprintf("mysql://%s:%s@tcp(%s)/%s", user, pass, host, *dbname)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		fmt.Printf("error: %s", err)
		os.Exit(1)
	}

	err = db.Ping()
	if err != nil {
		fmt.Printf("error: %s", err)
		os.Exit(1)
	}
	defer func() {
		err = db.Close()
		if err != nil {
			fmt.Printf("error: %s", err)
			os.Exit(1)
		}
		fmt.Println("Closed db connection")
	}()

	storage := mysql.NewStorage(db)
	woo := exporter.NewExporter(storage, "")
	r, w := io.Pipe()
	defer func() {
		err = r.Close()
		err = w.Close()
	}()

	woo.Export(w)
	/*if err := woo.Export(w); err != nil {
		fmt.Printf("Error in export: %v", err)
		os.Exit(1)
	}*/
	body, _ := ioutil.ReadAll(r)
	fmt.Printf("%s", body)
}
