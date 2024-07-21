package services

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/duke-git/lancet/v2/convertor"
	"github.com/duke-git/lancet/v2/strutil"
	"github.com/spf13/viper"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"gobot/backend/constants"
	"gobot/backend/forms"
	"gobot/backend/log"
	"gobot/backend/services/sys_exec"
	"gobot/backend/utils"
	"os"
	"path/filepath"
)

func StartCapture() (string, error) {
	windowsInspectCommand = sys_exec.BuildCmd(viper.GetString(constants.ConfigPythonPath), "-u", "-B", "-m", "robot_image.screen_shot")
	output, err := windowsInspectCommand.Output()
	if err != nil {
		log.Logger.Logger.Error().Err(err)
		return "", err
	}
	return string(output), err
}

func HighlightImage(id string, imageData forms.ImageData) error {
	project, err := QueryProjectById(id)
	if err != nil {
		return err
	}
	imageDirPath := filepath.Join(project.Path, constants.BaseDir, constants.DevDir, constants.SnapshotDir)
	if !utils.PathExist(imageDirPath) {
		_ = os.MkdirAll(imageDirPath, os.ModePerm)
	}
	imagePath := filepath.Join(imageDirPath, imageData.Id+".png")
	needDelete := false
	if !utils.PathExist(imagePath) {
		needDelete = true
		imageByte, err := base64.StdEncoding.DecodeString(imageData.Image)
		if err != nil {
			return err
		}
		_ = os.WriteFile(filepath.Join(imageDirPath, imageData.Id+".png"), imageByte, 0666)
	}
	windowsInspectCommand = sys_exec.BuildCmd(viper.GetString(constants.ConfigPythonPath), "-u", "-B", "-m", "robot_image.highlight_image",
		imagePath, convertor.ToString(imageData.Threshold), imageData.MatchType)
	output, err := windowsInspectCommand.Output()
	if needDelete {
		_ = os.Remove(imagePath)
	}
	if err != nil {
		log.Logger.Logger.Error().Err(err).Msg("图片匹配失败")
		return errors.New("图片未找到")
	}
	if strutil.Trim(string(output), "\r\n", "\n") == "OK" {
		return nil
	}
	return errors.New("图片未找到")
}

func SaveImage(ctx context.Context, id string, imageData forms.ImageData) error {
	project, err := QueryProjectById(id)
	if err != nil {
		return err
	}
	imageDirPath := filepath.Join(project.Path, constants.BaseDir, constants.DevDir, constants.SnapshotDir)
	if !utils.PathExist(imageDirPath) {
		_ = os.MkdirAll(imageDirPath, os.ModePerm)
	}

	imageByte, err := base64.StdEncoding.DecodeString(imageData.Image)
	if err != nil {
		return err
	}
	_ = os.WriteFile(filepath.Join(imageDirPath, imageData.Id+".png"), imageByte, 0666)
	imageData.Image = filepath.Join(constants.DevDir, constants.SnapshotDir, imageData.Id+".png")
	imageConfigPath := filepath.Join(project.Path, constants.BaseDir, constants.Image)
	selectorData := make(map[string]forms.ImageData)
	if utils.PathExist(imageConfigPath) {
		data, err := os.ReadFile(imageConfigPath)
		if err != nil {
			return err
		}
		err = json.Unmarshal(data, &selectorData)
		if err != nil {
			return err
		}
	}
	selectorData[imageData.Id] = imageData
	saveData, err := json.Marshal(selectorData)
	if err != nil {
		return err
	}
	err = os.WriteFile(imageConfigPath, saveData, 0666)
	runtime.EventsEmit(ctx, constants.ImageElementChange, "")
	return err
}

func RemoveImage(ctx context.Context, id, imageId string) error {
	project, err := QueryProjectById(id)
	if err != nil {
		return err
	}
	imageDirPath := filepath.Join(project.Path, constants.BaseDir, constants.DevDir, constants.SnapshotDir)
	_ = os.Remove(filepath.Join(imageDirPath, imageId+".png"))

	imageConfigPath := filepath.Join(project.Path, constants.BaseDir, constants.Image)
	selectorData := make(map[string]forms.ImageData)
	if utils.PathExist(imageConfigPath) {
		data, err := os.ReadFile(imageConfigPath)
		if err != nil {
			return err
		}
		err = json.Unmarshal(data, &selectorData)
		if err != nil {
			return err
		}
	}
	delete(selectorData, imageId)
	saveData, err := json.Marshal(selectorData)
	if err != nil {
		return err
	}
	err = os.WriteFile(imageConfigPath, saveData, 0666)
	runtime.EventsEmit(ctx, constants.ImageElementChange, "")
	return err
}

func GetImages(id string) (string, error) {
	project, err := QueryProjectById(id)
	if err != nil {
		return "", err
	}
	imageConfigPath := filepath.Join(project.Path, constants.BaseDir, constants.Image)
	imageDatas := make(map[string]forms.ImageData)
	if utils.PathExist(imageConfigPath) {
		data, err := os.ReadFile(imageConfigPath)
		if err != nil {
			return "", err
		}
		err = json.Unmarshal(data, &imageDatas)
		if err != nil {
			return "", err
		}
		datas := make([]forms.ImageData, 0)
		for _, imageData := range imageDatas {
			imagePath := filepath.Join(project.Path, constants.BaseDir, constants.DevDir, constants.SnapshotDir, imageData.Id+".png")
			if utils.PathExist(imagePath) {
				base64Data, err := os.ReadFile(imagePath)
				if err != nil {
					continue
				}
				imageData.Image = base64.StdEncoding.EncodeToString(base64Data)
				datas = append(datas, imageData)
			}
		}
		parsedData, err := json.Marshal(datas)
		if err != nil {
			return "", err
		}
		return string(parsedData), nil
	}

	return "", nil
}
