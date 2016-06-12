package client

import (
	"strings"

	"gopkg.in/redis.v3"
)

func List(client *redis.Client) ([]string, error) {
	clients, err := client.ClientList().Result()
	if err != nil {
		return nil, err
	}

	return strings.Split(clients, "\n"), nil
}
