package small

import (
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/nalgeon/redka"
)

func TestRedDB(t *testing.T) {
	db := NewRedDB()

	db.Str().Set("name", "alice")
	db.Str().Set("age", 25)

	count, err := db.Key().Count("name", "age", "city")
	t.Log("count", "count", count, "err", err)
	name, err := db.Str().Get("name")
	t.Log("get", "name", name, "err", err)

	updCount := 0
	err = db.Update(func(tx *redka.Tx) error {
		err := tx.Str().Set("name", "bob")
		if err != nil {
			return err
		}
		updCount++
		err = tx.Str().Set("age", 50)
		if err != nil {
			return err
		}
		updCount++
		return nil
	})

	t.Log("updated", "count", updCount, "err", err)
	defer db.Close()
}
