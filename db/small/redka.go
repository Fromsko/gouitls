package small

import (
	"log/slog"

	_ "github.com/mattn/go-sqlite3"
	"github.com/nalgeon/redka"
)

var (
	err         error
	RED_DB      *redka.DB
	CONNECT_DNS string
)

func checkParams() (have bool) {
	have = CONNECT_DNS == ""
	if !have {
		CONNECT_DNS = "_app.db"
	}
	return have
}

func init() {
	checkParams()
}

func NewRedDB() *redka.DB {
	RED_DB, err = redka.Open(CONNECT_DNS, nil)
	if err != nil {
		slog.Error("Connect RED_DB: ", "err", err)
	}
	return RED_DB
}
