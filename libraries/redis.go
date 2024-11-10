package libraries

import (
	"context"
	"log"
	"sync"

	"github.com/go-redis/redis/v8"
)

var lock = &sync.Mutex{}

type Singleton struct {
	client *redis.Client
}

var instance *Singleton

func GetInstance() *redis.Client {
	if instance == nil {
		lock.Lock()
		defer lock.Unlock()
		if instance == nil {
			instance = &Singleton{
				client: GetClient(),
			}
		}
	}
	return instance.client
}

func GetClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})
	ctx := context.Background()         // Add this line
	_, err := client.Ping(ctx).Result() // Update this line
	if err != nil {
		log.Fatal("SERVER - Error connecting to redis")
	}
	log.Print("SERVER - Connected to redis")
	return client
}
