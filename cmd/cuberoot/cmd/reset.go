package cmd

import (
	"Cubernetes/cmd/cuberoot/utils"
	"Cubernetes/pkg/apiserver/crudobj"
	"Cubernetes/pkg/cubenetwork/nodenetwork"
	"Cubernetes/pkg/cubeproxy/proxyruntime"
	"Cubernetes/pkg/object"
	"Cubernetes/pkg/utils/kafka"
	"Cubernetes/pkg/utils/localstorage"
	"github.com/spf13/cobra"
	"log"
)

// resetCmd represents reset a node
var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset an registered node",
	Long: `
Reset an registered node, clear local metadata
usage:
	cuberoot reset`,

	Run: func(cmd *cobra.Command, args []string) {
		meta, err := localstorage.TryLoadMeta()
		if err != nil {
			_ = localstorage.ClearMeta()
			return
		}

		if meta.Node.Spec.Type == object.Slave {
			// TODO: What if master is reset too
			nodenetwork.SetMasterIP(meta.MasterIP)
			err = crudobj.DeleteNode(meta.Node.UID)
			if err != nil {
				log.Println("[INFO] fail to delete node from apiserver, err: ", err)
				log.Println("[INFO] Reset before you stop the master")
			}
		} else {
			// only master clear all topics
			err = kafka.DeleteAllTopics(meta.Node.Status.Addresses.InternalIP)
			if err != nil {
				log.Println("[Error]: fail to clean kafka topics and msg")
			}

			err = utils.ClearData()
			if err != nil {
				log.Println("[INFO] fail to remove etcd data, err: ", err)
				log.Println("[INFO] Reset before you stop the master")
				return
			}
		}

		err = localstorage.ClearMeta()
		if err != nil {
			log.Println("[FATAL] fail to clear local metadata, err: ", err)
		}

		err = proxyruntime.CleanIptables()
		if err != nil {
			log.Println("[Error]: fail to clean iptables chain and rules")
		}
	},
}

func init() {
	rootCmd.AddCommand(resetCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
