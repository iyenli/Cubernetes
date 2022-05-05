/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"Cubernetes/pkg/apiserver/crudobj"
	"fmt"
	"gopkg.in/yaml.v2"
	"log"

	"github.com/spf13/cobra"
)

// describeCmd represents the describe command
var describeCmd = &cobra.Command{
	Use:   "describe",
	Short: "Describe detailed information of an object",
	Long: `
Describe detailed information of an object
for example:
	cubectl describe pod nginx:452cbd60-131c-4efa-9e06-7b364692a737
	cubectl describe [Object kind] [UID]`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			log.Fatal("[FATAL] lack arguments")
			return
		}
		UID := args[1]
		switch args[0] {
		case "pod", "Pod":
			pod, err := crudobj.GetPod(UID)
			if err != nil {
				log.Fatal("[FATAL] fail to get pod")
				return
			}
			str, err := yaml.Marshal(pod)
			if err != nil {
				log.Fatal("[FATAL] fail to marshall pod")
				return
			}
			fmt.Print(string(str))
		case "service", "Service", "svc":
			svc, err := crudobj.GetService(UID)
			if err != nil {
				log.Fatal("[FATAL] fail to get service")
				return
			}
			str, err := yaml.Marshal(svc)
			if err != nil {
				log.Fatal("[FATAL] fail to marshall service")
				return
			}
			fmt.Print(string(str))
		default:
			log.Fatal("[FATAL] Unknown kind: " + args[0])
		}
	},
}

func init() {
	rootCmd.AddCommand(describeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// describeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// describeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
