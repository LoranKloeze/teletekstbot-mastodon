package teletekst

import (
	"context"
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
	"github.com/mattn/go-mastodon"
)

var redisCli *redis.Client

func InitStore() *redis.Client {
	redisCli = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	return redisCli
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

func PageExists(store *redis.Client, p Page) bool {
	var title string
	key := fmt.Sprintf("teletekst:pages:%s:title", p.Nr)
	title, err := redisCli.Get(context.Background(), key).Result()
	if err != nil {
		title = ""
	}

	return title == p.Title
}

func InsertPage(store *redis.Client, p Page) {
	key := fmt.Sprintf("teletekst:pages:%s:title", p.Nr)
	err := redisCli.Set(context.Background(), key, p.Title, 0).Err()
	if err != nil {
		log.Fatal(err)
	}
}
