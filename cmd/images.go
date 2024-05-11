package cmd

import (
	"github.com/spf13/cobra"
	"github.com/taylormonacelli/myville/incus"
)

var imagessCmd = &cobra.Command{
	Use:   "imagess [filter] [container]",
	Short: "list local imagess",
	Long:  `equivilent to running 'incus images ls'`,
	Run: func(cmd *cobra.Command, args []string) {
		var filter, container string
		if len(args) > 0 {
			filter = args[0]
		}
		if len(args) > 1 {
			container = args[1]
		}
		verbose, _ := cmd.Flags().GetBool("verbose")
		incus.ProcessImageCommand(filter, container, verbose)
	},
}

func init() {
	rootCmd.AddCommand(imagessCmd)
	imagessCmd.Flags().BoolP("verbose", "v", false, "Enable verbose output")
}
