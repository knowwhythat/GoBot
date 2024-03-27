package backend

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"gobot/backend/dao"
	"gobot/backend/forms"
	"gobot/backend/models"
	"gobot/backend/plugin"
	"gobot/backend/services"
	"gobot/backend/services/sys_tray"
	"gobot/backend/utils"

	"github.com/spf13/viper"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

const aesKey = "I98NTHNPezFnbe8iCaSc1xMPAv8ZtTil"

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

//go:embed services/appicon.ico
var icon []byte

// Startup is called when the app starts. The context is saved,
// so we can call the runtime methods
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
	sys_tray.CreateTray(ctx, icon)
}

func (a *App) Shutdown(ctx context.Context) {
	_ = dao.MainDB.Close()
	_ = dao.LogDB.Close()
	services.StopWindowsInspectCommand()
}

func (a *App) OpenDialog(option map[string]string) string {
	if option["type"] == "file" {
		filter := runtime.FileFilter{DisplayName: option["display"], Pattern: option["pattern"]}
		filePath, _ := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{Title: "选择文件", Filters: []runtime.FileFilter{filter}})
		return filePath
	} else if option["type"] == "directory" {
		filePath, _ := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{Title: "选择目录"})
		return filePath
	}

	return ""
}

func (a *App) GetLoginData() forms.LoginForm {
	username := viper.GetString("login.username")
	pwd := viper.GetString("login.pwd")
	if pwd != "" {
		decryptPwd, err := utils.DecryptAES2Base64([]byte(aesKey), pwd)
		if err == nil {
			pwd = string(decryptPwd)
		} else {
			pwd = ""
		}
	}
	rememberMe := viper.GetBool("login.rememberMe")
	autoLogin := viper.GetBool("login.autoLogin")
	return forms.LoginForm{
		Username:   username,
		Pwd:        pwd,
		RememberMe: rememberMe,
		AutoLogin:  autoLogin,
	}
}

func (a *App) Login(form forms.LoginForm) error {
	sysInfo := services.GetSysInfo()
	fmt.Printf("sysInfo: %#v\n", sysInfo)
	if form.RememberMe {
		pwd, _ := utils.EncryptAES2Base64([]byte(aesKey), []byte(form.Pwd))
		viper.Set("login.username", form.Username)
		viper.Set("login.pwd", pwd)
		viper.Set("login.rememberMe", form.RememberMe)
		viper.Set("login.autoLogin", form.AutoLogin)
		_ = viper.WriteConfig()
	}
	services.InitSchedule(a.ctx)
	return nil
}

func (a *App) ListProject() (resultList []*models.Project, err error) {
	resultList, err = services.QueryProjectPage()
	return
}

func (a *App) SelectProject(id string) (project *models.Project, err error) {
	project, err = services.QueryProjectById(id)
	return project, err
}

func (a *App) DeleteProject(id string) error {
	err := services.RemoveProject(id)
	return err
}

func (a *App) AddOrUpdateProject(id, name, desc string, isFlow bool) error {
	_, err := services.AddOrUpdateProject(id, name, desc, isFlow)
	return err
}

func (a *App) RunMainFlow(id string) error {
	project, err := services.QueryProjectById(id)
	if err != nil {
		return err
	}
	if project.IsFlow {
		return errors.New("暂不支持流程图类型的项目运行")
	} else {
		return services.RunSubFlow(a.ctx, id, "main")
	}
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

func (a *App) ParseAllPlugin() ([]plugin.Activity, error) {
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

func (a *App) GetElementImage(id, elementId string) (string, error) {
	return services.GetElementImage(id, elementId)
}

func (a *App) StartPickWindowsElement() (string, error) {
	return services.StartPickWindowsElement(a.ctx)
}

func (a *App) StartCheckWindowsElement(paths string) (string, error) {
	return services.StartCheckWindowsElement(a.ctx, paths)
}

func (a *App) SaveWindowsElement(id, elementId, image, data string) error {
	return services.SaveWindowsElement(id, elementId, image, data)
}

func (a *App) GetWindowsElement(id, elementId string) (string, error) {
	return services.GetWindowsElement(id, elementId)
}

func (a *App) GetWindowsElementList(parentId string) (string, error) {
	return services.GetWindowsElementList(parentId)
}

func (a *App) HighlightCurrentElement(controlId string) error {
	return services.HighlightCurrentElement(controlId)
}

func (a *App) GetSelectedWindowsElement(controlId string) (string, error) {
	return services.GetSelectedWindowsElement(a.ctx, controlId)
}

func (a *App) GetProjectWindowsElements(projectId string) (string, error) {
	return services.GetProjectWindowsElements(projectId)
}

func (a *App) RemoveWindowsElement(projectId, controlId string) error {
	return services.RemoveWindowsElement(projectId, controlId)
}

func (a *App) QuerySchedulePage() (resultList []*models.Schedule, err error) {
	return services.QuerySchedulePage()
}

func (a *App) AddOrUpdateSchedule(schedule models.Schedule) (err error) {
	return services.AddOrUpdateSchedule(schedule)
}

func (a *App) ToggleScheduleById(id string, state bool) (err error) {
	return services.ToggleScheduleById(id, state)
}

func (a *App) GetNextTriggerTime(cron string) (result string, err error) {
	return services.GetNextTriggerTime(cron)
}

func (a *App) RemoveSchedule(id string) (err error) {
	return services.RemoveSchedule(id)
}
