package cmd

import (
	"github.com/spf13/cobra"
	"github.com/gkwa/myville/incus"
)

var imagesCmd = &cobra.Command{
	Use:     "images [filter] [container]",
	Short:   "list local images",
	Long:    `equivalent to running 'incus image ls'`,
	Aliases: []string{"images"},
	Run: func(cmd *cobra.Command, args []string) {
		var filter, container string
		if len(args) > 0 {
			filter = args[0]
		}
		if len(args) > 1 {
			container = args[1]
		}
		quiet, _ := cmd.Flags().GetBool("quiet")
		incus.ProcessImageCommand(filter, container, !quiet)
	},
}

func init() {
	rootCmd.AddCommand(imagesCmd)
	imagesCmd.Flags().BoolP("quiet", "q", false, "Disable verbose output")
}
