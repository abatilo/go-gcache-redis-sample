package cachedsvc

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	// Cmd is the exported cobra command which starts the webhook handler service
	Cmd = &cobra.Command{
		Use:   "svc",
		Short: "Runs the web service",
		Run: func(cmd *cobra.Command, args []string) {
			main()
		},
	}
)

func main() {
	fmt.Println("Hi")
}
