package refactor

import (
	"fmt"
	"html/template"
	"os"
	"path"
)

func DoMakeOverlay(sourceDir string, overlayList []string, targetDir string) error {

	shouldReturn, returnValue := validateArgs(sourceDir, targetDir, overlayList)
	if shouldReturn {
		return returnValue
	}

	// create base dir in target
	base := path.Join(targetDir, "base")
	err := os.MkdirAll(base, 0755)
	if err != nil {
		return err
	}

	// create overlay dir in target
	// foreach over overlays and create a dir for each one in target
	// foreach overlays create a kustomize using golang template replacing
	for _, ov := range overlayList {
		thisol := path.Join(targetDir, "overlay", ov)
		err = os.MkdirAll(thisol, 0755)

		if err != nil {
			return err
		}

		manifest := path.Join(thisol, "kustomization.yaml")
		doTemplate(manifest, ov, "gitea")
	}

	// move (only) yaml files from source into base
	// add missing kustomize files using golang template
	//  overlay specific details from overlay list

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

type tempargs struct {
	NamePrefix string
	NameSpace  string
}

func doTemplate(path string, np string, ns string) error {
	args := tempargs{NamePrefix: np, NameSpace: ns}
	templ := `namePrefix: {{.NamePrefix}}-
namespace: {{.NameSpace}}
commonLabels:
  environment: development
bases:
- ../../base
`
	t := template.Must(template.New("yaml-overlay").Parse(templ))

	file, err := os.Create(path)
	if err != nil {
		return err
	}

	err = t.Execute(file, args)
	if err != nil {
		return err
	}
	return nil
}
