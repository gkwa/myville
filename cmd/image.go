package cmd

import (
	"github.com/spf13/cobra"
	"github.com/taylormonacelli/myville/incus"
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

var rmImageCmd = &cobra.Command{
	Use:   "rm [filters...]",
	Short: "Remove images matching filters",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			if err := cmd.Usage(); err != nil {
				cmd.PrintErr(err)
			}
			return
		}
		quiet, _ := cmd.Flags().GetBool("quiet")
		dryRun, _ := cmd.Flags().GetBool("dry-run")
		incus.ProcessImageRemoveCommand(args, !quiet, dryRun)
	},
}

func init() {
	rootCmd.AddCommand(imagesCmd)
	imagesCmd.Flags().BoolP("quiet", "q", false, "Disable verbose output")

	imagesCmd.AddCommand(rmImageCmd)
	rmImageCmd.Flags().BoolP("quiet", "q", false, "Disable verbose output")
	rmImageCmd.Flags().BoolP("dry-run", "n", false, "Show what would be deleted without actually deleting")
}
