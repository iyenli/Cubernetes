/*
Copyright © 3022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"Cubernetes/pkg/apiserver/crudobj"
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"strings"
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
		switch strings.ToLower(args[0]) {

		case "pod", "pods":
			pods, err := crudobj.GetPods()
			if err != nil {
				log.Fatal("[FATAL] fail to get Pods")
				return
			}
			if len(pods) == 0 {
				fmt.Println("No Pods Found")
				return
			}
			fmt.Printf("%d Pods Found\n", len(pods))
			fmt.Printf("%-30s\t%-s\n", "Name", "UID")
			for _, pod := range pods {
				fmt.Printf("%-30s\t%-s\n", pod.Name, pod.UID)
			}

		case "service", "services", "svc", "svcs":
			svcs, err := crudobj.GetServices()
			if err != nil {
				log.Fatal("[FATAL] fail to get Services")
				return
			}
			if len(svcs) == 0 {
				fmt.Println("No Services Found")
				return
			}
			fmt.Printf("%d Services Found\n", len(svcs))
			fmt.Printf("%-30s\t%-s\n", "Name", "UID")
			for _, svc := range svcs {
				fmt.Printf("%-30s\t%-s\n", svc.Name, svc.UID)
			}

		case "replicaset", "replicasets", "rs", "rss":
			rss, err := crudobj.GetReplicaSets()
			if err != nil {
				log.Fatal("[FATAL] fail to get ReplicaSets")
				return
			}
			if len(rss) == 0 {
				fmt.Println("No ReplicaSets Found")
				return
			}
			fmt.Printf("%d ReplicaSets Found\n", len(rss))
			fmt.Printf("%-30s\t%-40s\t(%-v/%-v)\n", "Name", "UID", "running", "expected")
			for _, rs := range rss {
				var running int32
				if rs.Status != nil {
					running = rs.Status.RunningReplicas
				} else {
					running = 0
				}
				fmt.Printf("%-30s\t%-40s\t(%-v/%-v)\n", rs.Name, rs.UID, running, rs.Spec.Replicas)
			}

		case "node", "nodes":
			nodes, err := crudobj.GetNodes()
			if err != nil {
				log.Fatal("[FATAL] fail to get nodes")
				return
			}
			if len(nodes) == 0 {
				fmt.Println("No Nodes Found")
				return
			}
			fmt.Printf("%d Nodes found\n", len(nodes))
			fmt.Printf("%-30s\t%-40s\t%-v\n", "Name", "UID", "Ready")
			for _, node := range nodes {
				var ready bool
				if node.Status != nil {
					ready = node.Status.Condition.Ready
				} else {
					ready = false
				}
				fmt.Printf("%-30s\t%-40s\t%-v\n", node.Name, node.UID, ready)
			}

		case "dns", "dnses":
			dnses, err := crudobj.GetDnses()
			if err != nil {
				log.Fatal("[FATAL] fail to get Dnses")
				return
			}
			if len(dnses) == 0 {
				fmt.Println("No Dnses Found")
				return
			}
			fmt.Printf("%d Dnses Found\n", len(dnses))
			fmt.Printf("%-30s\t%-40s\t%-30s\t%-s\n", "Name", "UID", "Host", "PathCnt")
			for _, dns := range dnses {
				fmt.Printf("%-30s\t%-40s\t%-30s\t%-v\n", dns.Name, dns.UID, dns.Spec.Host, len(dns.Spec.Paths))
			}

		case "autoscaler", "autoscalers":
			autoScalers, err := crudobj.GetAutoScalers()
			if err != nil {
				log.Fatal("[FATAL] fail to get AutoScalers")
				return
			}
			if len(autoScalers) == 0 {
				fmt.Println("No AutoScalers Found")
				return
			}
			fmt.Printf("%d AutoScalers Found\n", len(autoScalers))
			fmt.Printf("%-30s\t%-40s\t(%-v/%-v)\n", "Name", "UID", "running", "expected")
			for _, as := range autoScalers {
				var running int
				var expected int
				if as.Status != nil {
					running = as.Status.ActualReplicas
					expected = as.Status.DesiredReplicas
				} else {
					running = 0
					expected = 0
				}
				fmt.Printf("%-30s\t%-40s\t(%-v/%-v)\n", as.Name, as.UID, running, expected)
			}

		case "job", "jobs", "gpujob", "gpujobs":
			jobs, err := crudobj.GetGpuJobs()
			if err != nil {
				log.Fatal("[FATAL] fail to get GpuJobs")
				return
			}
			if len(jobs) == 0 {
				fmt.Println("No GpuJobs Found")
				return
			}
			fmt.Printf("%d GpuJobs found\n", len(jobs))
			fmt.Printf("%-30s\t%-40s\t%-v\n", "Name", "UID", "Phase")
			for _, job := range jobs {
				fmt.Printf("%-30s\t%-40s\t%-v\n", job.Name, job.UID, job.Status.Phase)
			}

		case "action", "actions":
			actions, err := crudobj.GetActions()
			if err != nil {
				log.Fatal("[FATAL] fail to get Actions")
				return
			}
			if len(actions) == 0 {
				fmt.Println("No Actions Found")
				return
			}
			fmt.Printf("%d Actions found\n", len(actions))
			fmt.Printf("%-30s\t%-40s\t(%-v/%-v)\n", "Name", "UID", "running", "expected")
			for _, action := range actions {
				var running int
				var expected int
				if action.Status != nil {
					running = action.Status.ActualReplicas
					expected = action.Status.DesiredReplicas
				} else {
					running = 0
					expected = 0
				}
				fmt.Printf("%-30s\t%-40s\t(%-v/%-v)\n", action.Name, action.UID, running, expected)
			}

		case "ingress", "ingresses", "igs":
			ingresses, err := crudobj.GetIngresses()
			if err != nil {
				log.Fatal("[FATAL] fail to get Ingresses")
				return
			}
			if len(ingresses) == 0 {
				fmt.Println("No Ingresses Found")
				return
			}
			fmt.Printf("%d Ingresses found\n", len(ingresses))
			fmt.Printf("%-30s\t%-40s\t%-30s\t%-s\n", "Name", "UID", "Trigger", "InvokeAction")
			for _, ingress := range ingresses {
				fmt.Printf("%-30s\t%-40s\t%-30s\t%-v\n", ingress.Name, ingress.UID, ingress.Spec.TriggerPath, ingress.Spec.InvokeAction)
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
