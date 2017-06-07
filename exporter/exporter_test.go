package exporter

import (
	"testing"

	"github.com/imega-teleport/xml2db/commerceml"
	"github.com/stretchr/testify/assert"
)

func TestTerms_WithGroupLevel1_ReturnTerm(t *testing.T) {
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

	var id = 0
	terms := Terms(&id, 0, groups)

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

func TestTerms_WithGroupLevel2_ReturnTerm(t *testing.T) {
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

	var id = 0
	terms := Terms(&id, 0, groups)

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
			Group: 1,
		},
		term{
			ID:    3,
			Name:  "name1-2",
			Slug:  "name1-2",
			Group: 1,
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
			Group: 4,
		},
	}, terms)
}

func TestTerms_WithGroupLevel3_ReturnTerm(t *testing.T) {
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

	var id = 0
	terms := Terms(&id, 0, groups)

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
			Group: 1,
		},
		term{
			ID:    3,
			Name:  "name1-1-1",
			Slug:  "name1-1-1",
			Group: 2,
		},
		term{
			ID:    4,
			Name:  "name1-1-2",
			Slug:  "name1-1-2",
			Group: 2,
		},
		term{
			ID:    5,
			Name:  "name1-2",
			Slug:  "name1-2",
			Group: 1,
		},
		term{
			ID:    6,
			Name:  "name1-2-1",
			Slug:  "name1-2-1",
			Group: 5,
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
			Group: 7,
		},
	}, terms)
}
