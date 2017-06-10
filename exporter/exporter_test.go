package exporter

import (
	"testing"

	"github.com/imega-teleport/xml2db/commerceml"
	"github.com/stretchr/testify/assert"
	"gopkg.in/Masterminds/squirrel.v1"
)

func Test_Terms_WithGroupLevel1_ReturnTerm(t *testing.T) {
	groups := []commerceml.Group{
		{
			IdName: commerceml.IdName{
				Id:   "id1",
				Name: "name1",
			},
		},
		{
			IdName: commerceml.IdName{
				Id:   "id2",
				Name: "name2",
			},
		},
	}

	var startTermID, startTaxonomyID = 1, 0
	terms, _ := Terms(&startTermID, startTaxonomyID, groups)

	assert.Equal(t, []term{
		{
			ID:    1,
			Name:  "name1",
			Slug:  "name1",
			Group: 0,
		},
		{
			ID:    2,
			Name:  "name2",
			Slug:  "name2",
			Group: 0,
		},
	}, terms)
}

func Test_Terms_WithGroupLevel2_ReturnTerm(t *testing.T) {
	groups := []commerceml.Group{
		{
			IdName: commerceml.IdName{
				Id:   "id1",
				Name: "name1",
			},
			Groups: []commerceml.Group{
				{
					IdName: commerceml.IdName{
						Id:   "id1-1",
						Name: "name1-1",
					},
				},
				{
					IdName: commerceml.IdName{
						Id:   "id1-2",
						Name: "name1-2",
					},
				},
			},
		},
		{
			IdName: commerceml.IdName{
				Id:   "id2",
				Name: "name2",
			},
			Groups: []commerceml.Group{
				{
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
		{
			ID:    1,
			Name:  "name1",
			Slug:  "name1",
			Group: 0,
		},
		{
			ID:    2,
			Name:  "name1-1",
			Slug:  "name1-1",
			Group: 0,
		},
		{
			ID:    3,
			Name:  "name1-2",
			Slug:  "name1-2",
			Group: 0,
		},
		{
			ID:    4,
			Name:  "name2",
			Slug:  "name2",
			Group: 0,
		},
		{
			ID:    5,
			Name:  "name2-1",
			Slug:  "name2-1",
			Group: 0,
		},
	}, terms)
}

func Test_Terms_WithGroupLevel3_ReturnTerm(t *testing.T) {
	groups := []commerceml.Group{
		{
			IdName: commerceml.IdName{
				Id:   "id1",
				Name: "name1",
			},
			Groups: []commerceml.Group{
				{
					IdName: commerceml.IdName{
						Id:   "id1-1",
						Name: "name1-1",
					},
					Groups: []commerceml.Group{
						{
							IdName: commerceml.IdName{
								Id:   "id1-1-1",
								Name: "name1-1-1",
							},
						},
						{
							IdName: commerceml.IdName{
								Id:   "id1-1-2",
								Name: "name1-1-2",
							},
						},
					},
				},
				{
					IdName: commerceml.IdName{
						Id:   "id1-2",
						Name: "name1-2",
					},
					Groups: []commerceml.Group{
						{
							IdName: commerceml.IdName{
								Id:   "id1-2-1",
								Name: "name1-2-1",
							},
						},
					},
				},
			},
		},
		{
			IdName: commerceml.IdName{
				Id:   "id2",
				Name: "name2",
			},
			Groups: []commerceml.Group{
				{
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
		{
			ID:    1,
			Name:  "name1",
			Slug:  "name1",
			Group: 0,
		},
		{
			ID:    2,
			Name:  "name1-1",
			Slug:  "name1-1",
			Group: 0,
		},
		{
			ID:    3,
			Name:  "name1-1-1",
			Slug:  "name1-1-1",
			Group: 0,
		},
		{
			ID:    4,
			Name:  "name1-1-2",
			Slug:  "name1-1-2",
			Group: 0,
		},
		{
			ID:    5,
			Name:  "name1-2",
			Slug:  "name1-2",
			Group: 0,
		},
		{
			ID:    6,
			Name:  "name1-2-1",
			Slug:  "name1-2-1",
			Group: 0,
		},
		{
			ID:    7,
			Name:  "name2",
			Slug:  "name2",
			Group: 0,
		},
		{
			ID:    8,
			Name:  "name2-1",
			Slug:  "name2-1",
			Group: 0,
		},
	}, terms)
}

func Test_TermsTaxonomy_WithGroupLevel1_ReturnTermTaxonomy(t *testing.T) {
	groups := []commerceml.Group{
		{
			IdName: commerceml.IdName{
				Id:   "id1",
				Name: "name1",
			},
		},
		{
			IdName: commerceml.IdName{
				Id:   "id2",
				Name: "name2",
			},
		},
	}

	var startTermID, startTaxonomyID = 1, 0
	_, termsTaxonomy := Terms(&startTermID, startTaxonomyID, groups)

	assert.Equal(t, []termTaxonomy{
		{
			ID:          1,
			TermID:      1,
			Taxonomy:    "product_cat",
			Description: "name1",
			Parent:      1,
			Count:       0,
		},
		{
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
		{
			IdName: commerceml.IdName{
				Id:   "id1",
				Name: "name1",
			},
			Groups: []commerceml.Group{
				{
					IdName: commerceml.IdName{
						Id:   "id1-1",
						Name: "name1-1",
					},
				},
				{
					IdName: commerceml.IdName{
						Id:   "id1-2",
						Name: "name1-2",
					},
				},
			},
		},
		{
			IdName: commerceml.IdName{
				Id:   "id2",
				Name: "name2",
			},
			Groups: []commerceml.Group{
				{
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
		{
			ID:          1,
			TermID:      1,
			Taxonomy:    "product_cat",
			Description: "name1",
			Parent:      1,
			Count:       0,
		},
		{
			ID:          2,
			TermID:      2,
			Taxonomy:    "product_cat",
			Description: "name1-1",
			Parent:      1,
			Count:       0,
		},
		{
			ID:          3,
			TermID:      3,
			Taxonomy:    "product_cat",
			Description: "name1-2",
			Parent:      1,
			Count:       0,
		},
		{
			ID:          4,
			TermID:      4,
			Taxonomy:    "product_cat",
			Description: "name2",
			Parent:      4,
			Count:       0,
		},
		{
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
