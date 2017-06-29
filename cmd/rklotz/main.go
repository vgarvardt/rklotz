package main

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var version string

func main() {
	var versionFlag bool

	versionString := "rKlotz v" + version
	cobra.OnInitialize(func() {
		if versionFlag {
			fmt.Println(versionString)
			os.Exit(0)
		}

		log.SetFormatter(&log.JSONFormatter{})
	})

	var RootCmd = &cobra.Command{
		Use:   "rklotz",
		Short: "rKlotz is a simple one-user file-based blog engine.",
		Long:  versionString + ` rKlotz is a simple one-user file-based blog engine.`,
		Run:   RunServer,
	}
	RootCmd.Flags().BoolVarP(&versionFlag, "version", "v", false, "Print application version")

	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
