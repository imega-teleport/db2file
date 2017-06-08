package exporter

import (
	"testing"

	"math"

	"github.com/imega-teleport/xml2db/commerceml"
	"github.com/stretchr/testify/assert"
	"gopkg.in/Masterminds/squirrel.v1"
	"fmt"
)

func Test_Terms_WithGroupLevel1_ReturnTerm(t *testing.T) {
	groups := []commerceml.Group{
		commerceml.Group{
			IdName: commerceml.IdName{
				Id:   "id1",
				Name: "name1",
			},
		},
		commerceml.Group{
			IdName: commerceml.IdName{
				Id:   "id2",
				Name: "name2",
			},
		},
	}

	var startTermID, startTaxonomyID = 1, 0
	terms, _ := Terms(&startTermID, startTaxonomyID, groups)

	assert.Equal(t, []term{
		term{
			ID:    1,
			Name:  "name1",
			Slug:  "name1",
			Group: 0,
		},
		term{
			ID:    2,
			Name:  "name2",
			Slug:  "name2",
			Group: 0,
		},
	}, terms)
}

func Test_Terms_WithGroupLevel2_ReturnTerm(t *testing.T) {
	groups := []commerceml.Group{
		commerceml.Group{
			IdName: commerceml.IdName{
				Id:   "id1",
				Name: "name1",
			},
			Groups: []commerceml.Group{
				commerceml.Group{
					IdName: commerceml.IdName{
						Id:   "id1-1",
						Name: "name1-1",
					},
				},
				commerceml.Group{
					IdName: commerceml.IdName{
						Id:   "id1-2",
						Name: "name1-2",
					},
				},
			},
		},
		commerceml.Group{
			IdName: commerceml.IdName{
				Id:   "id2",
				Name: "name2",
			},
			Groups: []commerceml.Group{
				commerceml.Group{
					IdName: commerceml.IdName{
						Id:   "id2-1",
						Name: "name2-1",
					},
				},
			},
		},
	}

	var startTermID, startTaxonomyID = 1, 0
	terms, _ := Terms(&startTermID, startTaxonomyID, groups)

	assert.Equal(t, []term{
		term{
			ID:    1,
			Name:  "name1",
			Slug:  "name1",
			Group: 0,
		},
		term{
			ID:    2,
			Name:  "name1-1",
			Slug:  "name1-1",
			Group: 0,
		},
		term{
			ID:    3,
			Name:  "name1-2",
			Slug:  "name1-2",
			Group: 0,
		},
		term{
			ID:    4,
			Name:  "name2",
			Slug:  "name2",
			Group: 0,
		},
		term{
			ID:    5,
			Name:  "name2-1",
			Slug:  "name2-1",
			Group: 0,
		},
	}, terms)
}

func Test_Terms_WithGroupLevel3_ReturnTerm(t *testing.T) {
	groups := []commerceml.Group{
		commerceml.Group{
			IdName: commerceml.IdName{
				Id:   "id1",
				Name: "name1",
			},
			Groups: []commerceml.Group{
				commerceml.Group{
					IdName: commerceml.IdName{
						Id:   "id1-1",
						Name: "name1-1",
					},
					Groups: []commerceml.Group{
						commerceml.Group{
							IdName: commerceml.IdName{
								Id:   "id1-1-1",
								Name: "name1-1-1",
							},
						},
						commerceml.Group{
							IdName: commerceml.IdName{
								Id:   "id1-1-2",
								Name: "name1-1-2",
							},
						},
					},
				},
				commerceml.Group{
					IdName: commerceml.IdName{
						Id:   "id1-2",
						Name: "name1-2",
					},
					Groups: []commerceml.Group{
						commerceml.Group{
							IdName: commerceml.IdName{
								Id:   "id1-2-1",
								Name: "name1-2-1",
							},
						},
					},
				},
			},
		},
		commerceml.Group{
			IdName: commerceml.IdName{
				Id:   "id2",
				Name: "name2",
			},
			Groups: []commerceml.Group{
				commerceml.Group{
					IdName: commerceml.IdName{
						Id:   "id2-1",
						Name: "name2-1",
					},
				},
			},
		},
	}

	var startTermID, startTaxonomyID = 1, 0
	terms, _ := Terms(&startTermID, startTaxonomyID, groups)

	assert.Equal(t, []term{
		term{
			ID:    1,
			Name:  "name1",
			Slug:  "name1",
			Group: 0,
		},
		term{
			ID:    2,
			Name:  "name1-1",
			Slug:  "name1-1",
			Group: 0,
		},
		term{
			ID:    3,
			Name:  "name1-1-1",
			Slug:  "name1-1-1",
			Group: 0,
		},
		term{
			ID:    4,
			Name:  "name1-1-2",
			Slug:  "name1-1-2",
			Group: 0,
		},
		term{
			ID:    5,
			Name:  "name1-2",
			Slug:  "name1-2",
			Group: 0,
		},
		term{
			ID:    6,
			Name:  "name1-2-1",
			Slug:  "name1-2-1",
			Group: 0,
		},
		term{
			ID:    7,
			Name:  "name2",
			Slug:  "name2",
			Group: 0,
		},
		term{
			ID:    8,
			Name:  "name2-1",
			Slug:  "name2-1",
			Group: 0,
		},
	}, terms)
}

