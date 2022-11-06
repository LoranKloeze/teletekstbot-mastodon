// Copyright 2022 Codedivision - Loran Kloeze. All rights reserved.
// Use of this source code is governed by the MIT-license.

package main

import (
	"database/sql"
	"log"
)

func InitDb(file string) *sql.DB {
	db, err := sql.Open("sqlite3", file)
	if err != nil {
		log.Fatal(err)
	}

	sqlStmt := `
	CREATE TABLE IF NOT EXISTS hashes (nr text unique, hash text);
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return nil
	}
	return db
}

func PageExists(db *sql.DB, p Page) bool {
	var hash string

	stmt, _ := db.Prepare("SELECT hash FROM hashes WHERE nr = ?")
	defer stmt.Close()

	err := stmt.QueryRow(p.Nr).Scan(&hash)
	if err != nil {
		return false
	}

	return hash == p.Hash
}

func InsertPage(db *sql.DB, p Page) {
	stmt, _ := db.Prepare("INSERT OR REPLACE INTO HASHES (nr, hash) values(?, ?)")
	defer stmt.Close()

	stmt.Exec(p.Nr, p.Hash)
}
