package client

import (
	"gopkg.in/redis.v3"
)

func ListTasks(client *redis.Client) ([]string, error) {
	jobs, err := client.LRange("jobs", 0, -1).Result()
	if err != nil {
		return nil, err
	}

	return jobs, nil
}

func Add(client *redis.Client, cmd string) (int64, error) {
	count, err := client.LPush("jobs", cmd).Result()
	if err != nil {
		return 0, err
	}

	return count, nil
}

func Run(client *redis.Client) (string, error) {
	job, err := client.BRPopLPush("jobs", "processing", 0).Result()
	if err != nil {
		return "", err
	}

	removed, err := client.LRem("processing", 1, job).Result()
	if err != nil {
		return "", err
	}

	if removed != 1 {
		panic("no items removed after processing")
	}

	return job, nil
}
