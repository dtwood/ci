package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/urfave/cli"
	"gopkg.in/redis.v3"
)

func getConnection() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return client
}

func list(c *cli.Context) error {
	client := getConnection()
	jobs := client.LRange("jobs", 0, -1).Val()
	fmt.Println("[" + strings.Join(jobs, ", ") + "]")
	return nil
}

func add(C *cli.Context) error {
	if !C.Args().Present() {
		fmt.Println("No command given")
		return nil
	}

	client := getConnection()
	client.RPush("jobs", C.Args().First())
	return nil
}

func run(C *cli.Context) error {
	client := getConnection()
	job, err := client.BLPop(0, "jobs").Result()
	if err != nil {
		panic(err)
	}

	fmt.Println("running:", job[1])

	return nil
}

func main() {
	app := cli.NewApp()

	app.EnableBashCompletion = true
	app.Commands = []cli.Command{
		{
			Name:   "list",
			Usage:  "list all pending tasks",
			Action: list,
		},
		{
			Name:   "add",
			Usage:  "add a new task",
			Action: add,
		},
		{
			Name:   "run",
			Usage:  "run the oldest task",
			Action: run,
		},
	}

	app.Run(os.Args)
}
