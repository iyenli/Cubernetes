package cmd

import (
	"Cubernetes/cmd/cuberoot/options"
	"github.com/spf13/cobra"
	"log"
	"os"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cuberoot",
	Short: "Init your cubernetes node in an elegant way",
	Long: `
	cuberoot is the starter of cubernetes, you can init your
	master and slave by init or join command.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	// Check log files:
	if _, err := os.Stat(options.LOGDIR); err != nil {
		err := os.Mkdir(options.LOGDIR, 0666)
		if err != nil {
			log.Panicln("Mkdir log failed. please run in su mode.")
			return
		}
	}

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
