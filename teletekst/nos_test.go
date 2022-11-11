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

	value := DownloadPage("110", server.URL)
	if value.Content != "sample_content" {
		t.Errorf("Expected 'sample_content' as page content, got %s", value.Content)
	}
	if value.Nr != "110" {
		t.Errorf("Expected '110' as page number, got %s", value.Nr)
	}
}
