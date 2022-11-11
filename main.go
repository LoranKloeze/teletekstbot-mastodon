// Copyright 2022 Codedivision - Loran Kloeze. All rights reserved.
// Use of this source code is governed by the MIT-license.

package main

import (
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func logStart() {
	fmt.Println("Starting teletekst bot")
	if os.Getenv("TELETEKST_ENV") == "production" {
		fmt.Printf("Mode: %s\n", "production")
	} else {
		fmt.Printf("Mode: %s\n", "development")
	}
}

func main() {
	logStart()
	Notifications()
}
