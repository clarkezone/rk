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
			//oktetoLog.SetLevel(logLevel)
			//oktetoLog.SetOutputFormat(outputMode)
			//oktetoLog.Infof("started %s", strings.Join(os.Args, " "))

		},
		PersistentPostRun: func(ccmd *cobra.Command, args []string) {
			log.Printf("done")
			//oktetoLog.Infof("finished %s", strings.Join(os.Args, " "))
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
		//		oktetoLog.Fail(message)
		//		if uErr, ok := err.(oktetoErrors.UserError); ok {
		//			if len(uErr.Hint) > 0 {
		//				oktetoLog.Hint("    %s", uErr.Hint)
		//			}
		//		}
		os.Exit(1)
	}
}
