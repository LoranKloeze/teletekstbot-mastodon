// Copyright 2022 Codedivision - Loran Kloeze. All rights reserved.
// Use of this source code is governed by the MIT-license.

package teletekst

import (
	"encoding/json"
	"errors"
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
func DownloadPage(pageNr string, server string) (p Page, err error) {
	p.Nr = pageNr
	u := fmt.Sprintf("%s/json/%s?t=%d", server, pageNr, time.Now().UnixNano())
	r, err := http.Get(u)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not download teletekst page: %s", err)
	}
	defer r.Body.Close()
	if r.StatusCode == http.StatusNotFound {
		return Page{}, fmt.Errorf("page %s not found at NOS", pageNr)
	}

	json.NewDecoder(r.Body).Decode(&p)
	p.Hash = MD5Hash(p.Content)
	p.Title, err = extractTitle(p)
	if err != nil {
		return p, err
	}
	return p, nil
}

func ConstructPageNr(content string) (string, error) {
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

	return "", errors.New("could not find a page number in content")
}

func NOSHasPage(nr string) bool {
	_, err := DownloadPage(nr, "https://teletekst-data.nos.nl")
	return err == nil
}

func extractTitle(p Page) (string, error) {
	re := regexp.MustCompile(`<span class=\"yellow bg-blue doubleHeight \">(.+?)</span>`)
	res := re.FindAllStringSubmatch(p.Content, -1)
	if len(res) == 0 || len(res[0]) == 0 {
		return "", errors.New("could not extract a title from html")
	}

	return html.UnescapeString(removeTags(res[0][1])), nil
}

func removeTags(s string) string {
	re := regexp.MustCompile(`<.*?>`)
	return re.ReplaceAllString(s, "")
}
