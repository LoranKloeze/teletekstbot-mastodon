// Copyright 2022 Codedivision - Loran Kloeze. All rights reserved.
// Use of this source code is governed by the MIT-license.

package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"loran.cc/teletekstbot/teletekst"
)

func logStart() {
	fmt.Println("Starting teletekst bot for 101's")
	if os.Getenv("TELETEKST_ENV") == "production" {
		fmt.Printf("Mode: %s\n", "production")
	} else {
		fmt.Printf("Mode: %s\n", "development")
	}
}

func main() {
	firstPage, lastPage := 104, 150
	logStart()

	store := teletekst.InitStore("teletekst")
	defer store.Close()

	teletekst.InitMastodon()

	for i := firstPage; i < lastPage+1; i++ {

		// Let's not overflow the NOS servers
		time.Sleep(500 * time.Millisecond)

		page := teletekst.DownloadPage(strconv.Itoa(i), "https://teletekst-data.nos.nl")

		// The content of a page is empty if NOS told us there is no page
		if page.Content == "" || teletekst.PageExists(page) {
			log.Printf("Skipping %s\n", page.Nr)
			continue
		}

		log.Printf("--- Constructing a toot for %s --- \n", page.Nr)
		if os.Getenv("TELETEKST_ENV") == "production" {
			teletekst.PersistScreenshot101(page)
			teletekst.Post101Toot(page)
		} else {
			teletekst.Fake101Toot(page)
		}
		teletekst.InsertPage(page)
	}

}
