package refactor

import (
	"os/exec"
	"path"
	"testing"
)

func SimpleTest(t *testing.T, git_root string) (string, error) {
	testsource := path.Join(git_root, "testdata/simple/helloworldkustomize")
	testtarget := t.TempDir()

	cmd := exec.Command("cp", "-r", testsource+"/.", testtarget)
	testsource = testtarget

	_, err := cmd.CombinedOutput()
	if err != nil {
		t.Errorf("Failed")
	}

	overlays := []string{"dev", "stagig", "prod"}
	err = DoMakeOverlay(testsource, overlays, testsource, "ns", false)
	return testsource, err
}
