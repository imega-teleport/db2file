package exporter

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/imega-teleport/db2file/storage"
	"github.com/imega-teleport/xml2db/commerceml"
	"github.com/stretchr/testify/assert"
	squirrel "gopkg.in/Masterminds/squirrel.v1"
)

func EqualWithDuration(t *testing.T, expected, actual interface{}, delta time.Duration) bool {
	a := reflect.ValueOf(actual)
	e := reflect.ValueOf(expected)

	for i := 0; i < a.NumField(); i++ {
		if a.Field(i).Type() == reflect.TypeOf(time.Time{}) {
			assert.WithinDuration(t, e.Field(i).Interface().(time.Time), a.Field(i).Interface().(time.Time), delta, fmt.Sprintf("Field: %s", a.Type().Field(i).Name))
		} else {
			assert.Equal(t, e.Field(i).Interface(), a.Field(i).Interface())
		}
	}
	return true
}

func Test_makeTerms_WithGroupLevel1_ReturnTerm(t *testing.T) {
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

	var startTermID, startTaxonomyID = 0, 0
	terms, _ := makeTerms(&startTermID, startTaxonomyID, groups)

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

func Test_makeTerms_WithGroupLevel2_ReturnTerm(t *testing.T) {
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

	var startTermID, startTaxonomyID = 0, 0
	terms, _ := makeTerms(&startTermID, startTaxonomyID, groups)

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

func Test_makeTerms_WithGroupLevel3_ReturnTerm(t *testing.T) {
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

	var startTermID, startTaxonomyID = 0, 0
	terms, _ := makeTerms(&startTermID, startTaxonomyID, groups)

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

func Test_makeTerms_WithGroupLevel1_ReturnTermTaxonomy(t *testing.T) {
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

	var startTermID, startTaxonomyID = 0, 0
	_, termsTaxonomy := makeTerms(&startTermID, startTaxonomyID, groups)

	assert.Equal(t, []termTaxonomy{
		{
			ID:          1,
			TermID:      1,
			Taxonomy:    "product_cat",
			Description: "name1",
			Parent:      0,
			Count:       0,
		},
		{
			ID:          2,
			TermID:      2,
			Taxonomy:    "product_cat",
			Description: "name2",
			Parent:      0,
			Count:       0,
		},
	}, termsTaxonomy)
}

func Test_makeTerms_WithGroupLevel2_ReturnTermTaxonomy(t *testing.T) {
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

	var startTermID, startTaxonomyID = 0, 0
	_, termsTaxonomy := makeTerms(&startTermID, startTaxonomyID, groups)

	assert.Equal(t, []termTaxonomy{
		{
			ID:          1,
			TermID:      1,
			Taxonomy:    "product_cat",
			Description: "name1",
			Parent:      0,
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
			Parent:      0,
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
	assert.Equal(t, "INSERT INTO test_terms (term_id,name,slug,term_group) VALUES ('')", squirrel.DebugSqlizer(b))
}

func Test_builderTermTaxonomy_WithPrefix_ReturnBuilder(t *testing.T) {
	w := &woocommece{
		prefix: "test_",
	}
	b := w.builderTermTaxonomy().Values("")
	assert.Equal(t, "INSERT INTO test_term_taxonomy (term_taxonomy_id,term_id,taxonomy,description,parent,count) VALUES ('')", squirrel.DebugSqlizer(b))
}

func TestTaxonomyID_WithIndex_ReturnsVariableWithIndex(t *testing.T) {
	assert.Equal(t, "@max_term_taxonomy_id+1", taxonomyID(1).String())
}

func TestTaxonomyID_WithIndexZero_ReturnsZero(t *testing.T) {
	assert.Equal(t, "0", taxonomyID(0).String())
}

func Test_slug_WithString_ReturnsDecodedSlug(t *testing.T) {
	assert.Equal(t, "s_l-a-g", slug("с_л+а-г").String())
}

func Test_termID_WithIndex_ReturnsVariableWithIndex(t *testing.T) {
	assert.Equal(t, termID(1).String(), "@max_term_id+1")
}

func Test_NewExporter_ReturnsExporter(t *testing.T) {
	assert.Implements(t, (*Exporter)(nil), NewExporter((storage.Store)(nil), "", 1))
}

func Test_postID_WithIndex_ReturnsVariableWithIndex(t *testing.T) {
	assert.Equal(t, "@max_post_id+1", postID(1).String())
}

func Test_authorID_WithIndex_ReturnsVariable(t *testing.T) {
	assert.Equal(t, "@author_id", authorID(1).String())
}

func Test_makePosts_WithProducts_ReturnsPosts(t *testing.T) {
	products := []commerceml.Product{
		{
			IdName: commerceml.IdName{
				Id:   "id1",
				Name: "name1",
			},
		},
	}
	actual, _ := makePosts(products)

	expected := []post{
		{
			ID:            postID(1),
			Title:         "name1",
			Status:        "publish",
			CommentStatus: "open",
			PingStatus:    "open",
			Name:          "name1",
			Type:          "product",
			Date:          time.Now(),
			Modified:      time.Now(),
		},
	}

	for k, v := range actual {
		EqualWithDuration(t, expected[k], v, time.Second)
	}
}

func Test_makePosts_WithProducts_ReturnsTermRelationship(t *testing.T) {
	products := []commerceml.Product{
		{
			IdName: commerceml.IdName{
				Id:   "id1",
				Name: "name1",
			},
			Groups: []commerceml.Group{
				{
					IdName: commerceml.IdName{
						Id: "group1",
					},
				},
			},
		},
	}
	_, actual := makePosts(products)

	assert.Equal(t, []termRelationship{
		{
			ObjectType:     typePost,
			ObjectID:       uuid("id1"),
			TermTaxonomyID: uuid("group1"),
		},
	}, actual)
}

func Test_objectType_WithTypePost_ReturnsVariableName(t *testing.T) {
	assert.Equal(t, "max_post_id", objectType(typePost).String())
}

func Test_objectType_WithTypeTerm_ReturnsVariableName(t *testing.T) {
	assert.Equal(t, "max_term_id", objectType(typeTerm).String())
}

func Test_uuid_WithTypeTerm_ReturnsVariableName(t *testing.T) {
	assert.Equal(t, "@abc", uuid("a-b-c").ToVar())
}
