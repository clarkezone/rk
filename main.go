package main

import (
	refactorCMD "github.com/clarkezone/rk/cmd/refactor"
)

func main() {
	refactorCMD.MakeOverlay().Execute()

	//	root := &cobra.Command{
	//		Use:           fmt.Sprintf("%s COMMAND [ARG...]", "rk"),
	//		Short:         "Refactoring for kubernetes kustomize manifests",
	//		Long:          "Refactoring for kubernetes kustomize manifests",
	//		SilenceErrors: true,
	//		PersistentPreRun: func(ccmd *cobra.Command, args []string) {
	//			ccmd.SilenceUsage = false
	//
	//		},
	//		PersistentPostRun: func(ccmd *cobra.Command, args []string) {
	//			log.Printf("done")
	//		},
	//	}
	//
	//	root.AddCommand(refactorCMD.Refactor())
	//	root.AddCommand(refactorCMD.MakeOverlay())
	//
	//	err := root.Execute()
	//
	//	if err != nil {
	//		message := err.Error()
	//		if len(message) > 0 {
	//			tmp := []rune(message)
	//			tmp[0] = unicode.ToUpper(tmp[0])
	//			message = string(tmp)
	//		}
	//		log.Fatal(message)
	//		os.Exit(1)
	//	}
}
