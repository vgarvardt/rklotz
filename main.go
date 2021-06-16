package main

import (
	"context"
	"log"

	"github.com/vgarvardt/rklotz/cmd"
)

var version = "0.0.0-dev"

func main() {
	rootCmd := cmd.NewRootCmd(version)

	err := rootCmd.ExecuteContext(context.Background())
	if err != nil {
		log.Fatal(err)
	}
}
