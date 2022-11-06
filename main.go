// Copyright 2022 Codedivision - Loran Kloeze. All rights reserved.
// Use of this source code is governed by the MIT-license.

package main

import (
	"log"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	firstPage, lastPage := 104, 150

	db := InitDb("teletekstbot.db")
	defer db.Close()

	for i := firstPage; i < lastPage+1; i++ {

		// Let's not overflow the NOS servers
		time.Sleep(500 * time.Millisecond)

		page := DownloadPage(strconv.Itoa(i), "https://teletekst-data.nos.nl")

		// The content of a page is empty if NOS told us there is no page
		if page.Content == "" || PageExists(db, page) {
			log.Printf("Skipping %s\n", page.Nr)
			continue
		}

		log.Printf("--- Constructing a toot for %s --- \n", page.Nr)
		PersistScreenshot(page)
		PostToot(page)
		InsertPage(db, page)
	}

}
