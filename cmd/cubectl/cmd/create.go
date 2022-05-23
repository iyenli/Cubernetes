/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"Cubernetes/pkg/apiserver/crudobj"
	"Cubernetes/pkg/apiserver/objfile"
	"Cubernetes/pkg/object"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"

	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create objects from a config file",
	Long: `
Create objects from a config file
for example:
	cubectl create -f pod.yaml
	cubectl create -f [file path]`,
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
		}

		switch t.Kind {

		case object.KindPod:
			var pod object.Pod
			err = yaml.Unmarshal(file, &pod)
			if err != nil {
				log.Fatal("[FATAL] fail to parse Pod")
			}
			newPod, err := crudobj.CreatePod(pod)
			if err != nil {
				log.Fatal("[FATAL] fail to create new Pod")
			}
			log.Printf("Pod UID=%s created\n", newPod.UID)

		case object.KindService:
			var service object.Service
			err = yaml.Unmarshal(file, &service)
			if err != nil {
				log.Fatal("[FATAL] fail to parse Service")
			}
			newService, err := crudobj.CreateService(service)
			if err != nil {
				log.Fatal("[FATAL] fail to create new Service")
			}
			log.Printf("Service UID=%s created\n", newService.UID)

		case object.KindReplicaSet:
			var rs object.ReplicaSet
			err = yaml.Unmarshal(file, &rs)
			if err != nil {
				log.Fatal("[FATAL] fail to parse ReplicaSet")
			}
			newRs, err := crudobj.CreateReplicaSet(rs)
			if err != nil {
				log.Fatal("[FATAL] fail to create new ReplicaSet")
			}
			log.Printf("ReplicaSet UID=%s created\n", newRs.UID)

		case object.KindDns:
			var dns object.Dns
			err = yaml.Unmarshal(file, &dns)
			if err != nil {
				log.Fatal("[FATAL] fail to parse Dns", err)
			}
			newDns, err := crudobj.CreateDns(dns)
			if err != nil {
				log.Fatal("[FATAL] fail to create new Dns")
			}
			log.Printf("Dns UID=%s created\n", newDns.UID)

		case object.KindAutoScaler:
			var as object.AutoScaler
			err = yaml.Unmarshal(file, &as)
			if err != nil {
				log.Fatal("[FATAL] fail to parse AutoScaler", err)
			}
			newAs, err := crudobj.CreateAutoScaler(as)
			if err != nil {
				log.Fatal("[FATAL] fail to create new AutoScaler")
			}
			log.Printf("AutoScaler UID=%s created\n", newAs.UID)

		case object.KindGpuJob:
			var job object.GpuJob
			err = yaml.Unmarshal(file, &job)
			if err != nil {
				log.Fatal("[FATAL] fail to parse GpuJob", err)
			}
			newJob, err := crudobj.CreateGpuJob(job)
			if err != nil {
				log.Fatal("[FATAL] fail to create new GpuJob")
			}

			err = objfile.PostJobFile(newJob.UID, job.Spec.Filename)
			if err != nil {
				log.Fatal("[FATAL] fail to upload GpuJob file")
			}

			newJob.Status.Phase = object.JobCreated
			newJob, err = crudobj.UpdateGpuJob(newJob)
			if err != nil {
				log.Fatal("[FATAL] fail to update GpuJob phase")
			}

			log.Printf("GpuJob UID=%s created\n", newJob.UID)

		default:
			log.Fatal("[FATAL] Unknown kind: " + t.Kind)
		}
	},
}

func init() {
	rootCmd.AddCommand(createCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")
	createCmd.Flags().StringP("file", "f", "", "path of your config yaml file")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
