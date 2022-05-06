package cmd

import (
	"Cubernetes/cmd/cuberoot/options"
	"github.com/spf13/cobra"
	"log"
	"net"
	"os"
	"os/exec"
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

		_, err := os.Stat(options.APISERVER)
		if err != nil {
			log.Panicf("API Server not found: %v", err.Error())
			return
		}

		server := exec.Command(options.APISERVER, args[0])
		stdout, err := os.OpenFile(options.APISERVERLOG, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
		if err != nil {
			log.Println(os.Getpid(), ": open log file error", err)
		}
		server.Stderr = stdout
		server.Stdout = stdout
		err = server.Start()

		if err != nil {
			return
		}
		// TODO: Start cubeproxy
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
