/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"Cubernetes/pkg/apiserver/crudobj"
	"fmt"
	"github.com/spf13/cobra"
	"log"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get all objects of a certain kind",
	Long: `
Get all object of a certain kind
for example:
	cubectl get pods
	cubectl get svcs`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			log.Fatal("[FATAL] lack arguments")
			return
		}
		switch args[0] {
		case "pods", "pod", "Pods", "Pod":
			pods, err := crudobj.GetPods()
			if err != nil {
				log.Fatal("[FATAL] fail to get pods")
				return
			}
			if len(pods) == 0 {
				fmt.Println("No pods found")
				return
			}
			fmt.Printf("%d pods found\n", len(pods))
			fmt.Printf("%-30s\t%-s\n", "Name", "UID")
			for _, pod := range pods {
				fmt.Printf("%-30s\t%-s\n", pod.Name, pod.UID)
			}

		case "service", "services", "Service", "Services", "svc", "svcs":
			svcs, err := crudobj.GetServices()
			if err != nil {
				log.Fatal("[FATAL] fail to get pods")
				return
			}
			if len(svcs) == 0 {
				fmt.Println("No services found")
				return
			}
			fmt.Printf("%d services found\n", len(svcs))
			fmt.Printf("%-30s\t%-s\n", "Name", "UID")
			for _, svc := range svcs {
				fmt.Printf("%-30s\t%-s\n", svc.Name, svc.UID)
			}
		default:
			log.Fatal("[FATAL] Unknown kind: " + args[0])
		}
	},
}

func init() {
	rootCmd.AddCommand(getCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
