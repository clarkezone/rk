package refactor

import (
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/golang-collections/collections/stack"
	. "sigs.k8s.io/kustomize/kyaml/yaml"
)

func DoMakeOverlay(sourceDir string, overlayList []string, targetDir string, namespace string) error {
	shouldReturn, returnValue := validateArgs(sourceDir, overlayList)

	if shouldReturn {
		return returnValue
	}

	err := expandDir(&sourceDir)
	if err != nil {
		return err
	}

	err = expandDir(&targetDir)
	if err != nil {
		return err
	}

	if !anyManifests(sourceDir) {
		fmt.Printf("No manifests found in source directory %v\n", sourceDir)
		return nil
	}

	if !targetExists(targetDir) {
		err := os.MkdirAll(targetDir, 0755)
		if err != nil {
			return err
		}
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
	pname, err := findPrimaryName(base)
	if err != nil {
		return err
	}

	containernames, err := findContainerNames(base)
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

		err = writeOverlayKustTemplate(thisol, ov, namespace)
		if err != nil {
			return err
		}
		err = writeIncreaseReplicas(thisol, pname)
		if err != nil {
			return err
		}
		err = writeMemoryLimits(thisol, pname, containernames)
		if err != nil {
			return err
		}
	}

	err = writeRootKustTemplate(targetDir, overlayList)
	if err != nil {
		return err
	}

	baseKust := path.Join(targetDir, "base", "kustomization.yaml")

	err = editKustomize(namespace, baseKust)

	return err
}

func expandDir(s *string) error {
	if !path.IsAbs(*s) {
		wd, err := os.Getwd()
		if err != nil {
			return err
		}
		*s = path.Join(wd, *s)
	}
	return nil
}

func validateArgs(sourceDir string, overlayList []string) (bool, error) {
	_, err := os.Stat(sourceDir)
	if err != nil {
		return true, err
	}

	if len(overlayList) == 0 {
		return true, fmt.Errorf("no overlays provided")
	}
	return false, nil
}

func targetExists(targetDir string) bool {
	_, err := os.Stat(targetDir)
	return err == nil
}

type tempargs struct {
	NamePrefix string
	NameSpace  string
}

func writeOverlayKustTemplate(parentPath string, np string, ns string) error {
	manifestPath := path.Join(parentPath, "kustomization.yaml")
	args := tempargs{NamePrefix: np, NameSpace: ns}
	templ := `namePrefix: {{.NamePrefix}}-
namespace: {{.NameSpace}}
commonLabels:
  environment: development
bases:
- ../../base
`
	t := template.Must(template.New("yaml-overlay").Parse(templ))
	err := writeTemplate(manifestPath, t, args)
	return err
}

func writeIncreaseReplicas(parentPath string, pname string) error {
	manifestPath := path.Join(parentPath, "increase_replicas.yaml")
	templ := `apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{.}}
spec:
  replicas: 3
`
	t := template.Must(template.New("yaml-increasereplicas").Parse(templ))
	return writeTemplate(manifestPath, t, pname)
}

type memlimits struct {
	DeploymentName string
	ContainerNames []string
}

func writeMemoryLimits(parentPath string, deploymentName string, containerNames []string) error {
	foo := memlimits{deploymentName, containerNames}
	manifestPath := path.Join(parentPath, "set_memory_limits.yaml")
	templ := `apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{.DeploymentName}}
spec:
  template:
    spec:
      containers:
{{range $val := .ContainerNames}}      - name: {{$val}}
        resources:
          limits:
            memory: 512Mi
{{end}}`
	t := template.Must(template.New("yaml-set-memorylimits").Parse(templ))
	return writeTemplate(manifestPath, t, foo)
}

func writeRootKustTemplate(parentPath string, overlays []string) error {
	manifestPath := path.Join(parentPath, "kustomization.yaml")
	templ := `apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
bases:
{{range $val := .}}- overlay/{{$val}}
{{end}}`
	t := template.Must(template.New("yaml-overlay").Parse(templ))

	return writeTemplate(manifestPath, t, overlays)
}

func writeTemplate(path string, t *template.Template, object interface{}) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}

	err = t.Execute(file, object)
	if err != nil {
		return err
	}
	return nil
}

func anyManifests(baseDir string) bool {
	//TODO recurse
	//TODO verify that any yaml files found are k8s manifests

	dirs, err := ioutil.ReadDir(baseDir)
	if err != nil {
		return false
	}
	for _, e := range dirs {
		ext := strings.ToLower(path.Ext(e.Name()))
		if ext == ".yaml" || ext == ".yml" {
			return true
		}
	}
	return false
}

func findPrimaryName(baseDir string) (string, error) {
	//TODO: iterate over all manifests in base
	// find deployment, daemonset, statefulset
	pName := path.Join(baseDir, "deployment.yaml")
	// lookup name metadata
	return findName(pName)
}

func findContainerNames(baseDir string) ([]string, error) {
	// TODO: find deployment, daemonset, statefulset
	manifestpath := path.Join(baseDir, "deployment.yaml")
	// lookup name metadata
	result, err := findContainerNamesForDeployment(manifestpath)

	return result, err
}

func copyDir(source string, dest string, move bool) error {
	var s *stack.Stack = stack.New()
	err := filepath.Walk(source, func(walkSource string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("Error entering walk %v", err.Error())
		}

		if strings.HasPrefix(walkSource, dest) {
			return nil
		}

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

func loadFile(f string) (*RNode, error) {
	bytes, err := ioutil.ReadFile(f)

	if err != nil {
		return nil, err
	}

	obj, err := Parse(string(bytes))
	if err != nil {
		return nil, err
	}
	return obj, nil
}

func editKustomize(ns string, f string) error {
	bytes, err := ioutil.ReadFile(f)

	if err != nil {
		return err
	}

	obj, err := Parse(string(bytes))
	if err != nil {
		return err
	}
	_, err = obj.Pipe(SetField("namespace", NewScalarRNode(ns)))
	if err != nil {
		return err
	}

	resstr, err := obj.String()
	if err != nil {
		return err
	}

	bytes = []byte(resstr)

	err = ioutil.WriteFile(f, bytes, fs.FileMode(0644))

	return err
}

func findName(f string) (string, error) {
	bytes, err := ioutil.ReadFile(f)

	if err != nil {
		return "", err
	}

	obj, err := Parse(string(bytes))
	if err != nil {
		return "", err
	}
	node, err := obj.Pipe(Get("name"))
	if err != nil {
		return "", err
	}
	value, err := node.String()
	return value, err
}

func findContainerNamesForDeployment(f string) ([]string, error) {
	obj, err := loadFile(f)
	if err != nil {
		return []string{""}, err
	}
	node, err := obj.Pipe(Lookup("spec", "template", "spec", "containers"))
	if err != nil {
		return []string{""}, err
	}
	containerNames, err := node.ElementValues("name")
	if err != nil {
		return []string{""}, err
	}

	return containerNames, err
}
