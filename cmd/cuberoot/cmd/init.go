package cmd

import (
	"Cubernetes/cmd/cuberoot/options"
	"Cubernetes/cmd/cuberoot/utils"
	"github.com/spf13/cobra"
	"log"
	"net"
)

// getCmd represents the get command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Init an API Server as a cubernetes master",
	Long: `
Init an API Server as a cubernetes master
usage:
	cuberoot join local-ip
example:
	cuberoot join 192.168.1.5`,

	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			log.Fatal("[FATAL] lack arguments")
			return
		}
		if ip := net.ParseIP(args[0]); ip == nil {
			log.Fatalf("[FATAL] illegal ip address: %v", args[0])
			return
		}

		err := utils.StartDaemonProcess(options.APISERVERLOG, options.APISERVER, args[0])
		if err != nil {
			return
		}
		err = utils.StartDaemonProcess(options.CUBEPROXYLOG, options.CUBEPROXY, args[0])
		if err != nil {
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
