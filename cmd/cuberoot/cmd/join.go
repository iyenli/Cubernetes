package cmd

import (
	"Cubernetes/cmd/cuberoot/options"
	"Cubernetes/cmd/cuberoot/utils"
	"github.com/spf13/cobra"
	"log"
	"net"
)

// joinCmd represents join master as a slave
var joinCmd = &cobra.Command{
	Use:   "join",
	Short: "Join exist API Server as a slave",
	Long: `
Join exist API Server as a slave
usage:
	cuberoot join local-ip api-server-ip
example:
	cuberoot join 192.168.1.5 192.168.1.11`,

	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			log.Fatal("[FATAL] lack arguments")
			return
		}
		if ip := net.ParseIP(args[0]); ip == nil {
			log.Fatalf("[FATAL] illegal ip address: %v", args[0])
			return
		}
		if ip := net.ParseIP(args[1]); ip == nil {
			log.Fatalf("[FATAL] illegal ip address: %v", args[1])
			return
		}

		err := utils.StartDaemonProcess(options.CUBELETLOG, options.CUBELET, args[0], args[1])
		if err != nil {
			return
		}
		err = utils.StartDaemonProcess(options.CUBEPROXYLOG, options.CUBEPROXY, args[0], args[1])
		if err != nil {
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(joinCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
