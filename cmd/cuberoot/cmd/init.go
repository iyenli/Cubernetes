package cmd

import (
	"Cubernetes/cmd/cuberoot/utils"
	"Cubernetes/pkg/object"
	"Cubernetes/pkg/utils/localstorage"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"net"
	"time"
)

// initCmd represents the init cubernetes master
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Init an API Server as a cubernetes master",
	Long: `
Init as a cubernetes master
usage:
	cuberoot init -f [file path]
example:
	cuberoot init -f node.yaml`,

	Run: func(cmd *cobra.Command, args []string) {
		meta, err := localstorage.TryLoadMeta()
		if err == nil {
			if meta.Node.Spec.Type == object.Master {
				log.Fatal("[FATAL] already initialized as master, please reset first")
			} else {
				log.Fatalf("[FATAL] already joined %s as slave, please reset first", meta.MasterIP)
			}
		}

		f, err := cmd.Flags().GetString("file")
		if err != nil {
			log.Fatal("[FATAL] missing input config file")
		}
		file, err := ioutil.ReadFile(f)
		if err != nil {
			log.Fatal("[FATAL] cannot read input config file")
		}

		var node object.Node
		err = yaml.Unmarshal(file, &node)
		if err != nil {
			log.Fatal("[FATAL] fail to parse config file")
		}

		if net.ParseIP(node.Status.Addresses.InternalIP) == nil {
			log.Fatalf("[FATAL] illegal ip address: %v", node.Status.Addresses.InternalIP)
		}

		log.Println("Starting etcd & apiserver processes, this may take 4s")

		err = utils.PreStartMaster()
		if err != nil {
			log.Fatal("[FATAL] fail to pre-start master processes, err: ", err)
		}

		time.Sleep(4 * time.Second)

		log.Println("Registering as master...")
		err = utils.RegisterAsMaster(node)
		if err != nil {
			log.Fatal("[FATAL] fail to register as master, err: ", err)
		}

		time.Sleep(3 * time.Second)
		log.Println("Registered as master, starting processes...")

		meta, err = localstorage.TryLoadMeta()
		if err != nil {
			log.Fatal("[Fatal]: Meta file should have existed")
		}

		log.Printf("Starting Master, UID = %v, It may takes 15s...", meta.Node.UID)
		err = utils.StartMaster(node.Status.Addresses.InternalIP, meta.Node.UID)
		if err != nil {
			log.Fatal("[FATAL] fail to start master processes, err: ", err)
		}

		log.Printf("Master node launched successfully\n"+
			"To join Cubernetes cluster, execute:\n"+
			"\tcuberoot join %s -f [node config file]\n", node.Status.Addresses.InternalIP)
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
	initCmd.Flags().StringP("file", "f", "", "path of your node config yaml file")
}
