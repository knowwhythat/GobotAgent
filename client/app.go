package main

import (
	"browser-agent/backend/models"
	"browser-agent/backend/service"
	"context"
	"github.com/spf13/viper"
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
}

func (a *App) OpenFileDialog(filters []runtime.FileFilter) string {
	filePath, _ := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{Title: "选择文件", Filters: filters})
	return filePath
}

func (a *App) GetSettingData() models.SettingData {
	return models.SettingData{
		APIBaseUrl:  viper.GetString(service.LLMApiBase),
		APIKey:      viper.GetString(service.LLMApiKey),
		ModelName:   viper.GetString(service.LLMModelName),
		BrowserPath: viper.GetString(service.ConfigBrowserPath),
	}
}

func (a *App) SetSettingData(settingData models.SettingData) {
	viper.Set(service.LLMApiBase, settingData.APIBaseUrl)
	viper.Set(service.LLMApiKey, settingData.APIKey)
	viper.Set(service.LLMModelName, settingData.ModelName)
	viper.Set(service.ConfigBrowserPath, settingData.BrowserPath)
	_ = viper.WriteConfig()
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) error {
	return service.ExecuteUserTask(a.ctx, name)
}
