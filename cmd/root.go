package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"time"
)

var RetryCount int
var Sleep time.Duration
var rootCmd = &cobra.Command{
	Use:   "wait4x",
	Short: "wait4x allows waiting for a port or a service to enter into specify state",
	Long: `wait4x allows waiting for a port to enter into specify state or waiting for a service e.g. redis, mysql, postgres, ... to enter inter ready state`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("tmplBytes.String()")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
