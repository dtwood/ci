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
	defer client.Close()

	jobs, err := client.LRange("jobs", 0, -1).Result()
	if err != nil {
		fmt.Println("Fatal error:", err)
		return nil
	}

	fmt.Println("jobs:", "["+strings.Join(jobs, ", ")+"]")
	return nil
}

func add(C *cli.Context) error {
	if !C.Args().Present() {
		fmt.Println("No command given")
		return nil
	}

	client := getConnection()
	defer client.Close()

	cmd := C.Args().First()

	count, err := client.LPush("jobs", cmd).Result()
	if err != nil {
		fmt.Println("Fatal error:", err)
		return nil
	}

	fmt.Println("Added", cmd)
	fmt.Println(count, "jobs in queue")

	return nil
}

func run(C *cli.Context) error {
	client := getConnection()
	defer client.Close()

	job, err := client.BRPopLPush("jobs", "processing", 0).Result()
	if err != nil {
		fmt.Println("Fatal error:", err)
		return nil
	}

	fmt.Println("running:", job)

	removed, err := client.LRem("processing", 1, job).Result()
	if err != nil {
		fmt.Println("Fatal error:", err)
		return nil
	}

	if removed != 1 {
		panic("no items removed after processing")
	}

	fmt.Println("run", job)

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
