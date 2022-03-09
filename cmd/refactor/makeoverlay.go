package refactor

import (
	"fmt"

	refactorCMD "github.com/clarkezone/rk/pkg/cmd/refactor"
	"github.com/spf13/cobra"
)

func MakeOverlay() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "makeoverlay",
		Args:  GetOverlayList(),
		Short: "m",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("dd")
			return refactorCMD.DoMakeOverlay(
				"/home/james/src/github.com/clarkezone/JekyllPreview/k8s/full",
				[]string{"dev", "staging", "prod"},
				"/home/james/src/github.com/clarkezone/JekyllPreview/k8s/layered",
				"JekyllPreviewV2")
			//	getCurrentDir())
			//			return refactorCMD.DoMakeOverlay(
			//				getCurrentDir(),
			//				[]string{"dev", "staging", "prod"},
			//				getCurrentDir())
		},
	}
	return cmd
}

//func getCurrentDir() string {
//	return "/usr/james/"
//}

func GetOverlayList() cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {

		return nil
	}
}
