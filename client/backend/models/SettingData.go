package models

type SettingData struct {
	BrowserPath string `json:"browserPath"`
	APIBaseUrl  string `json:"apiBaseUrl"`
	APIKey      string `json:"apiKey"`
	ModelName   string `json:"modelName"`
}
