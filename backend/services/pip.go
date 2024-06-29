package services

import (
	"bufio"
	"context"
	"gobot/backend/forms"
	"gobot/backend/services/sys_exec"
	"gobot/backend/utils"
	"io"

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

func ListInstallPackage(id string) (string, error) {
	project, err := QueryProjectById(id)
	if err != nil {
		return "", nil
	}
	pythonPath := utils.GetVenvPython(project.Path)
	proc := sys_exec.BuildCmd(pythonPath, "-m", "pip", "list")
	output, err := proc.Output()
	return string(output), err
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

func monitorOutput(ctx context.Context, outPipe io.ReadCloser) {
	reader := bufio.NewReader(outPipe)
	for {
		line, err := reader.ReadString(')')
		if err != nil || io.EOF == err {
			break
		}
		runtime.EventsEmit(ctx, "pipOutput", line)
	}
}
