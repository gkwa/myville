package cmd

import (
	"github.com/spf13/cobra"
	"github.com/taylormonacelli/myville/incus"
)

var imagesCmd = &cobra.Command{
	Use:   "images [filter] [container]",
	Short: "list local images",
	Long:  `equivilent to running 'incus image ls'`,
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
	rootCmd.AddCommand(imagesCmd)
	imagesCmd.Flags().BoolP("verbose", "v", false, "Enable verbose output")
}
