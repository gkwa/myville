package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/taylormonacelli/myville/incus"
)

var containerCmd = &cobra.Command{
	Use:     "container filter container_name",
	Short:   "Create a container from a local image matching a filter",
	Aliases: []string{"c"},
	Long: `Create a new container from an existing local image that matches the provided filter.

Usage:
  myville container merry test-container

This will search for local images containing "merry" in their name and create
a new container named "test-container" from the matching image.

If multiple images match the filter, the most recent one will be used.
Use the --quiet flag to suppress output during container creation.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			err := cmd.Usage()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			return
		}
		filter := args[0]
		name := args[1]
		quiet, _ := cmd.Flags().GetBool("quiet")
		incus.ProcessContainerCommand(filter, name, !quiet)
	},
}


func init() {
	rootCmd.AddCommand(containerCmd)
	containerCmd.Flags().BoolP("quiet", "q", false, "Disable verbose output")
}
