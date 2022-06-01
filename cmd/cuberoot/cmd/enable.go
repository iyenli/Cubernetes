package cmd

import (
	"Cubernetes/cmd/cuberoot/utils"
	"Cubernetes/pkg/object"
	"Cubernetes/pkg/utils/localstorage"
	"github.com/spf13/cobra"
	"log"
)

var serverlessCmd = &cobra.Command{
	Use:   "serverless",
	Short: "Control the serverless module of Cubernetes",
	Long: `
Control the serverless module of Cubernetes
usage:
	cuberoot serverless enable
	cuberoot serverless disable`,

	Run: func(cmd *cobra.Command, args []string) {
		meta, err := localstorage.TryLoadMeta()
		if err != nil {
			log.Println("[Error]: register the machine as Cubernetes master and retry")
			return
		}
		if meta.Node.Spec.Type != object.Master {
			log.Println("[Error]: register the machine as Cubernetes master and retry")
			return
		}
		if len(args) != 1 {
			log.Println("[Error]: Too much or little args")
			return
		}

		if args[0] == "enable" {
			err = utils.EnableServerlessGateway(&meta)
			if err != nil {
				log.Println("[Error]: enable serverless gateway failed")
				return
			}
		} else {
			log.Println("[Error]: Unsupported operation to serverless modules")
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(serverlessCmd)
}
