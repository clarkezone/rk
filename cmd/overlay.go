package cmd

import (
	"fmt"

	refactorCMD "github.com/clarkezone/rk/pkg/cmd/refactor"
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
	var outdir string
	var namespace string
	command := &cobra.Command{
		Use:   "create <souredir>",
		Short: "restructure source as base and overlays.",
		Long: `rk overlay create restructures an existing flat set of kubernetes manifests and creates a base with a set of overlays.
	By default this is, dev, stage, prod
	Examples:
		rk overlay create .
		rk overlay create . --out \tmp\thing
		`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return refactorCMD.DoMakeOverlay(
				args[0],
				[]string{"dev", "staging", "prod"},
				outdir,
				namespace,
				true)
		},
	}
	command.Flags().StringVar(&outdir, "out", "output", "Specify an output directory")
	err := command.MarkFlagDirname("out")
	if err != nil {
		panic(err)
	}
	command.Flags().StringVar(&namespace, "namespace", "default", "Specify kubernetes namespace")
	return command
}

func NoArgsAccepted() cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 {
			return fmt.Errorf("no args")
		}
		return nil
	}
}
