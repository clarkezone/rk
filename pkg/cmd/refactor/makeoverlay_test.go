package refactor

import (
	"os"
	"os/exec"
	"path"
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

	err := DoMakeOverlay("", []string{}, "")
	if err == nil {
		t.Error("Arg validation incorrect")
	}

	err = DoMakeOverlay(testsource, []string{}, "")
	if err == nil {
		t.Error("Arg validation incorrect")
	}

	err = DoMakeOverlay(testsource, []string{}, testtarget)
	if err == nil {
		t.Error("Arg validation incorrect")
	}
}

func Test_simple(t *testing.T) {
	testsource := path.Join(git_root, "testdata/simple/helloworldkustomize")
	testtarget := t.TempDir()
	overlays := []string{"dev", "stagig", "prod"}
	err := DoMakeOverlay(testsource, overlays, testtarget)
	if err != nil {
		t.Errorf("Test fail %v", err)
	}
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

func shutdown() {

}
