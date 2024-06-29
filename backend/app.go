package backend

import (
	"context"
	_ "embed"
	"errors"
	"gobot/backend/constants"
	"gobot/backend/dao"
	"gobot/backend/forms"
	"gobot/backend/log"
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
	username := viper.GetString(constants.ConfigLoginUsername)
	pwd := viper.GetString(constants.ConfigLoginPwd)
	if pwd != "" {
		decryptPwd, err := utils.DecryptAES2Base64([]byte(aesKey), pwd)
		if err == nil {
			pwd = string(decryptPwd)
		} else {
			pwd = ""
		}
	}
	rememberMe := viper.GetBool(constants.ConfigLoginRemember)
	autoLogin := viper.GetBool(constants.ConfigLoginAuto)
	return forms.LoginForm{
		Username:   username,
		Pwd:        pwd,
		RememberMe: rememberMe,
		AutoLogin:  autoLogin,
	}
}

func (a *App) GetBasicConfigData() map[string]string {
	return viper.GetStringMapString("basic")
}

func (a *App) SaveBasicConfigData(data map[string]string) error {
	for k, v := range data {
		viper.Set("basic."+k, v)
	}
	return viper.WriteConfig()
}

func (a *App) Login(form forms.LoginForm) error {
	sysInfo := services.GetSysInfo()
	log.Logger.Infof("sysInfo: %#v", sysInfo)
	if form.RememberMe {
		pwd, _ := utils.EncryptAES2Base64([]byte(aesKey), []byte(form.Pwd))
		viper.Set(constants.ConfigLoginUsername, form.Username)
		viper.Set(constants.ConfigLoginPwd, pwd)
		viper.Set(constants.ConfigLoginRemember, form.RememberMe)
		viper.Set(constants.ConfigLoginAuto, form.AutoLogin)
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

func (a *App) AddOrUpdateProject(project models.Project) error {
	return services.AddOrUpdateProject(project)
}

func (a *App) RunMainFlow(id string) error {
	project, err := services.QueryProjectById(id)
	if err != nil {
		return err
	}
	if project.IsFlow {
		return errors.New("暂不支持流程图类型的项目运行")
	} else {
		return services.RunSequence(a.ctx, id, "main", "手动触发")
	}
}

func (a *App) TerminateMainFlow(id string) error {
	return services.TerminateFlow(id)
}

func (a *App) ReadProjectConfig(id string) (models.ProjectConfig, error) {
	result, err := services.ReadProjectConfig(id)
	return result, err
}

func (a *App) SaveProjectConfig(id, data string) error {
	err := services.SaveProjectConfig(id, data)
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

func (a *App) GetRunningFlows() (resultList []forms.RunningInstance, err error) {
	return services.GetRunningFlows()
}

func (a *App) DebugSubFlow(id, subId string) error {
	return services.DebugSubFlow(a.ctx, id, subId)
}

func (a *App) TerminateSubFlow() error {
	return services.TerminateSubFlow()
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
	return services.SaveWindowsElement(a.ctx, id, elementId, image, data)
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
	return services.RemoveWindowsElement(a.ctx, projectId, controlId)
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

func (a *App) QueryAllExecution() (resultList []*forms.ExecutionForm, err error) {
	return services.QueryAllExecution()
}

func (a *App) RemoveExecution(executions []forms.ExecutionForm) (err error) {
	for _, execution := range executions {
		err := services.RemoveExecution(execution.Id)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *App) SelectExecutionLog(id string) (logs string, err error) {
	return services.SelectExecutionLog(id)
}

func (a *App) PipInstallPackage(id string, pip forms.PipPackage) (err error) {
	return services.InstallPackage(a.ctx, id, pip)
}

func (a *App) PipListInstallPackage(id string) (output string, err error) {
	return services.ListInstallPackage(id)
}

func (a *App) PipUnInstallPackage(id, name string) (err error) {
	return services.UnInstallPackage(id, name)
}

func (a *App) SaveProjectDependency(id string, packages []string) (err error) {
	return services.SaveProjectDependency(id, packages)
}
