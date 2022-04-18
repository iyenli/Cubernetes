package app

import (
	"github.com/spf13/cobra"
)

const (
	componentCubelet = "cubelet"
)

func NewCubeletCommand() *cobra.Command {

	return &cobra.Command{
		Use: componentCubelet,
		// TODO: parse cubelet command
	}
}
