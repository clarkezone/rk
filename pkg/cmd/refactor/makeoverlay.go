package refactor

import (
	"fmt"
	"html/template"
	"io"
	"os"
	"path"
	"path/filepath"

	"github.com/golang-collections/collections/stack"
)

func DoMakeOverlay(sourceDir string, overlayList []string, targetDir string) error {

	shouldReturn, returnValue := validateArgs(sourceDir, targetDir, overlayList)
	if shouldReturn {
		return returnValue
	}

	match, err := filepath.Rel(sourceDir, targetDir)
	if err != nil {
		return err
	}

	var inplace = false
	if match == "." {
		inplace = true
	}

	// create base dir in target
	base := path.Join(targetDir, "base")
	err = os.MkdirAll(base, 0755)
	if err != nil {
		return err
	}

	// move (only) yaml files from source into base
	err = copyDir(sourceDir, base, inplace)
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
		err = writeOverlayKustTemplate(manifest, ov, "gitea")
		if err != nil {
			return err
		}
	}

	rootKustomize := path.Join(targetDir, "kustomization.yaml")
	err = writeRootKustTemplate(rootKustomize, overlayList)
	if err != nil {
		return err
	}

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

func writeOverlayKustTemplate(path string, np string, ns string) error {
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

func writeRootKustTemplate(path string, overlays []string) error {
	templ := `apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
bases:
{{range $val := .}}- overlay/{{$val}}
{{end}}`
	t := template.Must(template.New("yaml-overlay").Parse(templ))

	file, err := os.Create(path)
	if err != nil {
		return err
	}

	err = t.Execute(file, overlays)
	if err != nil {
		return err
	}
	return nil
}

func copyDir(source string, dest string, move bool) error {
	var s *stack.Stack = stack.New()
	err := filepath.Walk(source, func(walkSource string, info os.FileInfo, err error) error {
		isMovedPath, err := filepath.Rel(walkSource, dest)
		if err != nil {
			return err
		}

		// When moving, skip anything that is in the destination directory
		if move && isMovedPath == ".." {
			return nil
		}
		if !info.IsDir() {
			relPath, err := filepath.Rel(source, walkSource)
			if err != nil {
				return err
			}

			walkTarget := path.Join(dest, relPath)
			err = copyFile(walkSource, walkTarget)
			if err != nil {
				return err
			}

			if move {
				s.Push(walkSource)
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	//Delete items
	for s.Len() > 0 {
		delete := s.Pop()
		str := fmt.Sprintf("%v", delete)
		err = os.Remove(string(str))
		if err != nil {
			return err
		}
	}
	return nil
}

func copyFile(source string, dest string) error {
	from, err := os.Open(source)
	if err != nil {
		return err
	}
	defer from.Close()

	//destNamedPath := path.Base(source)
	//destNamedPath = path.Join(dest, destNamedPath)

	to, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer to.Close()

	_, err = io.Copy(to, from)

	return err
}
