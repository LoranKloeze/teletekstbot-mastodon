// Copyright 2022 Codedivision - Loran Kloeze. All rights reserved.
// Use of this source code is governed by the MIT-license.

package main

import (
	"fmt"
	"os"
	"testing"
)

var (
	tmpDir       = "tmp/"
	databaseFile = "test_db.db"
)

func TestDBInit(t *testing.T) {
	removeDBFile()
	db := InitDb(tmpDir + databaseFile)
	defer db.Close()

	inf, err := os.Stat(tmpDir + databaseFile)
	if err != nil {
		t.Errorf("Got error statting file: %s", err)
	}
	if inf.Name() != databaseFile {
		t.Errorf("Expected a file named %s, got %s", databaseFile, inf.Name())
	}

}

func TestInsertPage(t *testing.T) {
	removeDBFile()
	db := InitDb(tmpDir + databaseFile)
	defer db.Close()

	InsertPage(db, Page{"content", "199", "h1", ""})

	stmt, _ := db.Prepare("SELECT hash FROM hashes WHERE nr = ?")
	defer stmt.Close()

	var val string

	err := stmt.QueryRow("199").Scan(&val)
	if err != nil {
		t.Errorf("Got error searching page: %s", err)
	}

	if val != "h1" {
		t.Errorf("Expected page hash h1, got %s", val)
	}
	fmt.Println(val)
}

func TestPageExists(t *testing.T) {
	removeDBFile()
	db := InitDb(tmpDir + databaseFile)
	defer db.Close()

	stmt, _ := db.Prepare("INSERT OR REPLACE INTO HASHES (nr, hash) values(?, ?)")
	defer stmt.Close()

	stmt.Exec("150", "h10")

	page := Page{"content", "150", "h10", ""}
	if !PageExists(db, page) {
		t.Errorf("Expected page to exist: %s", page)
	}
}

func removeDBFile() {
	os.Remove(tmpDir + databaseFile)
}
