package paginator

import (
	"fmt"
	"io"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/Masterminds/squirrel.v1"
)

func Test_Paginator_Processing(t *testing.T) {
	type termTaxonomy struct {
		ID          int
		TermID      int
		Taxonomy    string
		Description string
		Parent      int
		Count       int
	}

	terms := []termTaxonomy{
		termTaxonomy{
			ID:          1,
			TermID:      1,
			Taxonomy:    "product_cat",
			Description: "name1",
			Parent:      1,
			Count:       0,
		},
		termTaxonomy{
			ID:          2,
			TermID:      2,
			Taxonomy:    "product_cat",
			Description: "name1-1",
			Parent:      1,
			Count:       0,
		},
		termTaxonomy{
			ID:          3,
			TermID:      3,
			Taxonomy:    "product_cat",
			Description: "name1-2",
			Parent:      1,
			Count:       0,
		},
		termTaxonomy{
			ID:          4,
			TermID:      4,
			Taxonomy:    "product_cat",
			Description: "name2",
			Parent:      4,
			Count:       0,
		},
		termTaxonomy{
			ID:          5,
			TermID:      5,
			Taxonomy:    "product_cat",
			Description: "name2-1",
			Parent:      4,
			Count:       0,
		},
		termTaxonomy{
			ID:          5,
			TermID:      5,
			Taxonomy:    "product_cat",
			Description: "name2-2",
			Parent:      4,
			Count:       0,
		},
	}

	values := make([]interface{}, len(terms))
	for i, v := range terms {
		values[i] = v
	}

	r := paginator{4}

	s := squirrel.Insert("term_taxonomy").Columns("description")

	pr, pw := io.Pipe()
	defer pr.Close()
	defer pw.Close()
	go func() {
		r.Processing(
			values,
			func(i interface{}) interface{} {
				s = s.Values(i.(termTaxonomy).Description)
				return false
			},
			func(p interface{}) interface{} {
				pw.Write([]byte(fmt.Sprintf("%s;", squirrel.DebugSqlizer(s))))
				s = squirrel.Insert("term_taxonomy").Columns("description")
				return false
			},
		)
		pw.Close()
	}()
	body, _ := ioutil.ReadAll(pr)

	assert.Equal(
		t,
		"INSERT INTO term_taxonomy (description) VALUES ('name1'),('name1-1'),('name1-2'),('name2');INSERT INTO term_taxonomy (description) VALUES ('name2-1'),('name2-2');",
		fmt.Sprintf("%s", body),
	)
}
