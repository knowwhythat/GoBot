package main

import (
	"context"
	"gobot/dao"
	"gobot/forms"
	"gobot/models"
	"gobot/plugin"
	"gobot/services"
	"gobot/services/sys_tray"

	"github.com/wailsapp/wails/v2/pkg/runtime"
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
	sys_tray.CreateTray(ctx, icon)
}

func (a *App) shutdown(ctx context.Context) {
	dao.MainDB.Close()
	dao.LogDB.Close()
}

func (a *App) beforeClose(ctx context.Context) bool {
	dialog, err := runtime.MessageDialog(ctx, runtime.MessageDialogOptions{
		Type:    runtime.QuestionDialog,
		Title:   "退出?",
		Message: "退出后所有的任务将停止,是否确定退出?",
	})

	if err != nil {
		return false
	}
	return dialog != "Yes"
}

func (a *App) OpenDialog(option map[string]string) string {
	if option["type"] == "file" {
		fileter := runtime.FileFilter{DisplayName: option["display"], Pattern: option["pattern"]}
		filePath, _ := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{Title: "选择文件", Filters: []runtime.FileFilter{fileter}})
		return filePath
	} else if option["type"] == "directory" {
		filePath, _ := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{Title: "选择目录"})
		return filePath
	}

	return ""
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
	return services.SaveMainFlow(id, data)
}

func (a *App) GetSubFlow(id, subId string) (string, error) {
	data, err := services.ReadSubFlow(id, subId)
	return data, err
}

func (a *App) SaveSubFlow(id, subId, data string) error {
	return services.SaveSubFlow(id, subId, data)
}

func (a *App) DeleteSubFlow(id, subId string) error {
	return services.DeleteSubFlow(id, subId)
}

func (a *App) RunSubFlow(id, subId string) error {
	return services.RunSubFlow(a.ctx, id, subId)
}

func (a *App) DebugSubFlow(id, subId string) error {
	return services.DeubugSubFlow(a.ctx, id, subId)
}

func (a *App) TerminateSubFlow(id, subId string) error {
	return services.TerminateSubFlow(id, subId)
}

func (a *App) DealDebugSignal(command string) error {
	return services.DealDebugSignal(command)
}

func (a *App) ParseAllPlugin() ([]plugin.Activitiy, error) {
	return plugin.ParseAllPlugin()
}

func (a *App) StartPick() (string, error) {
	return services.StartPick(a.ctx)
}

func (a *App) StartCheck(frame, selector string) (string, error) {
	return services.StartCheck(a.ctx, frame, selector)
}

func (a *App) RunActivity(id, req string) error {
	return services.RunActivity(a.ctx, id, req)
}

func (a *App) RestartReplCommand(id string) error {
	return services.RestartReplCommand(id)
}
