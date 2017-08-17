package main // import "github.com/imega-teleport/db2file"

import (
	"database/sql"
	"flag"
	"fmt"
	"io"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/imega-teleport/db2file/exporter"
	"github.com/imega-teleport/db2file/mysql"
)

func main() {
	user, pass, host := os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_HOST")

	dbname := flag.String("db", "", "Database name")
	path := flag.String("path", "", "Save to path")
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

	complete, err := storage.CheckCompleteAllTasks()
	if err != nil {
		fmt.Printf("Fail check complete task: %s", err)
		os.Exit(1)
	}

	if !complete {
		os.Exit(1)
	}

	woo := exporter.NewExporter(storage, "wp_", 1)
	r, w := io.Pipe()
	defer func() {
		err = r.Close()
		err = w.Close()
	}()

	woo.Export(w)

	file, err := os.Create(fmt.Sprintf("%s%c%s", *path, os.PathSeparator, "output.sql"))
	if err != nil {
		fmt.Printf("Could not create file: %s", err)
		os.Exit(1)
	}
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Printf("Fail close file: %s", err)
			os.Exit(1)
		}
	}()
	buf := make([]byte, 1024)
	for {
		n, err := r.Read(buf)
		if err != nil && err != io.EOF {
			fmt.Printf("Fail read from buffer: %s", err)
			os.Exit(1)
		}
		if n == 0 {
			break
		}
		if _, err := file.Write(buf[:n]); err != nil {
			fmt.Printf("Fail write to file: %s", err)
			os.Exit(1)
		}
	}
}
