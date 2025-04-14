package main

import (
	"flag"

	"template/internal/app"
)

var (
	envParse bool
	envPath  string
)

func init() {
	flag.BoolVar(&envParse, "env.parse", true, "Whether parse envs from file or not")
	flag.StringVar(&envPath, "env.path", "internal/app/config/local.env", "Path to env file")
}

func main() {
	flag.Parse()

	files := make([]string, 0)
	if envParse {
		files = append(files, envPath)
	}

	app.Run(files...)
}
