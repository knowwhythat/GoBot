package main

import (
	"context"
	"fmt"
	"gobot/dao"
	"gobot/forms"
	"gobot/models"
	"gobot/services"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) shutdown(ctx context.Context) {
	dao.MainDB.Close()
	dao.LogDB.Close()
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

func (a *App) ListProject(query string, pageNum int) (result map[string]interface{}, err error) {
	total, resultList, err := services.QueryProjectPage(forms.QueryForm{
		Query: query,
		PageForm: forms.PageForm{
			PageNum:  1,
			PageSize: 12,
		},
	})
	result = map[string]interface{}{"total": total, "list": resultList}
	return result, err
}

func (a *App) SelectProject(id string) (project *models.Project, err error) {
	project, err = services.QueryProjectById(id)
	return project, err
}

func (a *App) UpdateProject(id string, name string) error {
	_, err := services.ModifyProject(id, models.Project{Name: name})
	return err
}

func (a *App) DeleteProject(id string) error {
	err := services.RemoveProject(id)
	return err
}

func (a *App) CreateProject(name string) error {
	_, err := services.AddProject(name)
	return err
}

func (a *App) GetMainFlow(id string) (map[string]string, error) {
	name, data, err := services.ReadMainFlow(id)
	result := map[string]string{"name": name, "data": data}
	return result, err
}

func (a *App) SaveMainFlow(id, data string) error {
	return services.SaveMainFLow(id, data)
}
