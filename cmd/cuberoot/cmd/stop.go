package cmd

import (
	"Cubernetes/cmd/cuberoot/options"
	"Cubernetes/cmd/cuberoot/utils"
	"github.com/spf13/cobra"
	"log"
)

// stopCmd Stop all process about Cubernetes
var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop all process about Cubernetes",
	Long: `
Stop all process about Cubernetes
usage:
	cuberoot stop`,

	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 2 {
			log.Fatal("[FATAL] too much arguments")
			return
		}

		err := utils.KillDaemonProcess(options.APISERVER)
		if err != nil {
			//log.Printf("Error when killing process: %v", err.Error())
		}

		err = utils.KillDaemonProcess(options.CUBELET)
		if err != nil {
			//log.Printf("Error when killing process: %v", err.Error())
		}

		err = utils.KillDaemonProcess(options.CUBEPROXY)
		if err != nil {
			//log.Printf("Error when killing process: %v", err.Error())
		}

		err = utils.KillDaemonProcess(options.MANAGER)
		if err != nil {
			//log.Printf("Error when killing process: %v", err.Error())
		}

		err = utils.KillDaemonProcess(options.ETCD)
		if err != nil {
			//log.Printf("Error when killing process: %v", err.Error())
		}

		err = utils.KillDaemonProcess(options.SCHEDULER)
		if err != nil {
			//log.Printf("Error when killing process: %v", err.Error())
		}

		//err = utils.KillDaemonProcess(options.GATEWAY)
		//if err != nil {
		//	//log.Printf("Error when killing process: %v", err.Error())
		//}

		err = utils.KillDaemonProcess(options.BRAIN)
		if err != nil {
			//log.Printf("Error when killing process: %v", err.Error())
		}

	},
}

func init() {
	rootCmd.AddCommand(stopCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
