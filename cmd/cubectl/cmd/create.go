/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"Cubernetes/pkg/apiserver/crudobj"
	"Cubernetes/pkg/apiserver/objfile"
	"Cubernetes/pkg/object"
	"github.com/google/uuid"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"os"

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
	cubectl create -f action.yaml -s ./myaction.py
	cubectl create -f gpujob.yaml -j ./myjob.tar.gz
	cubectl create -f [file path] [options]`,
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
			// Host path of corresponding gpu job file
			filePath, err := cmd.Flags().GetString("job")
			if err != nil {
				log.Fatal("[FATAL] missing gpu job file")
			}

			if _, err = os.Stat(filePath); err != nil && os.IsNotExist(err) {
				log.Fatal("[FATAL] cannot open gpu job file")
			}

			var job object.GpuJob
			err = yaml.Unmarshal(file, &job)
			if err != nil {
				log.Fatal("[FATAL] fail to parse GpuJob", err)
			}
			newJob, err := crudobj.CreateGpuJob(job)
			if err != nil {
				log.Fatal("[FATAL] fail to create new GpuJob")
			}

			err = objfile.PostJobFile(newJob.UID, filePath)
			if err != nil {
				log.Fatal("[FATAL] fail to upload GpuJob file")
			}

			newJob.Status.Phase = object.JobCreated
			newJob, err = crudobj.UpdateGpuJob(newJob)
			if err != nil {
				log.Fatal("[FATAL] fail to update GpuJob phase")
			}

			log.Printf("GpuJob UID=%s created\n", newJob.UID)

		case object.KindAction:
			scriptPath, err := cmd.Flags().GetString("script")
			if err != nil {
				log.Fatal("[FATAL] missing action script file")
			}

			scriptUID := uuid.New().String()
			err = objfile.PostActionFile(scriptUID, scriptPath)
			if err != nil {
				log.Fatal("[FATAL] fail to upload Action script")
			}

			var action object.Action
			err = yaml.Unmarshal(file, &action)
			if err != nil {
				log.Fatal("[FATAL] fail to parse Action", err)
			}
			action.Spec.ScriptUID = scriptUID

			newAction, err := crudobj.CreateAction(action)
			if err != nil {
				log.Fatal("[FATAL] fail to create new Action, err: ", err)
			}
			log.Printf("Action UID=%s created (or updated)\n", newAction.UID)

		case object.KindIngress:
			var ingress object.Ingress
			err = yaml.Unmarshal(file, &ingress)
			if err != nil {
				log.Fatal("[FATAL] fail to parse Ingress", err)
			}
			newIngress, err := crudobj.CreateIngress(ingress)
			if err != nil {
				log.Fatal("[FATAL] fail to create new Ingress")
			}
			log.Printf("Ingress UID=%s created\n", newIngress.UID)

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
	createCmd.Flags().StringP("script", "s", "", "path of your action script file")
	createCmd.Flags().StringP("job", "j", "", "path of your gpu job file")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
