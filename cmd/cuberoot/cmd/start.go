package cmd

import (
	"Cubernetes/cmd/cuberoot/utils"
	"github.com/spf13/cobra"
)

// startCmd represents start as a node from local registry
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start an registered node",
	Long: `
Start as a node from local registry
usage:
	cuberoot start`,

	Run: func(cmd *cobra.Command, args []string) {
		utils.StartFromRegistry()
	},
}

func init() {
	rootCmd.AddCommand(startCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
