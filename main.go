package main

import (
	"fmt"
	"log"
	"os"
	"unicode"

	"github.com/spf13/cobra"

	refactorCMD "github.com/clarkezone/rk/cmd/refactor"
)

func main() {
	root := &cobra.Command{
		Use:           fmt.Sprintf("%s COMMAND [ARG...]", "rk"),
		Short:         "Refactoring for kubernetes kustomize manifests",
		Long:          "Refactoring for kubernetes kustomize manifests",
		SilenceErrors: true,
		PersistentPreRun: func(ccmd *cobra.Command, args []string) {
			ccmd.SilenceUsage = false

		},
		PersistentPostRun: func(ccmd *cobra.Command, args []string) {
			log.Printf("done")
		},
	}

	root.AddCommand(refactorCMD.Refactor())

	err := root.Execute()

	if err != nil {
		message := err.Error()
		if len(message) > 0 {
			tmp := []rune(message)
			tmp[0] = unicode.ToUpper(tmp[0])
			message = string(tmp)
		}
		log.Fatal(message)
		os.Exit(1)
	}
}
