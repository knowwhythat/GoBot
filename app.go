package main

import (
	"context"
	"fmt"
	"gobot/dao"
	"gobot/forms"
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
	result = map[string]interface{}{"total": total, "result": resultList}
	return result, err
}

func (a *App) NewWorkflow(name string) error {
	_, err := services.AddProject(name)
	return err
}
