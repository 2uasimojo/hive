package terraform

import (
	"fmt"
	"path"

	"github.com/pkg/errors"
)

func terraformExec(clusterDir string, args ...string) error {
	// Create an executor
	ex, err := newExecutor()
	if err != nil {
		return errors.Wrap(err, "failed to create Terraform executor")
	}

	err = ex.execute(clusterDir, args...)
	if err != nil {
		return errors.Wrap(err, "failed to execute Terraform")
	}
	return nil
}

// Apply runs "terraform apply" in the given directory. It returns the absolute
// path of the tfstate file, rooted in the specified directory, along with any
// errors from Terraform.
func Apply(dir string, extraArgs ...string) (string, error) {
	stateFileName := "terraform.tfstate"
	defaultArgs := []string{
		"apply",
		"-auto-approve",
		"-input=false",
		"-no-color",
		fmt.Sprintf("-state=%s", stateFileName),
	}
	args := append(defaultArgs, extraArgs...)

	return path.Join(dir, stateFileName), terraformExec(dir, args...)
}

// Init runs "terraform init" in the given directory.
func Init(dir string) error {
	return terraformExec(dir, "init", "-input=false", "-no-color")
}
