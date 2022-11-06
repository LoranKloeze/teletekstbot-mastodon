// Copyright 2022 Codedivision - Loran Kloeze. All rights reserved.
// Use of this source code is governed by the MIT-license.

package main

import (
	"encoding/json"
	"fmt"
	"html"
	"net/http"
	"os"
	"regexp"
	"time"
)

type Page struct {
	Content string `json:"content"`
	Nr      string
	Hash    string
	Title   string
}

// DownloadPage downloads the textual representation of a teletekst page from the NOS
// using the provided server e.g. https://teletekst-data.nos.nl
//
// It prevents server side caching by using a timestamp in the query
func DownloadPage(pageNr string, server string) (p Page) {
	p.Nr = pageNr
	u := fmt.Sprintf("%s/json/%s?t=%d", server, pageNr, time.Now().UnixNano())
	r, err := http.Get(u)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not download teletekst page: %s", err)
	}
	defer r.Body.Close()

	json.NewDecoder(r.Body).Decode(&p)
	p.Hash = MD5Hash(p.Content)
	p.Title = extractTitle(p)
	return
}

func extractTitle(p Page) string {
	re := regexp.MustCompile(`<span class=\"yellow bg-blue doubleHeight \">(.+?)</span>`)
	res := re.FindAllStringSubmatch(p.Content, -1)
	if len(res) == 0 || len(res[0]) == 0 {
		return "Onbekende titel"
	}

	return html.UnescapeString(removeTags(res[0][1]))
}

func removeTags(s string) string {
	re := regexp.MustCompile(`<.*?>`)
	return re.ReplaceAllString(s, "")
}