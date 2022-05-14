/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"Cubernetes/pkg/apiserver/crudobj"
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"strings"

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
		switch strings.ToLower(args[0]) {
		case "pod":
			pod, err := crudobj.GetPod(UID)
			if err != nil {
				log.Fatal("[FATAL] fail to get Pod")
				return
			}
			str, err := yaml.Marshal(pod)
			if err != nil {
				log.Fatal("[FATAL] fail to marshall Pod")
				return
			}
			fmt.Print(string(str))
		case "service", "svc":
			svc, err := crudobj.GetService(UID)
			if err != nil {
				log.Fatal("[FATAL] fail to get Service")
				return
			}
			str, err := yaml.Marshal(svc)
			if err != nil {
				log.Fatal("[FATAL] fail to marshall Service")
				return
			}
			fmt.Print(string(str))
		case "replicaset", "rs":
			rs, err := crudobj.GetReplicaSet(UID)
			if err != nil {
				log.Fatal("[FATAL] fail to get ReplicaSet")
				return
			}
			str, err := yaml.Marshal(rs)
			if err != nil {
				log.Fatal("[FATAL] fail to marshall ReplicaSet")
				return
			}
			fmt.Print(string(str))
		case "node":
			node, err := crudobj.GetNode(UID)
			if err != nil {
				log.Fatal("[FATAL] fail to get Node")
				return
			}
			str, err := yaml.Marshal(node)
			if err != nil {
				log.Fatal("[FATAL] fail to marshall Node")
				return
			}
			fmt.Print(string(str))
		case "dns":
			dns, err := crudobj.GetDns(UID)
			if err != nil {
				log.Fatal("[FATAL] fail to get Dns")
				return
			}
			str, err := yaml.Marshal(dns)
			if err != nil {
				log.Fatal("[FATAL] fail to marshall Dns")
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
