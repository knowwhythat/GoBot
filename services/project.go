package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"gobot/constants"
	"gobot/dao"
	"gobot/forms"
	"gobot/models"
	"gobot/services/sys_exec"
	"gobot/utils"
	"io/fs"
	"os"
	"strings"

	"github.com/ahmetb/go-linq"
	"github.com/spf13/viper"
)

func QueryProjectPage(page forms.QueryForm) (total int, resultList []*models.Project, err error) {
	if err := dao.ReadTx(dao.MainDB, func(tx dao.Tx) error {
		allProjects, err := tx.SelectAllProject()
		if err != nil {
			return err
		}

		query := linq.From(allProjects).Where(func(project interface{}) bool {
			projectVlue := project.(*models.Project)
			return strings.Contains(projectVlue.Name, page.Query)
		}).OrderByDescendingT(func(project *models.Project) int64 {
			return project.UpdateTs.UnixNano()
		}).Query
		total = query.Count()

		if page.PageNum > 0 && page.PageSize > 0 {
			start, end := page.Range()
			query = query.Skip(start).Take(end - start)
		}
		projects := query.Results()

		for _, project := range projects {
			project := project.(*models.Project)
			resultList = append(resultList, project)
		}
		return nil
	}); err != nil {
		return 0, nil, err
	}

	return total, resultList, nil
}

func CountProject() (total int, err error) {
	if err := dao.ReadTx(dao.MainDB, func(tx dao.Tx) error {
		allProjects, err := tx.SelectAllProject()
		if err != nil {
			return err
		}
		total = len(allProjects)
		return nil
	}); err != nil {
		return 0, err
	}
	return total, nil
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

func AddOrUpdateProject(id, name, desc string, isFlow bool) (result *models.Project, err error) {
	path := viper.GetString("python.path")
	dataPath := viper.GetString("data.path")
	projectDir := dataPath + constants.BaseDir + string(os.PathSeparator) + id
	if !utils.PathExist(projectDir) {
		os.MkdirAll(projectDir, fs.ModeDir)
	}
	if id == "" {
		result = &models.Project{
			Name:        name,
			Description: desc,
			IsFlow:      isFlow,
			Path:        projectDir,
		}
		command := sys_exec.BuildCmd(path, "-m", "venv", result.Path+"\\venv")
		output, err := command.Output()
		if err != nil {
			fmt.Println(string(output))
			return nil, err
		}
		if err := dao.WriteTx(dao.MainDB, func(tx dao.Tx) error {
			if err := tx.InsertProject(result); err != nil {
				return err
			}
			return nil
		}); err != nil {
			return nil, err
		}
	} else {
		project, err := QueryProjectById(id)
		if err != nil {
			return nil, err
		}
		project.Name = name
		project.Description = desc
		project.IsFlow = isFlow
		ModifyProject(*project)
	}

	return result, nil
}

func ModifyProject(form models.Project) (err error) {
	if err := dao.WriteTx(dao.MainDB, func(tx dao.Tx) error {
		project, err := tx.SelectProject(form.Id.String())
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
		}
		return nil
	}); err != nil {
		return err
	}

	return nil
}

func ReadMainFlow(id string) (string, string, error) {
	project, err := QueryProjectById(id)
	if err != nil {
		return "", "", err
	}
	mainFlowPath := project.Path + string(os.PathSeparator) + constants.BaseDir + string(os.PathSeparator) + constants.MainFlow
	if !utils.PathExist(mainFlowPath) {
		return "", "", nil
	}
	data, err := os.ReadFile(mainFlowPath)
	if err != nil {
		return "", "", err
	}
	return project.Name, string(data), nil
}

func SaveMainFlow(id string, data string) error {
	project, err := QueryProjectById(id)
	if err != nil {
		return err
	}
	ModifyProject(*project)

	mainFlowPath := project.Path + string(os.PathSeparator) + constants.BaseDir
	if !utils.PathExist(mainFlowPath) {
		if err = os.MkdirAll(mainFlowPath, fs.ModeDir); err != nil {
			return err
		}
	}
	err = os.WriteFile(mainFlowPath+string(os.PathSeparator)+constants.MainFlow, []byte(data), 0666)
	return err
}

func ReadSubFlow(id string, subId string) (string, error) {
	project, err := QueryProjectById(id)
	if err != nil {
		return "", err
	}
	subFlowPath := project.Path + string(os.PathSeparator) + constants.BaseDir + string(os.PathSeparator) + constants.DevDir + string(os.PathSeparator) + subId + ".flow"
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
	ModifyProject(*project)
	projectPath := project.Path + string(os.PathSeparator) + constants.BaseDir
	devPath := projectPath + string(os.PathSeparator) + constants.DevDir
	if !utils.PathExist(devPath) {
		if err = os.MkdirAll(devPath, fs.ModeDir); err != nil {
			return err
		}
	}
	err = os.WriteFile(devPath+string(os.PathSeparator)+subId+".flow", []byte(data), 0666)
	if err != nil {
		return err
	}
	flow := Flow{}
	err = json.Unmarshal([]byte(data), &flow)
	if err != nil {
		return nil
	}
	err = flow.GeneratePythonCode(projectPath + string(os.PathSeparator) + subId + ".py")
	return err
}

func DeleteSubFlow(id string, subId string) error {
	project, err := QueryProjectById(id)
	if err != nil {
		return err
	}
	projectPath := project.Path + string(os.PathSeparator) + constants.BaseDir
	devPath := projectPath + string(os.PathSeparator) + constants.DevDir
	if utils.PathExist(devPath + string(os.PathSeparator) + subId + ".flow") {
		err = os.Remove(devPath + string(os.PathSeparator) + subId + ".flow")
		if err != nil {
			return err
		}
	}
	if utils.PathExist(projectPath + string(os.PathSeparator) + subId + ".py") {
		err = os.Remove(projectPath + string(os.PathSeparator) + subId + ".py")
		if err != nil {
			return err
		}
	}
	return nil
}
