// Copyright 2022 Codedivision - Loran Kloeze. All rights reserved.
// Use of this source code is governed by the MIT-license.

package teletekst

import (
	"context"
	"testing"

	"github.com/mattn/go-mastodon"
)

func TestClearStore(t *testing.T) {
	InitStore("teletekst_test")
	key := "teletekst_test:imnothere"

	redisCli.Set(context.Background(), key, "xxx", 0)
	v1, _ := redisCli.Get(context.Background(), key).Result()
	if v1 != "xxx" {
		t.Errorf("Expected key '%s' to have 'xxx'", key)
	}

	ClearStore()

	v2, _ := redisCli.Get(context.Background(), key).Result()
	if v2 != "" {
		t.Errorf("Expected key '%s' to have nothing", key)
	}

}
func TestSetLastNotificationId(t *testing.T) {
	InitStore("teletekst_test")
	ClearStore()
	id := mastodon.ID("1337")
	SetLastNotificationId(id)
	key := "teletekst_test:last_notification_id"
	value, err := redisCli.Get(context.Background(), key).Result()
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	if value != "1337" {
		t.Errorf("Expected '1337', got %s", value)
	}

}

func TestLastNotificationId(t *testing.T) {
	InitStore("teletekst_test")
	ClearStore()
	key := "teletekst_test:last_notification_id"
	redisCli.Set(context.Background(), key, "1337", 0)
	id := LastNotificationId()
	if id != "1337" {
		t.Errorf("Expected '1337', got %s", id)
	}
}

func TestPageExists(t *testing.T) {
	InitStore("teletekst_test")
	ClearStore()

	key := "teletekst_test:pages:200:title"
	redisCli.Set(context.Background(), key, "aTitle", 0).Err()

	if !PageExistsInDb(Page{Nr: "200", Title: "aTitle"}) {
		t.Errorf("Expected page 200 to exist")
	}

}

func TestInsertPage(t *testing.T) {
	InitStore("teletekst_test")
	ClearStore()

	InsertPage(Page{Nr: "200", Title: "aTitle"})

	key := "teletekst_test:pages:200:title"
	value, _ := redisCli.Get(context.Background(), key).Result()
	if value != "aTitle" {
		t.Errorf("Expected 'aTitle', got %s", value)
	}

}
