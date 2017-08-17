package integration

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/imega-teleport/db2file/mysql"
	"github.com/imega-teleport/xml2db/commerceml"
	"github.com/stretchr/testify/assert"
)

type dbunit struct {
	db *sql.DB
	t  *testing.T
}

func (u *dbunit) setup(t *testing.T, tableName []string, fixture func(db *sql.DB) (err error)) (db *sql.DB, teardown func()) {
	user, pass, host, dbname := os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_HOST"), os.Getenv("DB_NAME")
	dsn := fmt.Sprintf("mysql://%s:%s@tcp(%s)/%s", user, pass, host, dbname)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		t.Fatalf("Could not open db: %s", err)
	}
	u.db = db

	for _, n := range tableName {
		if _, err := db.Exec(fmt.Sprintf("TRUNCATE %s", n)); err != nil {
			t.Fatalf("Could not truncate table %s: %s", n, err)
		}
	}

	if err := fixture(db); err != nil {
		t.Fatalf("Could not load fixtures: %s", err)
	}

	teardown = func() {
		for _, n := range tableName {
			if _, err := db.Exec(fmt.Sprintf("TRUNCATE %s", n)); err != nil {
				t.Fatalf("Could not truncate table %s: %s", n, err)
			}
		}

		if err := db.Close(); err != nil {
			t.Fatalf("Could not close db: %s", err)
		}
	}
	return
}

var dbUnit = &dbunit{}

func Test_Groups_ReturnsGroups(t *testing.T) {
	db, teardown := dbUnit.setup(t, []string{"groups"}, func(db *sql.DB) (err error) {
		_, err = db.Query("INSERT groups VALUES (?,?,?)", "ecc82696-3f98-11de-991a-001c23888998", "", "Group 1")
		return
	})
	defer teardown()

	s := mysql.NewStorage(db)
	groups, err := s.Groups("")
	assert.NoError(t, err)
	expected := []commerceml.Group{
		{
			IdName: commerceml.IdName{
				Id:   "ecc82696-3f98-11de-991a-001c23888998",
				Name: "Group 1",
			},
			Groups: []commerceml.Group{},
		},
	}

	assert.Equal(t, expected, groups)
}

func Test_Groups_NotExistsGroup_ReturnsEmptyGroups(t *testing.T) {
	db, teardown := dbUnit.setup(t, []string{"groups"}, func(db *sql.DB) (err error) {
		return
	})
	defer teardown()

	s := mysql.NewStorage(db)
	groups, err := s.Groups("")
	assert.NoError(t, err)
	expected := []commerceml.Group{}

	assert.Equal(t, expected, groups)
}

func Test_Groups_WithChildGroup_ReturnsGroups(t *testing.T) {
	db, teardown := dbUnit.setup(t, []string{"groups"}, func(db *sql.DB) (err error) {
		_, err = db.Query("INSERT groups VALUES (?,?,?)", "ecc82696-3f98-11de-991a-001c23888998", "", "Group 1")
		if err != nil {
			return
		}
		_, err = db.Query("INSERT groups VALUES (?,?,?)", "7077e5f0-f2a5-11de-bc7e-0022b0527b2e", "ecc82696-3f98-11de-991a-001c23888998", "Child Group 1")
		return
	})
	defer teardown()

	s := mysql.NewStorage(db)
	groups, err := s.Groups("")
	assert.NoError(t, err)
	expected := []commerceml.Group{
		{
			IdName: commerceml.IdName{
				Id:   "ecc82696-3f98-11de-991a-001c23888998",
				Name: "Group 1",
			},
			Groups: []commerceml.Group{
				{
					IdName: commerceml.IdName{
						Id:   "7077e5f0-f2a5-11de-bc7e-0022b0527b2e",
						Name: "Child Group 1",
					},
					Groups: []commerceml.Group{},
				},
			},
		},
	}

	assert.Equal(t, expected, groups)
}

func Test_Products_ReturnsProducts(t *testing.T) {
	db, teardown := dbUnit.setup(t, []string{"products", "products_groups"}, func(db *sql.DB) (err error) {
		_, err = db.Query(
			"INSERT products VALUES (?,?,?,?,?,?,?,?)",
			"b9f7eba5-ae8b-11e3-8162-003048f2904a",
			"prod1_name",
			"prod1_description",
			"prod1_barcode",
			"prod1_article",
			"prod1_fullname",
			"prod1_country",
			"prod1_brand",
		)
		if err != nil {
			return
		}
		_, err = db.Query(
			"INSERT products VALUES (?,?,?,?,?,?,?,?)",
			"7077e5f0-f2a5-11de-bc7e-0022b0527b2e",
			"prod2_name",
			"prod2_description",
			"prod2_barcode",
			"prod2_article",
			"prod2_fullname",
			"prod2_country",
			"prod2_brand",
		)
		if err != nil {
			return
		}
		_, err = db.Query(
			"INSERT products_groups VALUES (?,?)",
			"b9f7eba5-ae8b-11e3-8162-003048f2904a",
			"de258629-6b29-11e4-8220-005056b9f84b",
		)
		if err != nil {
			return
		}
		_, err = db.Query(
			"INSERT products_groups VALUES (?,?)",
			"7077e5f0-f2a5-11de-bc7e-0022b0527b2e",
			"cf0c4f35-b32c-11e3-8162-003048f2904a",
		)
		return
	})
	defer teardown()

	s := mysql.NewStorage(db)
	products, err := s.Products()
	assert.NoError(t, err)

	expected := []commerceml.Product{
		{
			IdName: commerceml.IdName{
				Id:   "b9f7eba5-ae8b-11e3-8162-003048f2904a",
				Name: "prod1_name",
			},
			Description: commerceml.Description{
				Value: "prod1_description",
			},
			BarCode:  "prod1_barcode",
			Article:  "prod1_article",
			FullName: "prod1_fullname",
			Country:  "prod1_country",
			Brand:    "prod1_brand",
			Groups: []commerceml.Group{
				{
					IdName: commerceml.IdName{
						Id: "de258629-6b29-11e4-8220-005056b9f84b",
					},
				},
			},
		},
		{
			IdName: commerceml.IdName{
				Id:   "7077e5f0-f2a5-11de-bc7e-0022b0527b2e",
				Name: "prod2_name",
			},
			Description: commerceml.Description{
				Value: "prod2_description",
			},
			BarCode:  "prod2_barcode",
			Article:  "prod2_article",
			FullName: "prod2_fullname",
			Country:  "prod2_country",
			Brand:    "prod2_brand",
			Groups: []commerceml.Group{
				{
					IdName: commerceml.IdName{
						Id: "cf0c4f35-b32c-11e3-8162-003048f2904a",
					},
				},
			},
		},
	}
	assert.Equal(t, expected, products)
}

func Test_CheckCompleteAllTasks_ReturnsTrue(t *testing.T) {
	db, teardown := dbUnit.setup(t, []string{"tasks"}, func(db *sql.DB) (err error) {
		_, err = db.Query(
			"INSERT tasks VALUES (?,?)",
			"store",
			1,
		)
		if err != nil {
			return
		}
		_, err = db.Query(
			"INSERT tasks VALUES (?,?)",
			"offer",
			1,
		)
		if err != nil {
			return
		}
		return
	})
	defer teardown()

	s := mysql.NewStorage(db)
	complete, err := s.CheckCompleteAllTasks()
	assert.NoError(t, err)

	assert.True(t, complete)
}
