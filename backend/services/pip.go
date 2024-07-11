package services

import (
	"bufio"
	"context"
	"github.com/duke-git/lancet/v2/slice"
	"github.com/spf13/viper"
	"gobot/backend/constants"
	"gobot/backend/forms"
	"gobot/backend/services/sys_exec"
	"gobot/backend/utils"
	"io"
	"regexp"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func InstallPackage(ctx context.Context, id string, pip forms.PipPackage) error {
	project, err := QueryProjectById(id)
	if err != nil {
		return nil
	}
	pythonPath := utils.GetVenvPython(project.Path)
	var stderr io.ReadCloser
	var stdout io.ReadCloser
	if pip.UpgradePip {
		proc := sys_exec.BuildCmd(pythonPath, "-m", "pip", "install", "--upgrade", "pip")
		stdout, err = proc.StdoutPipe()
		if err != nil {
			return err
		}
		stderr, err = proc.StderrPipe()
		if err != nil {
			return err
		}
		go monitorOutput(ctx, stderr)
		go monitorOutput(ctx, stdout)
		err = proc.Run()
		if err != nil {
			return err
		}
	}
	var installArgs []string
	installArgs = append(installArgs, pythonPath)
	installArgs = append(installArgs, "-m")
	installArgs = append(installArgs, "pip")
	installArgs = append(installArgs, "install")

	if pip.Upgrade {
		installArgs = append(installArgs, "--upgrade")
	}

	if pip.InputVersion && pip.Version != "" {
		installArgs = append(installArgs, pip.Name+"=="+pip.Version)
	} else {
		installArgs = append(installArgs, pip.Name)
	}

	if pip.InputMirror && pip.MirrorUrl != "" {
		installArgs = append(installArgs, "i")
		installArgs = append(installArgs, pip.MirrorUrl)
	}
	proc := sys_exec.BuildCmd(installArgs...)
	stdout, err = proc.StdoutPipe()
	if err != nil {
		return err
	}
	stderr, err = proc.StderrPipe()
	if err != nil {
		return err
	}
	go monitorOutput(ctx, stderr)
	go monitorOutput(ctx, stdout)
	err = proc.Run()

	return err
}

func ListInstallPackage(id string, onlyProject bool) ([]forms.Package, error) {
	project, err := QueryProjectById(id)
	if err != nil {
		return nil, nil
	}
	pythonPath := utils.GetVenvPython(project.Path)
	proc := sys_exec.BuildCmd(pythonPath, "-m", "pip", "list")
	output, err := proc.Output()
	if err != nil {
		return nil, err
	}
	allPackages := extractPackages(string(output))
	if onlyProject {
		proc = sys_exec.BuildCmd(viper.GetString(constants.ConfigPythonPath), "-m", "pip", "list")
		output, err = proc.Output()
		if err != nil {
			return nil, err
		}
		basePackages := extractPackages(string(output))
		allPackages = slice.DifferenceWith(allPackages, basePackages, func(item1, item2 forms.Package) bool {
			return strings.ReplaceAll(item1.Name, "-", "_") == strings.ReplaceAll(item2.Name, "-", "_")
		})
	}

	return allPackages, err
}

func UnInstallPackage(id, name string) error {
	project, err := QueryProjectById(id)
	if err != nil {
		return nil
	}
	pythonPath := utils.GetVenvPython(project.Path)
	proc := sys_exec.BuildCmd(pythonPath, "-m", "pip", "uninstall", "-y", name)
	return proc.Run()
}

func extractPackages(lines string) []forms.Package {
	var packages []forms.Package
	reg := regexp.MustCompile(`\n|\r\n`)
	result := reg.Split(lines, -1)
	for i := 2; i < len(result); i++ {
		words := strings.Fields(result[i])
		if len(words) >= 2 {
			packages = append(packages, forms.Package{Name: words[0], Version: words[1]})
		}
	}
	return packages
}

func monitorOutput(ctx context.Context, outPipe io.ReadCloser) {
	reader := bufio.NewReader(outPipe)
	for {
		line, err := reader.ReadString('\n')
		if err != nil || io.EOF == err {
			break
		}
		runtime.EventsEmit(ctx, "pipOutput", line)
	}
}
