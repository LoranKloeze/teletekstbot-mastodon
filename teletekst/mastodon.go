// Copyright 2022 Codedivision - Loran Kloeze. All rights reserved.
// Use of this source code is governed by the MIT-license.

package teletekst

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/mattn/go-mastodon"
	"golang.org/x/net/html"
)

var cli *mastodon.Client

func InitMastodon() {
	godotenv.Load()
	cli = mastodon.NewClient(&mastodon.Config{
		Server:       os.Getenv("MASTODON_SERVER"),
		ClientID:     os.Getenv("MASTODON_CLIENT_ID"),
		ClientSecret: os.Getenv("MASTODON_CLIENT_SECRET"),
		AccessToken:  os.Getenv("MASTODON_ACCESS_TOKEN"),
	})
}

// Fake101Toot logs a toot that would've been posted
//
// Mainly here for debugging purposes
func Fake101Toot(p Page) {
	log.Printf(">>> Would've posted a 101 toot for %s with title '%s'... ", p.Nr, p.Title)
}

// FakeReplyToot logs a toot that would've been posted
//
// Mainly here for debugging purposes
func FakeReplyToot(p Page, n *mastodon.Notification) {
	log.Printf(">>> Would've posted a reply toot for %s for %s... ", p.Nr, n.Account.Acct)
}

func Post101Toot(p Page) {
	log.Printf(">>> Posting a 101 toot for %s... ", p.Nr)
	ctx := context.Background()

	// To send a toot with an attachment, we first need to upload that attachment
	att := uploadScreenshot(ctx, p, "101")

	url := "https://nos.nl/teletekst#" + p.Nr
	text := fmt.Sprintf("[%s] %s\n%s", p.Nr, p.Title, url)
	_, err := cli.PostStatus(ctx, &mastodon.Toot{Status: text, Visibility: "unlisted", MediaIDs: []mastodon.ID{att.ID}})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Done!\n")
}

func Post404ReplyToot(p Page, n *mastodon.Notification) {
	log.Printf(">>> Posting a 404 toot for %s... ", p.Nr)
	ctx := context.Background()

	text := fmt.Sprintf("@%s Je vroeg om pagina %s maar die kon ik niet vinden...", n.Account.Acct, p.Nr)
	_, err := cli.PostStatus(ctx, &mastodon.Toot{Status: text, Visibility: "unlisted", InReplyToID: n.Status.ID})
	if err != nil {
		log.Fatal(err)
	}
}

func PostReplyToot(p Page, n *mastodon.Notification) {
	log.Printf(">>> Posting a reply toot for %s... ", p.Nr)
	ctx := context.Background()

	// To send a toot with an attachment, we first need to upload that attachment
	att := uploadScreenshot(ctx, p, "reply")

	text := fmt.Sprintf("@%s Je vroeg om pagina %s, hierbij.", n.Account.Acct, p.Nr)
	_, err := cli.PostStatus(ctx, &mastodon.Toot{Status: text, Visibility: "unlisted", InReplyToID: n.Status.ID, MediaIDs: []mastodon.ID{att.ID}})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Done!\n")
}

func NewNotifications() []*mastodon.Notification {
	ctx := context.Background()
	ns, err := cli.GetNotifications(ctx, &mastodon.Pagination{SinceID: LastNotificationId()})
	if err != nil {
		log.Fatal(err)
	}
	return ns
}

func uploadScreenshot(ctx context.Context, p Page, prefix string) *mastodon.Attachment {
	path := fmt.Sprintf("/tmp/gowitness/screenshots/%s_%s_cropped.png", prefix, p.Nr)
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	media := mastodon.Media{
		File:        f,
		Focus:       "0.0,1,0",
		Description: nosHtmlToText(p),
	}
	att, err := cli.UploadMediaFromMedia(ctx, &media)

	if err != nil {
		log.Fatal(err)
	}

	return att
}

func nosHtmlToText(p Page) string {
	var bld strings.Builder
	b := bytes.NewBufferString(p.Content)

	doc, err := html.Parse(b)
	if err != nil {
		log.Fatal(err)
	}

	var rowCnt int
	var fn func(*html.Node)
	fn = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "span" {
			txt := n.FirstChild.Data
			if rowCnt > 4 && !strings.HasPrefix(txt, "") && txt != "a" {
				bld.WriteString(txt)
				bld.WriteString("\n")
			}
			rowCnt++
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			fn(c)
		}
	}
	fn(doc)

	return bld.String()
}
