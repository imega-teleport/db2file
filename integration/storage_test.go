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

func (u *dbunit) setup(t *testing.T, tableName string, fixture func(db *sql.DB) (err error)) (db *sql.DB, teardown func()) {
	user, pass, host, dbname := os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_HOST"), os.Getenv("DB_NAME")
	dsn := fmt.Sprintf("mysql://%s:%s@tcp(%s)/%s", user, pass, host, dbname)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		t.Fatalf("Could not open db: %s", err)
	}
	u.db = db

	if _, err := db.Exec(fmt.Sprintf("TRUNCATE %s", tableName)); err != nil {
		t.Fatalf("Could not truncate table %s: %s", tableName, err)
	}

	if err := fixture(db); err != nil {
		t.Fatalf("Could not load fixtures: %s", err)
	}

	teardown = func() {
		if _, err := db.Exec(fmt.Sprintf("TRUNCATE %s", tableName)); err != nil {
			t.Fatalf("Could not truncate table %s: %s", tableName, err)
		}

		if err := db.Close(); err != nil {
			t.Fatalf("Could not close db: %s", err)
		}
	}
	return
}

var dbUnit = &dbunit{}

func Test_Groups_ReturnsGroups(t *testing.T) {
	db, teardown := dbUnit.setup(t, "groups", func(db *sql.DB) (err error) {
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
	db, teardown := dbUnit.setup(t, "groups", func(db *sql.DB) (err error) {
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
	db, teardown := dbUnit.setup(t, "groups", func(db *sql.DB) (err error) {
		_, err = db.Query("INSERT groups VALUES (?,?,?)", "ecc82696-3f98-11de-991a-001c23888998", "", "Group 1")
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
