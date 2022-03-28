package refactor

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"testing"
)

var git_root string

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	shutdown()
	os.Exit(code)
}

func Test_ValidateArgs(t *testing.T) {
	testsource := path.Join(git_root, "testdata/simple/helloworldkustomize")
	testtarget := t.TempDir()

	err := DoMakeOverlay("", []string{}, "", "ns")
	if err == nil {
		t.Error("Arg validation incorrect")
	}

	err = DoMakeOverlay(testsource, []string{}, "", "ns")
	if err == nil {
		t.Error("Arg validation incorrect")
	}

	err = DoMakeOverlay(testsource, []string{}, testtarget, "ns")
	if err == nil {
		t.Error("Arg validation incorrect")
	}

	//TODO how to make path.Join fail
	badsource := path.Join(testsource, strings.Repeat("ssssssssssssssssssssssssss", 400))
	err = DoMakeOverlay(badsource, []string{}, testtarget, "ns")
	if err == nil {
		t.Error("Arg validation incorrect")
	}
}

func Test_simple(t *testing.T) {
	testsource := path.Join(git_root, "testdata/simple/helloworldkustomize")
	testtarget := t.TempDir()
	correctOutput := path.Join(git_root, "testdata/correctoutput/Test_simple/001")
	overlays := []string{"dev", "stagig", "prod"}
	err := DoMakeOverlay(testsource, overlays, testtarget, "ns")
	if err != nil {
		t.Errorf("Test fail %v", err)
	}

	files, err := ioutil.ReadDir(testtarget)
	if len(files) != 3 {
		t.Errorf("Target file count is wrong")
	}
	if err != nil {
		t.Errorf("Expected error not returned")
	}

	basepath := path.Join(testtarget, "base")
	files, err = ioutil.ReadDir(basepath)
	if len(files) != 3 {
		t.Errorf("Target file count is wrong")
	}
	if err != nil {
		t.Errorf("Expected error not returned")
	}

	overlaypath := path.Join(testtarget, "overlay")
	files, err = ioutil.ReadDir(overlaypath)
	if len(files) != 3 {
		t.Errorf("Target file count is wrong")
	}
	if err != nil {
		t.Errorf("Expected error not returned")
	}

	err = compareTree(testtarget, correctOutput)
	if err != nil {
		t.Errorf("compareTree failed")
	}

	testInput := path.Join(testtarget, "base/kustomization.yaml")
	testOutput := path.Join(correctOutput, "/base/kustomization.yaml")
	//TODO dyfrecurse over correct
	err = dyffFiles(testInput, testOutput)
	if err != nil {
		t.Errorf("Test output doesn't match %v", err)
	}
}

func Test_simple_inplace(t *testing.T) {
	testsource, err := SimpleTest(t, git_root)
	if err != nil {
		t.Errorf("Test fail %v", err)
	}

	files, err := ioutil.ReadDir(testsource)
	if len(files) != 3 {
		t.Errorf("Target file count is wrong")
	}
	if err != nil {
		t.Errorf("Expected error not returned")
	}

	basepath := path.Join(testsource, "base")
	files, err = ioutil.ReadDir(basepath)
	if len(files) != 3 {
		t.Errorf("Target file count is wrong")
	}
	if err != nil {
		t.Errorf("Expected error not returned")
	}

	overlaypath := path.Join(testsource, "overlay")
	files, err = ioutil.ReadDir(overlaypath)
	if len(files) != 3 {
		t.Errorf("Target file count is wrong")
	}
	if err != nil {
		t.Errorf("Expected error not returned")
	}

	//TODO verify kustomize manifests present and correct
}

func dyffFiles(input string, outputPath string) error {
	//TODO: verify success and failure cases
	cmd := exec.Command("dyff", "between", input, outputPath)

	_, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("dyff validation vailed files didn't match %v", err)
	}
	return nil
}

