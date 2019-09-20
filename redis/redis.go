package redis

import (
	"fmt"
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

func Exists(ClientID, RefreshToken string) (bool, error) {
	client := RedisClient()
	defer client.Close()

	dbValue, err := client.Get(ClientID).Result()
	if err == redis.Nil {
		return false, fmt.Errorf("redis: refresh_token does not exist or expired")
	}
	if dbValue == "" {
		return false, fmt.Errorf("redis: refresh_token is expired")
	}
	// check value
	if dbValue == RefreshToken {
		return true, nil
	}
	return false, fmt.Errorf("redis: refresh_token does not exist")
}

func SetRefreshToken(ClientID, RefreshToken string, ExpireIn time.Duration) (bool, error) {
	client := RedisClient()
	defer client.Close()

	err := client.Set(ClientID, RefreshToken, ExpireIn).Err()
	if err != nil {
		return false, fmt.Errorf("redis: error set new refresh_token ")
	}
	return true, nil
}
