package cmd

import (
	"github.com/spf13/cobra"
	"github.com/taylormonacelli/myville/incus"
)

var imageCmd = &cobra.Command{
	Use:   "image [filter] [container]",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
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
	rootCmd.AddCommand(imageCmd)
	imageCmd.Flags().BoolP("verbose", "v", false, "Enable verbose output")
}
