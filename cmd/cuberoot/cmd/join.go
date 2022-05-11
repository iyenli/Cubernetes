package cmd

import (
	"Cubernetes/cmd/cuberoot/utils"
	"Cubernetes/pkg/cubenetwork/nodenetwork"
	"Cubernetes/pkg/object"
	"Cubernetes/pkg/utils/localstorage"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"net"
)

// joinCmd represents join master as a slave
var joinCmd = &cobra.Command{
	Use:   "join",
	Short: "Join exist API Server as a slave",
	Long: `
Join an existed master as a slave
usage:
	cuberoot join [Master IP] -f [file path]
example:
	cuberoot join 192.168.1.11 -f node.yaml`,

	Run: func(cmd *cobra.Command, args []string) {
		meta, err := localstorage.TryLoadMeta()
		if err == nil {
			if meta.Node.Spec.Type == object.Master {
				log.Fatal("[FATAL] already initialized as master, please reset first")
			} else {
				log.Fatalf("[FATAL] already joined %s as slave, please reset first", meta.MasterIP)
			}
		}

		if len(args) < 1 {
			log.Fatal("[FATAL] lack arguments")
		}
		if net.ParseIP(args[0]) == nil {
			log.Fatalf("[FATAL] illegal ip address: %v", args[0])
		}
		masterIP := args[0]

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

		nodenetwork.SetMasterIP(masterIP)

		log.Println("Registering as slave...")
		err = utils.RegisterAsSlave(node, masterIP)
		if err != nil {
			log.Fatal("[FATAL] fail to register as slave, err: ", err)
		}

		log.Println("Registered as slave, starting processes...")

		err = utils.StartSlave(node.Status.Addresses.InternalIP, masterIP, node.UID)
		if err != nil {
			log.Fatal("[FATAL] fail to start slave processes, err: ", err)
		}
		log.Println("Slave node launched successfully")
	},
}

func init() {
	rootCmd.AddCommand(joinCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	joinCmd.Flags().StringP("file", "f", "", "path of your node config yaml file")
}
