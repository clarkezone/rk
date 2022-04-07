package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

//Layer
func Overlay() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "overlay transform, add",
		Short: "restructure source as base and overlays.",
		Args:  NoArgsAccepted(),
	}
	cmd.AddCommand(Create())
	return cmd
}

//Create creates a layer
func Create() *cobra.Command {
	return &cobra.Command{
		Use:   "create",
		Short: "restructure source as base and overlays.",
		Long: `rk overlay create restructures an existing flat set of kubernetes manifests and creates a base with a set of overlays.
	By default this is, dev, stage, prod
	Examples:
		rk overlay create .
		rk overlay create . --out \tmp\thing
		`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("Running overlay create command with source dir of %v\n", args[0])
			return nil
		},
	}
}

func NoArgsAccepted() cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 {
			return fmt.Errorf("no args")
		}
		return nil
	}
}
