package services

import (
	"encoding/json"
	"errors"
	"gobot/backend/constants"
	"gobot/backend/dao"
	"gobot/backend/models"
	"gobot/backend/services/sys_exec"
	"gobot/backend/utils"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/spf13/viper"
)

func QueryProjectPage() (resultList []*models.Project, err error) {
	if err := dao.ReadTx(dao.MainDB, func(tx dao.Tx) error {
		resultList, err = tx.SelectAllProject()
		return err

	}); err != nil {
		return nil, err
	}

	return resultList, nil
}

func QueryProjectById(id string) (result *models.Project, err error) {
	if err := dao.ReadTx(dao.MainDB, func(tx dao.Tx) error {
		result, err = tx.SelectProject(id)
		return err
	}); err != nil {
		return nil, err
	}
	return result, nil
}

func AddOrUpdateProject(project models.Project) (err error) {
	if project.Id == "" {
		project.Id = uuid.New().String()
		path := viper.GetString(constants.ConfigPythonPath)
		dataPath := viper.GetString(constants.ConfigDataPath)
		projectDir := dataPath + string(os.PathSeparator) + project.Id
		if !utils.PathExist(projectDir) {
			_ = os.MkdirAll(projectDir, fs.ModeDir)
		}
		project.Path = projectDir
		command := sys_exec.BuildCmd(path, "-m", "venv", "--system-site-packages", project.Path+"\\venv")
		_, err := command.Output()
		if err != nil {
			return err
		}
		projectDir = projectDir + string(os.PathSeparator) + constants.BaseDir
		if !utils.PathExist(projectDir) {
			if err = os.MkdirAll(projectDir, fs.ModeDir); err != nil {
				return err
			}
		}
		err = os.WriteFile(projectDir+string(os.PathSeparator)+"__init__.py", []byte("from .import package\n"), 0666)
		if err != nil {
			return err
		}
		err = os.WriteFile(projectDir+string(os.PathSeparator)+"package.py", []byte("variables = {} \n"), 0666)
		if err != nil {
			return err
		}
		if err := dao.WriteTx(dao.MainDB, func(tx dao.Tx) error {
			if err := tx.InsertProject(&project); err != nil {
				return err
			}
			return nil
		}); err != nil {
			return err
		}
	} else {
		result, err := QueryProjectById(project.Id)
		if err != nil {
			return err
		}
		if result == nil {
			return errors.New("未找到原始数据")
		}
		result.Name = project.Name
		result.Description = project.Description
		if project.InputParam != nil {
			result.InputParam = project.InputParam
		}
		_ = ModifyProject(*result)
	}

	return nil
}

