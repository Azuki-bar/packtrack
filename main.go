package main

import (
	"log"
	"os"

	"github.com/azuki-bar/packtrack/cli"
)

func main() {
	log.Panic(cli.Main(os.Stdin, os.Stdout, os.Stderr))
}