func Test_TermsTaxonomy_WithGroupLevel1_ReturnTermTaxonomy(t *testing.T) {
	groups := []commerceml.Group{
		commerceml.Group{
			IdName: commerceml.IdName{
				Id:   "id1",
				Name: "name1",
			},
		},
		commerceml.Group{
			IdName: commerceml.IdName{
				Id:   "id2",
				Name: "name2",
			},
		},
	}

	var startTermID, startTaxonomyID = 1, 0
	_, termsTaxonomy := Terms(&startTermID, startTaxonomyID, groups)

	assert.Equal(t, []termTaxonomy{
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
			Description: "name2",
			Parent:      2,
			Count:       0,
		},
	}, termsTaxonomy)
}

func Test_TermsTaxonomy_WithGroupLevel2_ReturnTerm(t *testing.T) {
	groups := []commerceml.Group{
		commerceml.Group{
			IdName: commerceml.IdName{
				Id:   "id1",
				Name: "name1",
			},
			Groups: []commerceml.Group{
				commerceml.Group{
					IdName: commerceml.IdName{
						Id:   "id1-1",
						Name: "name1-1",
					},
				},
				commerceml.Group{
					IdName: commerceml.IdName{
						Id:   "id1-2",
						Name: "name1-2",
					},
				},
			},
		},
		commerceml.Group{
			IdName: commerceml.IdName{
				Id:   "id2",
				Name: "name2",
			},
			Groups: []commerceml.Group{
				commerceml.Group{
					IdName: commerceml.IdName{
						Id:   "id2-1",
						Name: "name2-1",
					},
				},
			},
		},
	}

	var startTermID, startTaxonomyID = 1, 0
	_, termsTaxonomy := Terms(&startTermID, startTaxonomyID, groups)

	assert.Equal(t, []termTaxonomy{
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
	}, termsTaxonomy)
}

func Test_builderTerm_WithPrefix_ReturnBuilder(t *testing.T) {
	w := &woocommece{
		prefix: "test_",
	}
	b := w.builderTerm().Values("")
	assert.Equal(t, "INSERT INTO test_terms (term_id,name,slug,parent) VALUES ('')", squirrel.DebugSqlizer(b))
}

func Test_builderTermTaxonomy_WithPrefix_ReturnBuilder(t *testing.T) {
	w := &woocommece{
		prefix: "test_",
	}
	b := w.builderTermTaxonomy().Values("")
	assert.Equal(t, "INSERT INTO test_term_taxonomy (term_taxonomy_id,term_id,taxonomy,description,parent,count) VALUES ('')", squirrel.DebugSqlizer(b))
}

func Te1st_split(t *testing.T) {
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
		/*termTaxonomy{
			ID:          5,
			TermID:      5,
			Taxonomy:    "product_cat",
			Description: "name2-2",
			Parent:      4,
			Count:       0,
		},*/
	}
	piece := 6
	pages := math.Ceil(float64(len(terms)) / float64(piece))
	for n := 0; n < int(pages); n++ {
		start, end := 0+n*piece, n*piece
		end = end + piece
		if end > len(terms) {
			end = len(terms)
		}
		items := terms[start:end]
		for _, i := range items {
			fmt.Printf("%s,", i.Description)
		}
		fmt.Printf("%s\n", "")
	}
}
