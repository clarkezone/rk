package refactor

import (
	"fmt"
	"os"
)

func DoMakeOverlay(sourceDir string, overlayList []string, targetDir string) error {
	// create base dir in target
	// create overlay dir in target
	// foreach over overlays and create a dir for each one in target
	// foreach overlays create a kustomize using golang template replacing
	// move (only) yaml files from source into base
	// add missing kustomize files using golang template
	//  overlay specific details from overlay list

	shouldReturn, returnValue := validateArgs(sourceDir, targetDir, overlayList)
	if shouldReturn {
		return returnValue
	}

	return nil
}

func validateArgs(sourceDir string, targetDir string, overlayList []string) (bool, error) {
	_, err := os.Stat(sourceDir)
	if err != nil {
		return true, err
	}

	_, err = os.Stat(targetDir)
	if err != nil {
		return true, err
	}

	if len(overlayList) == 0 {
		return true, fmt.Errorf("no overlays provided")
	}
	return false, nil
}
