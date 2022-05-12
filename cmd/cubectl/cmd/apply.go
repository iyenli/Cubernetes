/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"Cubernetes/pkg/apiserver/crudobj"
	"Cubernetes/pkg/object"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
)

// applyCmd represents the apply command
var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply a yaml configuration file to Cubernetes",
	Long: `
Apply a yaml configuration file to Cubernetes
for example:
	cubectl apply -f pod.yaml
	cubectl apply -f [file path]`,
	Run: func(cmd *cobra.Command, args []string) {
		f, err := cmd.Flags().GetString("file")
		if err != nil {
			log.Fatal("[FATAL] missing input config file")
		}
		file, err := ioutil.ReadFile(f)
		if err != nil {
			log.Fatal("[FATAL] cannot read input config file")
		}

		var t object.TypeMeta
		err = yaml.Unmarshal(file, &t)
		if err != nil {
			log.Fatal("[FATAL] fail to unmarshal config file")
			return
		}
		switch t.Kind {
		case "Pod":
			var pod object.Pod
			err = yaml.Unmarshal(file, &pod)
			if err != nil {
				log.Fatal("[FATAL] fail to parse Pod")
				return
			}
			newPod, err := crudobj.CreatePod(pod)
			if err != nil {
				log.Fatal("[FATAL] fail to create new Pod")
				return
			}
			log.Printf("Pod UID=%s created\n", newPod.UID)

		case "Service":
			var service object.Service
			err = yaml.Unmarshal(file, &service)
			if err != nil {
				log.Fatal("[FATAL] fail to parse Service")
				return
			}
			newService, err := crudobj.CreateService(service)
			if err != nil {
				log.Fatal("[FATAL] fail to create new Service")
				return
			}
			log.Printf("Service UID=%s created\n", newService.UID)
		case "ReplicaSet":
			var rs object.ReplicaSet
			err = yaml.Unmarshal(file, &rs)
			if err != nil {
				log.Fatal("[FATAL] fail to parse ReplicaSet")
				return
			}
			newRs, err := crudobj.CreateReplicaSet(rs)
			if err != nil {
				log.Fatal("[FATAL] fail to create new ReplicaSet")
				return
			}
			log.Printf("ReplicaSet UID=%s created\n", newRs.UID)
		default:
			log.Fatal("[FATAL] Unknown kind: " + t.Kind)
		}
	},
}

func init() {
	rootCmd.AddCommand(applyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// applyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// applyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	applyCmd.Flags().StringP("file", "f", "", "path of your config yaml file")
}
