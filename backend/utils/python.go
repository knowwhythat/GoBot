package utils

import (
	"github.com/duke-git/lancet/v2/system"
	"gobot/backend/constants"
	"path/filepath"
)

func GetVenvPython(path string) string {
	if system.IsWindows() {
		return filepath.Join(path, constants.VenvDir, "Scripts", "python.exe")
	} else {
		return filepath.Join(path, constants.VenvDir, "bin", "python")
	}
}
