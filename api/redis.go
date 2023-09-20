package api

import (
	"context"
	"fmt"
	"strconv"

	"github.com/adarsh-a-tw/tt-backend/service"
	"github.com/redis/go-redis/v9"
)

const channel string = "match_update_notifications"

func PublishMatchChange(matchId int, rdb *redis.Client) {
	err := rdb.Publish(context.Background(), channel, fmt.Sprintf("%d", matchId)).Err()
	if err != nil {
		panic(err)
	}
}

func SubscribeToMatchChanges(rdb *redis.Client, svc service.Service) {
	pubsub := rdb.Subscribe(context.Background(), channel)

	defer pubsub.Close()
	ch := pubsub.Channel()

	for msg := range ch {
		matchId, err := strconv.Atoi(msg.Payload)
		if err != nil {
			panic(err)
		}
		NotifySubscribers(matchId, svc)
	}
}
