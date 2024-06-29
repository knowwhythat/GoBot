package utils

import (
	"github.com/duke-git/lancet/v2/system"
	"os"
)

func GetVenvPython(path string) string {
	if system.IsWindows() {
		return path + string(os.PathSeparator) + "venv" + string(os.PathSeparator) + "Scripts" + string(os.PathSeparator) + "python.exe"
	} else {
		return path + string(os.PathSeparator) + "venv" + string(os.PathSeparator) + "Scripts" + string(os.PathSeparator) + "python"
	}
}
