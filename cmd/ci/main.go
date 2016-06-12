package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/dtwood/ci/client"
	"github.com/urfave/cli"
	"gopkg.in/redis.v3"
)

var version string

func getConnection(C *cli.Context) *redis.Client {
	return client.GetConnection(C.GlobalString("server"), C.GlobalString("password"), int64(C.GlobalInt("database-id")))
}

func jobsList(C *cli.Context) error {
	cl := getConnection(C)
	defer cl.Close()

	res, err := client.ListTasks(cl)
	if err != nil {
		return cli.NewExitError(fmt.Sprint("Error:", err), 1)
	}

	fmt.Println(strings.Join(res, "\n"))

	return nil
}

func jobsAdd(C *cli.Context) error {
	cl := getConnection(C)
	defer cl.Close()

	if !C.Args().Present() {
		return cli.NewExitError("No command specified", 2)
	}

	count, err := client.Add(cl, C.Args().First())
	if err != nil {
		return cli.NewExitError(fmt.Sprint("Error:", err), 1)
	}

	fmt.Println("Queued:", count)

	return nil
}

func jobsRun(C *cli.Context) error {
	cl := getConnection(C)
	defer cl.Close()

	res, err := client.Run(cl)
	if err != nil {
		return cli.NewExitError(fmt.Sprint("Error:", err), 1)
	}

	fmt.Println("Run:", res)

	return nil
}

func botsList(C *cli.Context) error {
	cl := getConnection(C)
	defer cl.Close()

	res, err := client.List(cl)
	if err != nil {
		return cli.NewExitError(fmt.Sprint("Error:", err), 1)
	}

	fmt.Println(strings.Join(res, "\n"))

	return nil
}

func main() {
	app := cli.NewApp()

	app.Commands = []cli.Command{
		{
			Name:    "jobs",
			Aliases: []string{"j"},
			Subcommands: []cli.Command{
				{
					Name:    "list",
					Aliases: []string{"l"},
					Usage:   "list all pending tasks",
					Action:  jobsList,
				},
				{
					Name:    "add",
					Aliases: []string{"a"},
					Usage:   "add a new task",
					Action:  jobsAdd,
				},
				{
					Name:    "run",
					Aliases: []string{"r"},
					Usage:   "run the oldest task",
					Action:  jobsRun,
				},
			},
		},
		{
			Name: "bots",
			Subcommands: []cli.Command{
				{
					Name:    "list",
					Aliases: []string{"l"},
					Usage:   "list all connected build bots",
					Action:  botsList,
				},
			},
		},
	}
	app.EnableBashCompletion = true
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "server",
			Value:  "localhost:6379",
			EnvVar: "CI_SERVER",
		},
		cli.StringFlag{
			Name:   "password",
			Value:  "",
			EnvVar: "CI_PASSWORD",
		},
		cli.IntFlag{
			Name:   "database-id",
			Value:  0,
			EnvVar: "CI_DATABASE_ID",
		},
	}

	if version == "" {
		app.Version = "devel"
	} else {
		app.Version = version
	}

	app.Run(os.Args)
}
