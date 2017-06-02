package main

import (
	"fmt"
	"os"
	squirrel "gopkg.in/Masterminds/squirrel.v1"
)

func main() {
	builder := squirrel.Insert("terms")
	builder = builder.Columns("term", "value").Values("1", "2")
	builder = builder.Values([]byte("A"),"6")

	q, args, err := builder.ToSql()
	fmt.Println(q, args, err)
	query := squirrel.DebugSqlizer(builder)
	fmt.Println(query)
}

func check(e error) {
	if e != nil {
		fmt.Printf("error: %v", e)
		os.Exit(1)
	}
}
