package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/taylormonacelli/myville/incus"
)

var containerCmd = &cobra.Command{
	Use:   "container [filter] [name]",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
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
		verbose, _ := cmd.Flags().GetBool("verbose")
		incus.ProcessContainerCommand(filter, name, verbose)
	},
}

func init() {
	rootCmd.AddCommand(containerCmd)
	containerCmd.Flags().BoolP("verbose", "v", false, "Enable verbose output")
}
