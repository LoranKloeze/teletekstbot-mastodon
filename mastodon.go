// Copyright 2022 Codedivision - Loran Kloeze. All rights reserved.
// Use of this source code is governed by the MIT-license.

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"regexp"

	"github.com/joho/godotenv"
	"github.com/mattn/go-mastodon"
)

var cli *mastodon.Client

func init() {
	godotenv.Load()
	cli = mastodon.NewClient(&mastodon.Config{
		Server:       os.Getenv("MASTODON_SERVER"),
		ClientID:     os.Getenv("MASTODON_CLIENT_ID"),
		ClientSecret: os.Getenv("MASTODON_CLIENT_SECRET"),
		AccessToken:  os.Getenv("MASTODON_ACCESS_TOKEN"),
	})
}

// FakeToot logs the toot that would've been posted
//
// Mainly here for debugging purposes
func FakeToot(p Page) {
	log.Printf(">>> Would've posted a toot for %s with title '%s'... ", p.Nr, p.Title)
}

func PostToot(p Page) {
	log.Printf(">>> Posting a toot for %s... ", p.Nr)
	ctx := context.Background()

	// To send a toot with an attachment, we first need to upload that attachment
	att := uploadScreenshot(ctx, p)

	url := "https://nos.nl/teletekst#" + p.Nr
	text := fmt.Sprintf("[%s] %s\n%s", p.Nr, p.Title, url)
	_, err := cli.PostStatus(ctx, &mastodon.Toot{Status: text, Visibility: "unlisted", MediaIDs: []mastodon.ID{att.ID}})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Done!\n")
}

func Notifications() {
	ctx := context.Background()
	ns, err := cli.GetNotifications(ctx, &mastodon.Pagination{SinceID: LastNotificationId()})
	if err != nil {
		log.Fatal(err)
	}

	for _, n := range ns {
		if n.Type != "mention" {
			continue
		}
		re := regexp.MustCompile(`Pagina\s(\d{3})`)
		m := re.FindAllStringSubmatch(n.Status.Content, 1)
		if len(m) > 0 && len(m[0]) > 0 {
			p := Page{Nr: m[0][1]}
			fmt.Printf("Page %s asked\n", p.Nr)
			createScreenshot(p)
			cropScreenshot(p)
			att := uploadScreenshot(ctx, p)
			text := fmt.Sprintf("@%s Je vroeg om pagina %s, hierbij.", n.Account.Acct, p.Nr)
			cli.PostStatus(ctx, &mastodon.Toot{Status: text, Visibility: "unlisted", InReplyToID: n.Status.ID, MediaIDs: []mastodon.ID{att.ID}})
			SetLastNotificationId(n.ID)
		}
	}

}

func uploadScreenshot(ctx context.Context, p Page) *mastodon.Attachment {
	f, err := os.Open("/tmp/gowitness/screenshots/" + p.Nr + "_cropped.png")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	media := mastodon.Media{
		File:  f,
		Focus: "0.0,1,0",
	}
	att, err := cli.UploadMediaFromMedia(ctx, &media)

	if err != nil {
		log.Fatal(err)
	}

	return att
}
