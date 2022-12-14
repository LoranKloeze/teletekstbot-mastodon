// Copyright 2022 Codedivision - Loran Kloeze. All rights reserved.
// Use of this source code is governed by the MIT-license.

package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
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
	pageNr, err := teletekst.ConstructPageNr(notification.Status.Content)
	if err != nil {
		return teletekst.Page{}, false
	}
	return teletekst.Page{Nr: pageNr}, true
}

func main() {
	logStart()

	store := teletekst.InitStore("teletekst")
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

			currentId, err := strconv.Atoi(string(n.ID))
			if err != nil {
				log.Fatal("notification ID should always be convertable to integer")
			}

			lastId, err := strconv.Atoi(string(teletekst.LastNotificationId()))
			if err != nil || currentId > lastId {
				teletekst.SetLastNotificationId(n.ID)
			}

			fmt.Printf("Page %s asked by %s\n", p.Nr, n.Account.Acct)
			if !teletekst.NOSHasPage(p.Nr) {
				teletekst.Post404ReplyToot(p, n)
				continue
			}

			if os.Getenv("TELETEKST_ENV") == "production" {
				teletekst.PersistScreenshotReply(p)
				teletekst.PostReplyToot(p, n)
			} else {
				teletekst.FakeReplyToot(p, n)
			}

		}
		time.Sleep(15 * time.Second)
	}
}