func ModifyProject(form models.Project) (err error) {
	if err := dao.WriteTx(dao.MainDB, func(tx dao.Tx) error {
		project, err := tx.SelectProject(form.Id)
		if err != nil {
			return err
		}
		if project == nil {
			return errors.New("project not found")
		}

		if err := tx.UpdateProject(&form); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}

	return nil
}

func RemoveProject(id string) (err error) {
	if err := dao.WriteTx(dao.MainDB, func(tx dao.Tx) error {
		if project, err := tx.SelectProject(id); err != nil {
			return err
		} else if project == nil {
			return errors.New("project not found")
		} else {
			if utils.PathExist(project.Path) {
				if err := os.RemoveAll(project.Path); err != nil {
					return err
				}
			}
			if err = tx.DeleteProject(id); err != nil {
				return err
			}
			if err = tx.DeleteAllSchedulesWhereProjectId(id); err != nil {
				return nil
			}
			if err = tx.DeleteAllExecutionsWhereProjectId(id); err != nil {
				return nil
			}
		}
		return nil
	}); err != nil {
		return err
	}

	return nil
}

func ReadProjectConfig(id string) (models.ProjectConfig, error) {
	project, err := QueryProjectById(id)
	if err != nil {
		return models.ProjectConfig{}, err
	}
	projectPath := project.Path + string(os.PathSeparator) + constants.BaseDir
	devPath := projectPath + string(os.PathSeparator) + constants.DevDir
	if !utils.PathExist(devPath) {
		if err = os.MkdirAll(devPath, fs.ModeDir); err != nil {
			return models.ProjectConfig{}, err
		}
	}
	projectConfigPath := devPath + string(os.PathSeparator) + constants.ProjectConfig
	if !utils.PathExist(projectConfigPath) {
		if project.IsFlow {
			proConfig := models.ProjectConfig{Key: project.Id, Label: project.Name, NodeType: "flow", Children: []*models.Flow{}, Variables: []*models.Variable{}}
			value, err := json.Marshal(&proConfig)
			if err != nil {
				return models.ProjectConfig{}, err
			}
			err = os.WriteFile(projectConfigPath, []byte(value), 0666)
			return proConfig, err
		} else {
			proConfig := models.ProjectConfig{Key: project.Id, Label: project.Name, NodeType: "sequence", Variables: []*models.Variable{}, Children: []*models.Flow{{Key: "main", Label: "主流程", NodeType: "sequence", Opened: true}}}
			value, err := json.Marshal(&proConfig)
			if err != nil {
				return models.ProjectConfig{}, err
			}
			err = os.WriteFile(projectConfigPath, []byte(value), 0666)
			return proConfig, err
		}
	}
	data, err := os.ReadFile(projectConfigPath)
	if err != nil {
		return models.ProjectConfig{}, err
	}
	proConfig := models.ProjectConfig{}
	if err = json.Unmarshal(data, &proConfig); err != nil {
		return proConfig, err
	}
	return proConfig, nil
}

func SaveProjectConfig(id, data string) error {
	project, err := QueryProjectById(id)
	if err != nil {
		return err
	}
	projectPath := project.Path + string(os.PathSeparator) + constants.BaseDir
	devPath := projectPath + string(os.PathSeparator) + constants.DevDir
	if !utils.PathExist(devPath) {
		if err = os.MkdirAll(devPath, fs.ModeDir); err != nil {
			return err
		}
	}
	projectConfigPath := devPath + string(os.PathSeparator) + constants.ProjectConfig
	err = os.WriteFile(projectConfigPath, []byte(data), 0666)
	if err != nil {
		return err
	}
	proConfig := models.ProjectConfig{}
	if err = json.Unmarshal([]byte(data), &proConfig); err != nil {
		return err
	}
	err = GenerateGlobalVariable(proConfig.Variables, projectPath+string(os.PathSeparator)+"package.py")
	return err
}

func SaveProjectDependency(id string, packages []string) error {
	project, err := QueryProjectById(id)
	if err != nil {
		return err
	}
	projectConfigPath := project.Path + string(os.PathSeparator) + constants.BaseDir + string(os.PathSeparator) + constants.DevDir + string(os.PathSeparator) + constants.ProjectConfig
	data, err := os.ReadFile(projectConfigPath)
	if err != nil {
		return err
	}
	proConfig := models.ProjectConfig{}
	if err = json.Unmarshal(data, &proConfig); err != nil {
		return err
	}
	proConfig.ExternalDependencies = packages
	jsonData, err := json.Marshal(proConfig)
	if err != nil {
		return err
	}
	err = os.WriteFile(projectConfigPath, jsonData, 0666)
	return err
}

func ReadMainFlow(id string) (string, error) {
	project, err := QueryProjectById(id)
	if err != nil {
		return "", err
	}
	mainFlowPath := filepath.Join(project.Path, constants.BaseDir, constants.MainFlow)
	if !utils.PathExist(mainFlowPath) {
		return "", nil
	}
	data, err := os.ReadFile(mainFlowPath)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func SaveMainFlow(id string, data string) error {
	project, err := QueryProjectById(id)
	if err != nil {
		return err
	}
	_ = ModifyProject(*project)

	mainFlowPath := filepath.Join(project.Path, constants.BaseDir)
	if !utils.PathExist(mainFlowPath) {
		if err = os.MkdirAll(mainFlowPath, fs.ModeDir); err != nil {
			return err
		}
	}
	err = os.WriteFile(filepath.Join(mainFlowPath, constants.MainFlow), []byte(data), 0666)
	return err
}

func ReadSubFlow(id string, subId string) (string, error) {
	project, err := QueryProjectById(id)
	if err != nil {
		return "", err
	}
	subFlowPath := filepath.Join(project.Path, constants.BaseDir, constants.DevDir, subId+".flow")
	if !utils.PathExist(subFlowPath) {
		return "", nil
	}
	data, err := os.ReadFile(subFlowPath)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func SaveSubFlow(id string, subId string, data string) error {
	project, err := QueryProjectById(id)
	if err != nil {
		return err
	}
	_ = ModifyProject(*project)
	projectPath := filepath.Join(project.Path, constants.BaseDir)
	devPath := filepath.Join(projectPath, constants.DevDir)
	if !utils.PathExist(devPath) {
		if err = os.MkdirAll(devPath, fs.ModeDir); err != nil {
			return err
		}
	}
	err = os.WriteFile(filepath.Join(devPath, subId+".flow"), []byte(data), 0666)
	if err != nil {
		return err
	}
	flow := Flow{}
	err = json.Unmarshal([]byte(data), &flow)
	if err != nil {
		return nil
	}
	err = flow.GeneratePythonCode(filepath.Join(projectPath, subId+".py"))
	return err
}

func DeleteSubFlow(id string, subId string) error {
	project, err := QueryProjectById(id)
	if err != nil {
		return err
	}
	projectPath := filepath.Join(project.Path, constants.BaseDir)
	devPath := filepath.Join(projectPath, constants.DevDir)
	if utils.PathExist(filepath.Join(devPath, subId+".flow")) {
		err = os.Remove(filepath.Join(devPath, subId+".flow"))
		if err != nil {
			return err
		}
	}
	if utils.PathExist(filepath.Join(projectPath, subId+".py")) {
		err = os.Remove(filepath.Join(projectPath, subId+".py"))
		if err != nil {
			return err
		}
	}
	return nil
}

func ReadSubModule(id string, subId string) (string, error) {
	project, err := QueryProjectById(id)
	if err != nil {
		return "", err
	}
	code := ""
	subModuleCodePath := filepath.Join(project.Path, constants.BaseDir, subId+".py")
	if utils.PathExist(subModuleCodePath) {
		data, err := os.ReadFile(subModuleCodePath)
		if err != nil {
			return "", err
		}
		code = string(data)
	}
	subModulePath := filepath.Join(project.Path, constants.BaseDir, constants.DevDir, subId+".flow")
	if !utils.PathExist(subModulePath) {
		return "", nil
	}
	data, err := os.ReadFile(subModulePath)
	if err != nil {
		return "", err
	}
	moduleConfig := make(map[string]interface{})
	if err = json.Unmarshal(data, &moduleConfig); err != nil {
		return "", err
	}
	moduleConfig["code"] = code
	jsonData, err := json.Marshal(moduleConfig)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

func SaveSubModule(id string, subId string, data string) error {
	project, err := QueryProjectById(id)
	if err != nil {
		return err
	}
	_ = ModifyProject(*project)
	projectPath := filepath.Join(project.Path, constants.BaseDir)
	devPath := filepath.Join(projectPath, constants.DevDir)
	if !utils.PathExist(devPath) {
		if err = os.MkdirAll(devPath, fs.ModeDir); err != nil {
			return err
		}
	}
	moduleConfig := make(map[string]interface{})
	if err = json.Unmarshal([]byte(data), &moduleConfig); err != nil {
		return err
	}

	code := moduleConfig["code"].(string)
	delete(moduleConfig, "code")
	subModuleCodePath := filepath.Join(project.Path, constants.BaseDir, subId+".py")
	err = os.WriteFile(filepath.Join(subModuleCodePath), []byte(code), 0666)
	if err != nil {
		return err
	}
	path := viper.GetString(constants.ConfigPythonPath)
	command := sys_exec.BuildCmd(path, "-m", "robot_core.parse_function_definitions", subModuleCodePath)
	funcDefine, err := command.Output()
	if err != nil {
		return err
	}
	funcs := make([]interface{}, 0)
	if err = json.Unmarshal(funcDefine, &funcs); err != nil {
		return err
	}
	moduleConfig["funcDefine"] = funcs
	jsonData, err := json.Marshal(moduleConfig)
	if err != nil {
		return err
	}
	err = os.WriteFile(filepath.Join(devPath, subId+".flow"), jsonData, 0666)

	return err
}

func DeleteSubModule(id string, subId string) error {
	project, err := QueryProjectById(id)
	if err != nil {
		return err
	}
	projectPath := filepath.Join(project.Path, constants.BaseDir)
	devPath := filepath.Join(projectPath, constants.DevDir)
	if utils.PathExist(filepath.Join(devPath, subId+".flow")) {
		err = os.Remove(filepath.Join(devPath, subId+".flow"))
		if err != nil {
			return err
		}
	}
	if utils.PathExist(filepath.Join(projectPath, subId+".py")) {
		err = os.Remove(filepath.Join(projectPath, subId+".py"))
		if err != nil {
			return err
		}
	}
	return nil
}