func setup() {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")

	output, err := cmd.CombinedOutput()
	if err != nil {
		panic("couldn't read output from git command get gitroot")
	}
	git_root = string(output)
	git_root = strings.TrimSuffix(git_root, "\n")
}

func compareTree(source string, dest string) error {
	somethingFailed := false
	// run with go test -v -run Test_simple called from the correct dir
	err := filepath.Walk(source, func(walkSource string, info os.FileInfo, err error) error {
		//fmt.Printf("Walk: %v ", walkSource)
		if err != nil {
			log.Printf("Error entering walk %v", err.Error())
		}
		_, err = filepath.Rel(walkSource, dest)
		if err != nil {
			return err
		}

		if !info.IsDir() {
			relPath, err := filepath.Rel(source, walkSource)
			if err != nil {
				return err
			}

			destFull := path.Join(dest, relPath)
			areEqual := dyffFiles(walkSource, destFull)
			fmt.Printf("Good:%v dest %v", areEqual == nil, relPath)
			if areEqual != nil {
				fmt.Printf("\n  dyff between %v %v", walkSource, destFull)
				somethingFailed = true
			}
		}
		fmt.Println("")
		return nil
	})
	if err != nil {
		return err
	}

	if somethingFailed {
		return fmt.Errorf("Dyff failed")
	}

	return nil
}

func Test_compareTree(t *testing.T) {

}

func Test_copyFile(t *testing.T) {
	testsource := path.Join(git_root, "testdata/simple/helloworldkustomize/deployment.yaml")
	testtarget := t.TempDir()
	testtarget = path.Join(testtarget, "deployment.yaml")
	err := copyFile(testsource, testtarget)
	if err != nil {
		t.Errorf("Test fail %v", err)
	}

	sourceStat, _ := os.Stat(testsource)
	destStat, _ := os.Stat(testtarget)
	sourceSize := sourceStat.Size()
	destSize := destStat.Size()

	if sourceSize != destSize {
		t.Errorf("copied file sizes didn't match %v, %v\n", sourceSize, destSize)
	}

	err = copyFile("", "bar")
	if err == nil {
		t.Errorf("Expected error not returned")
	}

	err = copyFile(testsource, "")
	if err == nil {
		t.Errorf("Expected error not returned")
	}
}

func Test_copyDir(t *testing.T) {
	testsource := path.Join(git_root, "testdata/simple/helloworldkustomize")
	testtarget := t.TempDir()
	err := copyDir(testsource, testtarget, false)

	if err != nil {
		t.Errorf("Test fail %v", err)
	}

	files, err := ioutil.ReadDir(testtarget)
	if len(files) != 3 {
		t.Errorf("Target file count is wrong")
	}
	if err != nil {
		t.Errorf("Expected error not returned")
	}

	//TODO.. test scenario with hierarchical source files
	//TODO.. tree walk to ensure structure is correct
}

func Test_copy_inplace(t *testing.T) {
	testsource := path.Join(git_root, "testdata/simple/helloworldkustomize")
	testtarget := t.TempDir()

	cmd := exec.Command("cp", "-r", testsource+"/.", testtarget)

	_, err := cmd.CombinedOutput()
	if err != nil {
		t.Errorf("Failed")
	}

	testsource = testtarget
	testtarget = path.Join(testsource, "moved")

	err = os.Mkdir(testtarget, 0755)
	if err != nil {
		t.Errorf("Expected error not returned")
	}

	err = copyDir(testsource, testtarget, true)
	if err != nil {
		t.Errorf("Expected error not returned")
	}

	files, err := ioutil.ReadDir(testtarget)
	if len(files) != 3 {
		t.Errorf("Target file count is wrong")
	}
	if err != nil {
		t.Errorf("Expected error not returned")
	}

	files, err = ioutil.ReadDir(testsource)
	if len(files) != 1 {
		t.Errorf("Target file count is wrong")
	}
	if err != nil {
		t.Errorf("Expected error not returned")
	}
}

func shutdown() {

}
