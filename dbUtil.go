package util

import (
	"database/sql"
	"errors"
	"strings"
)

type DbUtil struct{}

func (dbUtil *DbUtil) GetDb(user, password, host, database string, flagsMap map[string]string) (*sql.DB, error) {
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
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

// TODO: where and params
func (dbUtil *DbUtil) ClearTables(db *sql.DB, tables ...string) error {
	prefix := "dbUtil.ClearTables"
	tx, err := db.Begin()
	if err != nil {
		return errors.New(prefix + " (tx begin): " + err.Error())
	}

	for _, table := range tables {
		_, err := db.Exec("DELETE FROM " + table)
		if err != nil {
			err2 := tx.Rollback()
			if err2 != nil {
				return errors.New(prefix + " (tx rollback at " + table + "): " + err.Error())
			}

			return errors.New(prefix + " (db exec at " + table + "): " + err.Error())
		}
	}

	err = tx.Commit()
	if err != nil {
		return errors.New(prefix + " (tx commit): " + err.Error())
	}

	return nil
}

func (dbUtil *DbUtil) GetRowCount(db *sql.DB, table string, where string, params []interface{}) (int, error) {
	prefix := "dbUtil.AssertRowCount"
	query := "SELECT COUNT(*) FROM " + table
	if where != "" {
		query += " WHERE " + where
	}

	rows, err := db.Query(query, params...)
	if err != nil {
		return 1, errors.New(prefix + " (query count): " + err.Error())
	}
	defer rows.Close()

	var count int
	if rows.Next() {
		rows.Scan(&count)
	}

	return count, nil
}
