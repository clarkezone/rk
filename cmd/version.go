package cmd

import (
	"fmt"

	"github.com/clarkezone/rk/pkg/config"
	"github.com/spf13/cobra"
)

//Show shows the current Okteto CLI version
func Version() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show rk CLI version",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("rk version %s hash %s\n", config.VersionHash, config.VersionString)
			return nil
		},
	}
}
