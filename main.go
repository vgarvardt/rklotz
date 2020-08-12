package main

import (
	"log"

	"github.com/vgarvardt/rklotz/cmd"
)

var version = "0.0.0-dev"

func main() {
	rootCmd := cmd.NewRootCmd(version)

	err := rootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
