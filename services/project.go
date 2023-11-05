package services

import (
	"errors"
	"fmt"
	"gobot/dao"
	"gobot/forms"
	"gobot/models"
	"gobot/utils"
	"io/fs"
	"os"
	"os/exec"
	"os/user"
	"strings"

	. "github.com/ahmetb/go-linq"
	uuid "github.com/google/uuid"
	"github.com/spf13/viper"
)

func QueryProjectPage(page forms.QueryForm) (total int, resultList []*models.Project, err error) {
	if err := dao.ReadTx(dao.MainDB, func(tx dao.Tx) error {
		allProjects, err := tx.SelectAllProject()
		if err != nil {
			return err
		}

		query := From(allProjects).Where(func(project interface{}) bool {
			projectVlue := project.(*models.Project)
			return strings.Contains(projectVlue.Name, page.Query)
		}).OrderByDescendingT(func(project *models.Project) int64 {
			return project.CreateTs.UnixNano()
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

func QueryProjectById(id uuid.UUID) (result *models.Project, err error) {
	if err := dao.ReadTx(dao.MainDB, func(tx dao.Tx) error {
		result, err = tx.SelectProject(id)
		return err
	}); err != nil {
		return nil, err
	}
	return result, nil
}

func AddProject(name string) (result *models.Project, err error) {
	path := viper.GetString("python.path")
	currentUser, err := user.Current()
	if err != nil {
		return nil, err
	}
	id := uuid.New()
	projectDir := currentUser.HomeDir + "\\gobot\\"
	if !utils.PathExist(projectDir) {
		os.MkdirAll(projectDir, fs.ModeDir)
	}
	result = &models.Project{
		Id:          id,
		Name:        name,
		Description: "",
		Path:        projectDir + id.String(),
	}
	command := exec.Command(path, "-m", "venv", result.Path)
	// if err = command.Run(); err != nil {
	// 	return nil, err
	// }
	output, err := command.Output()
	if err != nil {
		fmt.Println(string(output))
	}
	if err := dao.WriteTx(dao.MainDB, func(tx dao.Tx) error {
		if err := tx.InsertProject(result); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return result, nil
}

func ModifyProject(id uuid.UUID, form models.Project) (result *models.Project, err error) {
	if err := dao.WriteTx(dao.MainDB, func(tx dao.Tx) error {
		project, err := tx.SelectProject(id)
		if err != nil {
			return err
		}
		if project == nil {
			return errors.New("project not found")
		}

		if err := tx.UpdateProject(project); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return result, nil
}

func RemoveSchedule(id uuid.UUID) (res interface{}, err error) {
	if err := dao.WriteTx(dao.MainDB, func(tx dao.Tx) error {
		if project, err := tx.SelectProject(id); err != nil {
			return err
		} else if project == nil {
			return errors.New("project not found")
		} else {
			if err = tx.DeleteProject(id); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return nil, nil
}
