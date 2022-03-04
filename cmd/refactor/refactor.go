package refactor

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	refactorCMD "github.com/clarkezone/rk/pkg/cmd/refactor"
)

func Refactor() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "refactor",
		Args:  NoArgsAccepted(),
		Short: "fake command",
		RunE: func(cmd *cobra.Command, args []string) error {
			_ = context.Background()
			var val = refactorCMD.DoRefactor()
			fmt.Printf("RUNNING: Refactor: %v\n", val)
			return nil
		},
	}

	return cmd
}

func NoArgsAccepted() cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 {
			return fmt.Errorf("arguments not supported")
		}
		return nil
	}
}
