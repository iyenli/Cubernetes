/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"Cubernetes/pkg/apiserver/crudobj"
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"strings"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete an object from Cubernetes",
	Long: `
Delete an object from Cubernetes
for example:
	cubectl delete pod nginx:452cbd60-131c-4efa-9e06-7b364692a737
	cubectl delete [Object kind] [UID]
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			log.Fatal("[FATAL] lack arguments")
		}
		switch strings.ToLower(args[0]) {
		case "pod":
			err := crudobj.DeletePod(args[1])
			if err != nil {
				log.Fatal("[FATAL] fail to delete Pod")
			} else {
				fmt.Printf("Pod UID=%s deleted\n", args[1])
			}
		case "service", "svc":
			err := crudobj.DeleteService(args[1])
			if err != nil {
				log.Fatal("[FATAL] fail to delete Service")
			} else {
				fmt.Printf("Service UID=%s deleted\n", args[1])
			}
		case "replicaset", "rs":
			err := crudobj.DeleteReplicaSet(args[1])
			if err != nil {
				log.Fatal("[FATAL] fail to delete ReplicaSet")
			} else {
				fmt.Printf("ReplicaSet UID=%s deleted\n", args[1])
			}
		case "dns":
			err := crudobj.DeleteDns(args[1])
			if err != nil {
				log.Fatal("[FATAL] fail to delete Dns")
			} else {
				fmt.Printf("Dns UID=%s deleted\n", args[1])
			}
		case "autoscaler":
			err := crudobj.DeleteAutoScaler(args[1])
			if err != nil {
				log.Fatal("[FATAL] fail to delete AutoScaler")
			} else {
				fmt.Printf("AutoScaler UID=%s deleted\n", args[1])
			}
		case "job", "gpujob":
			err := crudobj.DeleteGpuJob(args[1])
			if err != nil {
				log.Fatal("[FATAL] fail to delete GpuJob")
			} else {
				fmt.Printf("GpuJob UID=%s deleted\n", args[1])
			}
		default:
			log.Fatal("[FATAL] Unknown kind: " + args[0])
		}
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deleteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deleteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
