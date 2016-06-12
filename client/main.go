package client

import (
	"gopkg.in/redis.v3"
)

func GetConnection(server string, password string, database_id int64) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     server,
		Password: password,
		DB:       database_id,
	})

	return client
}
