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

	// Table 'pages'
	sqlStmtPages := `
	CREATE TABLE IF NOT EXISTS pages (nr text unique, title text);
	
	`
	_, err = db.Exec(sqlStmtPages)
	if err != nil {
		log.Fatalf("%q: %s\n", err, sqlStmtPages)
		return nil
	}

	// Table 'key_vals'
	sqlStmt := `
	CREATE TABLE IF NOT EXISTS key_vals (key text unique, value text);
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return nil
	}
	return db
}

func PageExists(db *sql.DB, p Page) bool {
	var title string

	stmt, _ := db.Prepare("SELECT title FROM pages WHERE nr = ?")
	defer stmt.Close()

	err := stmt.QueryRow(p.Nr).Scan(&title)
	if err != nil {
		return false
	}

	return title == p.Title
}

func InsertPage(db *sql.DB, p Page) {
	stmt, _ := db.Prepare("INSERT OR REPLACE INTO pages (nr, title) values(?, ?)")
	defer stmt.Close()

	stmt.Exec(p.Nr, p.Title)
}
