// Copyright 2022 Codedivision - Loran Kloeze. All rights reserved.
// Use of this source code is governed by the MIT-license.

package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"time"

	"github.com/mattn/go-mastodon"
	"loran.cc/teletekstbot/teletekst"
)

func logStart() {
	fmt.Println("Starting teletekst bot for replies")
	if os.Getenv("TELETEKST_ENV") == "production" {
		fmt.Printf("Mode: %s\n", "production")
	} else {
		fmt.Printf("Mode: %s\n", "development")
	}
}

func pageNumber(content string) (string, error) {
	var re *regexp.Regexp
	var m [][]string

	// Check for xxx-xx like queries
	re = regexp.MustCompile(`(?i)pagina\s(\d{3}-\d{1,2})`)
	m = re.FindAllStringSubmatch(content, 1)
	if len(m) > 0 {
		return m[0][1], nil
	}

	// Check for xxx like queries
	re = regexp.MustCompile(`(?i)pagina\s(\d{3})`)
	m = re.FindAllStringSubmatch(content, 1)
	if len(m) > 0 {
		return m[0][1], nil
	}

	return "", nil
}

func constructPage(notification *mastodon.Notification) (page teletekst.Page, ok bool) {
	if notification.Type != "mention" {
		return teletekst.Page{}, false
	}
	pageNr, err := pageNumber(notification.Status.Content)
	if err != nil {
		return teletekst.Page{}, false
	}
	return teletekst.Page{Nr: pageNr}, false
}

func main() {
	logStart()

	store := teletekst.InitStore()
	defer store.Close()

	teletekst.InitMastodon()
	for {
		log.Printf("Checking notifications...")
		ns := teletekst.NewNotifications()

		for _, n := range ns {
			if n.Type != "mention" {
				continue
			}

			p, ok := constructPage(n)
			if !ok {
				continue
			}

			fmt.Printf("Page %s asked by %s\n", p.Nr, n.Account.Acct)
			if os.Getenv("TELETEKST_ENV") == "production" {
				teletekst.PersistScreenshotReply(p)
				teletekst.PostReplyToot(p, n)
			} else {
				teletekst.FakeReplyToot(p, n)
			}
			teletekst.SetLastNotificationId(n.ID)
		}
		time.Sleep(15 * time.Second)
	}
}
