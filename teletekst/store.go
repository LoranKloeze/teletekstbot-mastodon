package teletekst

import (
	"context"
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
	"github.com/mattn/go-mastodon"
)

var redisCli *redis.Client
var redisNamespace string

func InitStore(namespace string) *redis.Client {
	redisNamespace = namespace
	redisCli = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	return redisCli
}

func ClearStore() {
	cursor := uint64(0)
	for {
		keys, cursor, err := redisCli.Scan(context.Background(), cursor, redisNamespace+"*", 0).Result()
		if err != nil {
			log.Fatal(err)
		}

		for _, key := range keys {
			redisCli.Del(context.Background(), key)
		}
		if cursor == 0 {
			break
		}
	}

}

func LastNotificationId() mastodon.ID {
	key := fmt.Sprintf("%s:last_notification_id", redisNamespace)
	id, err := redisCli.Get(context.Background(), key).Result()
	if err != nil {
		id = ""
	}
	return mastodon.ID(id)
}

func SetLastNotificationId(id mastodon.ID) {
	key := fmt.Sprintf("%s:last_notification_id", redisNamespace)
	err := redisCli.Set(context.Background(), key, string(id), 0).Err()
	if err != nil {
		log.Fatal(err)
	}
}

func PageExists(p Page) bool {
	var title string
	key := fmt.Sprintf("%s:pages:%s:title", redisNamespace, p.Nr)
	title, err := redisCli.Get(context.Background(), key).Result()
	if err != nil {
		title = ""
	}

	return title == p.Title
}

func InsertPage(p Page) {
	key := fmt.Sprintf("%s:pages:%s:title", redisNamespace, p.Nr)
	err := redisCli.Set(context.Background(), key, p.Title, 0).Err()
	if err != nil {
		log.Fatal(err)
	}
}
