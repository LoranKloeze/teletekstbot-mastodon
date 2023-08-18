// Copyright 2022 Codedivision - Loran Kloeze. All rights reserved.
// Use of this source code is governed by the MIT-license.

package teletekst

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDownloadPage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/json/110" {
			t.Errorf("Expected to request '/json', got: %s", r.URL.Path)
		}

		if !strings.HasPrefix(r.URL.RawQuery, "t=") {
			t.Errorf("Expected query to have t=, got: %s", r.URL.RawQuery)
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"content":"sample_content"}`))
	}))
	defer server.Close()

	value, _ := DownloadPage("110", server.URL)

	if value.Content != "sample_content" {
		t.Errorf("Expected 'sample_content' as page content, got %s", value.Content)
	}
	if value.Nr != "110" {
		t.Errorf("Expected '110' as page number, got %s", value.Nr)
	}
}

func TestDownloadNonExistPage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	_, err := DownloadPage("888", server.URL)
	if err == nil {
		t.Errorf("Expected downloading a non-existing page to yield a 404 status")
	}

}

func TestConstructPageNr(t *testing.T) {
	tests := [][]string{
		{"@fewuihf Pagina 200", "200"},
		{"@fewuihf Pagina 404-1", "404-1"},
		{"@fewuihf Pagina 301", "301"},
		{"@fewuihf Pagina 501-22", "501-22"},
	}

	for _, c := range tests {
		nr, err := ConstructPageNr(c[0])
		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}
		if nr != c[1] {
			t.Errorf("Expected '%s' to return page %s but got %s", c[0], c[1], nr)
		}
	}

}
