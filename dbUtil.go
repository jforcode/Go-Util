package util

import (
	"database/sql"
	"strings"

	"github.com/jforcode/DeepError"
)

type IDbUtil interface {
	GetDb(user, password, host, database string, flagsMap map[string]string) (*sql.DB, error)
	ClearTables(db *sql.DB, tables ...string) error
	GetRowCount(db *sql.DB, table string, where string, params []interface{}) (int, error)
	PrepareAndExec(db *sql.DB, query string, parameters ...interface{}) (int64, error)
}

type DbUtil struct{}

func (dbUtil *DbUtil) GetDb(user, password, host, database string, flagsMap map[string]string) (*sql.DB, error) {
	fnName := "util.DbUtil.GetDb"

	flags := make([]string, len(flagsMap))
	index := 0
	for key, value := range flagsMap {
		flags[index] = key + "=" + value
		index++
	}

	var flagsS string
	if len(flags) == 0 {
		flagsS = ""
	} else {
		flagsS = "?" + strings.Join(flags, "&")
	}

	datasource := user +
		":" + password +
		"@" + host +
		"/" + database +
		flagsS

	db, err := sql.Open("mysql", datasource)
	if err != nil {
		return nil, deepError.New(fnName, "opening db", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, deepError.New(fnName, "pinging db", err)
	}

	return db, nil
}

// TODO: where and params
func (dbUtil *DbUtil) ClearTables(db *sql.DB, tables ...string) error {
	fnName := "util.DbUtil.ClearTables"

	tx, err := db.Begin()
	if err != nil {
		return deepError.New(fnName, "tx begin", err)
	}

	for _, table := range tables {
		_, err := db.Exec("DELETE FROM " + table)
		if err != nil {
			err2 := tx.Rollback()
			if err2 != nil {
				return deepError.DeepErr{
					Function: fnName,
					Action:   "rolling back transaction",
					Message:  "Couldn't roll back transaction for table",
					Params:   []interface{}{table},
					Cause:    err2,
				}
			}

			return deepError.DeepErr{
				Function: fnName,
				Action:   "exec",
				Message:  "Couldn't exec for table",
				Params:   []interface{}{table},
				Cause:    err,
			}
		}
	}

	err = tx.Commit()
	if err != nil {
		return deepError.New(fnName, "commiting transaction", err)
	}

	return nil
}

func (dbUtil *DbUtil) GetRowCount(db *sql.DB, table string, where string, params []interface{}) (int, error) {
	fnName := "util.dbUtil.AssertRowCount"

	query := "SELECT COUNT(*) FROM " + table
	if where != "" {
		query += " WHERE " + where
	}

	rows, err := db.Query(query, params...)
	if err != nil {
		return -1, deepError.New(fnName, "querying", err)
	}
	defer rows.Close()

	var count int
	if rows.Next() {
		err = rows.Scan(&count)

		if err != nil {
			return -1, deepError.New(fnName, "scanning", err)
		}
	}

	return count, nil
}

func (dbUtil *DbUtil) PrepareAndExec(db *sql.DB, query string, parameters ...interface{}) (int64, error) {
	fnName := "util.dbUtil.PrepareAndExec"

	stmt, err := db.Prepare(query)
	if err != nil {
		return -1, deepError.New(fnName, "preparing", err)
	}

	res, err := stmt.Exec(parameters...)
	if err != nil {
		return -1, deepError.New(fnName, "exec", err)
	}

	numInserted, err := res.RowsAffected()
	if err != nil {
		return -1, deepError.New(fnName, "getting rows affected", err)
	}

	return numInserted, nil
}
