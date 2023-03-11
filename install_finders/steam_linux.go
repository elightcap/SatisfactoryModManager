package install_finders

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

func FindInstallationsLinuxSteam() ([]*Installation, []error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, []error{errors.Wrap(err, "failed to get user home dir")}
	}
	steamPath := filepath.Join(homeDir, ".steam", "steam")
	if _, err := os.Stat(steamPath); os.IsNotExist(err) {
		return nil, []error{errors.New("steam not installed")}
	}
	return findInstallationsSteam(
		steamPath,
		"Steam",
		[]string{
			"cmd",
			"/C",
			"start",
			"",
		},
	)
}