package app

import "github.com/spf13/cobra"

const (
	componentCubeproxy = "cubeproxy"
)

func NewCubeletCommand() *cobra.Command {
	return &cobra.Command{
		Use: componentCubeproxy,
		// TODO: parse cubelet command
	}
}
