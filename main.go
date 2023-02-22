package main

import (
	"log"
	"os"

	"github.com/azuki-bar/packtrack/cli"
)

func main() {
	err := (cli.Main(os.Stdin, os.Stdout, os.Stderr))
	if err != nil {
		log.Print(err)
	}
}
