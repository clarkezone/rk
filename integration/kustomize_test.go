package integration

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"

	refactorCMD "github.com/clarkezone/rk/pkg/cmd/refactor"
)

var git_root string

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
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

func Test_base(t *testing.T) {
	testsource, err := refactorCMD.SimpleTest(t, git_root)
	if err != nil {
		t.Errorf("SimpleTest failed as part of integration run %v\n", err)
	}
	err = invokeKubectl(testsource)
	if err != nil {
		t.Errorf("Kubectl kustomize failed: %v\n", err)
	}
}

func invokeKubectl(d string) error {
	cmd := exec.Command("kubectl", "kustomize", "base")
	cmd.Dir = d

	yaml, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("Kustomize failed %v\n", err.Error())
	}

	const expected = 1079
	var actual = len(yaml)
	if actual != expected {
		return fmt.Errorf("yaml had incorrect length: expected:%v actual %v \n %v",
			expected,
			actual,
			yaml)
	}

	return nil
}
