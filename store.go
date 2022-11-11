package main

import (
	"context"
	"log"

	"github.com/go-redis/redis/v8"
	"github.com/mattn/go-mastodon"
)

var redisCli *redis.Client

func InitStore() {
	redisCli = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
}

func LastNotificationId() mastodon.ID {
	id, err := redisCli.Get(context.Background(), "teletekst:last_notification_id").Result()
	if err != nil {
		id = ""
	}
	return mastodon.ID(id)
}

func SetLastNotificationId(id mastodon.ID) {
	err := redisCli.Set(context.Background(), "teletekst:last_notification_id", string(id), 0).Err()
	if err != nil {
		log.Fatal(err)
	}
}
