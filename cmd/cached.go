package main

import (
	"github.com/abatilo/go-gcache-redis-sample/cmd/cachedsvc"
	"github.com/spf13/cobra"
)

func main() {
	var (
		rootCmd = &cobra.Command{
			Use:   "cached",
			Short: "A small sample application with requests backed by Redis",
		}
	)

	rootCmd.AddCommand(cachedsvc.Cmd)
	rootCmd.Execute()
}
