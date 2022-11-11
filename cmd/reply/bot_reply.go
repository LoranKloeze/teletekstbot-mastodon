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

func constructPage(notification *mastodon.Notification) (page teletekst.Page, ok bool) {
	if notification.Type != "mention" {
		return teletekst.Page{}, false
	}
	re := regexp.MustCompile(`(?i)pagina\s(\d{3})`)
	m := re.FindAllStringSubmatch(notification.Status.Content, 1)
	if len(m) > 0 && len(m[0]) > 0 {
		return teletekst.Page{Nr: m[0][1]}, true
	} else {
		return teletekst.Page{}, false
	}
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
