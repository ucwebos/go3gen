package main

import (
	"github.com/spf13/cobra"
	"github.com/ucwebos/go3gen/handler"
	"log"
)

var rootCmd = &cobra.Command{
	Use:     "go3gen",
	Short:   "go3gen: An toolkit for golang code generate.",
	Long:    "go3gen: An toolkit for golang code generate..",
	Version: "0.0.1",
}

func init() {
	rootCmd.AddCommand(handler.CmdList()...)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
