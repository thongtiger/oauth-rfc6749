package redis

import (
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis"
)

// RedisClient connector
func RedisClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	pong, err := client.Ping().Result()
	fmt.Println(pong, err)
	// Output: PONG <nil>
	return client
}

func Exists(ClientID, RefreshToken string) bool {
	client := RedisClient()
	defer client.Close()

	dbValue, err := client.Get(ClientID).Result()
	if err == redis.Nil {
		log.Println("redis: refresh_token does not exist or expired")
		return false
	}
	if dbValue == "" {
		log.Println("redis: refresh_token is expired")
		return false
	}
	// check value
	if dbValue == RefreshToken {
		return true
	}
	log.Println("redis: refresh_token does not exist")
	return false
}

func SetRefreshToken(ClientID, RefreshToken string, ExpireIn time.Duration) (bool, error) {
	client := RedisClient()
	defer client.Close()

	err := client.Set(ClientID, RefreshToken, ExpireIn).Err()
	if err != nil {
		log.Println("redis: error set new refresh_token")
		return false, fmt.Errorf("redis: error set new refresh_token")
	}
	return true, nil
}
