// Copyright 2022 Codedivision - Loran Kloeze. All rights reserved.
// Use of this source code is governed by the MIT-license.

package main

import (
	"testing"
)

func TestMD5Hash(t *testing.T) {
	value := MD5Hash("teletekst")
	if value != "544d7f701865fc2971eda7bd211d1d1c" {
		t.Errorf("Expected '544d7f701865fc2971eda7bd211d1d1c', got %s", value)
	}
}
