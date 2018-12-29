package util

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/google/go-cmp/cmp"
	"github.com/jforcode/Go-DeepError"
)

type IDbUtil interface {
	GetDb(user, password, host, database string, flagsMap map[string]string) (*sql.DB, error)
	ClearTables(db *sql.DB, tables ...string) error
	GetRowCount(db *sql.DB, table string, where string, params []interface{}) (int, error)
	PrepareAndExec(db *sql.DB, query string, parameters ...interface{}) (sql.Result, error)
	CompareDbData(db *sql.DB, query string, args []interface{}, numCols int, expected [][]string, debug bool) (bool, error)
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

	for _, table := range tables {
		_, err := db.Exec("DELETE FROM " + table)
		if err != nil {
			return deepError.DeepErr{
				Function: fnName,
				Action:   "exec",
				Message:  "Couldn't exec for table",
				Params:   []interface{}{table},
				Cause:    err,
			}
		}
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

func (dbUtil *DbUtil) PrepareAndExec(db *sql.DB, query string, parameters ...interface{}) (sql.Result, error) {
	fnName := "util.dbUtil.PrepareAndExec"

	stmt, err := db.Prepare(query)
	if err != nil {
		return nil, deepError.New(fnName, "preparing", err)
	}

	res, err := stmt.Exec(parameters...)
	if err != nil {
		return nil, deepError.New(fnName, "exec", err)
	}

	return res, nil
}

func (dbUtil *DbUtil) CompareDbData(db *sql.DB, query string, args []interface{}, numCols int, expected [][]string, debug bool) (bool, error) {
	fn := "CompareDbData"
	rows, err := db.Query(query, args...)
	if err != nil {
		return false, deepError.New(fn, "query", err)
	}
	defer rows.Close()

	actual := make([][]string, 0)
	for rows.Next() {
		rowData := make([]string, numCols)
		rowDataPtr := make([]interface{}, numCols)
		for i := 0; i < numCols; i++ {
			rowDataPtr[i] = &rowData[i]
		}

		err := rows.Scan(rowDataPtr...)
		if err != nil {
			return false, deepError.New(fn, "scan", err)
		}

		actual = append(actual, rowData)
	}

	if debug {
		fmt.Printf("\nActual\n%+v\nExpected\n%+v\n", actual, expected)
	}

	if len(actual) != len(expected) {
		return false, nil
	}

	for ind, act := range actual {
		if !cmp.Equal(act, expected[ind]) {
			return false, nil
		}
	}

	return true, nil
}
