package plugin

import (
	"embed"
	"encoding/json"
	"io"
)

//go:embed config
var f embed.FS

var fileList = []string{
	"base_control.json",
	"tools.json",
	"data_process.json",
	"browser.json",
	"desktop_application.json",
	"mouse_keyboard.json",
	"image.json",
	"ai.json",
}

type PluginConfig struct {
	Id          string     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Version     string     `json:"version"`
	Activities  []Activity `json:"activities"`
}

type Activity struct {
	Key             string          `json:"key"`
	Label           string          `json:"label"`
	IconPath        string          `json:"icon_path"`
	Description     string          `json:"description"`
	Namespace       string          `json:"namespace,omitempty"`
	Method          string          `json:"method,omitempty"`
	IsPseudo        bool            `json:"isPseudo,omitempty"`
	ParameterDefine ParameterDefine `json:"parameter_define,omitempty"`
	Component       string          `json:"component,omitempty"`
	Children        []Activity      `json:"children,omitempty"`
}

type ParameterDefine struct {
	Error   bool     `json:"error,omitempty"`
	Inputs  []Input  `json:"inputs,omitempty"`
	Extra   []Input  `json:"extra,omitempty"`
	Outputs []Output `json:"outputs,omitempty"`
}

type Input struct {
	Key          string   `json:"key"`
	Label        string   `json:"label"`
	Placeholder  string   `json:"placeholder"`
	Category     string   `json:"category,omitempty"`
	Type         string   `json:"type"`
	DefaultValue string   `json:"default_value"`
	Required     bool     `json:"required"`
	Options      any      `json:"options"`
	EditorType   string   `json:"editor_type"`
	ShowIf       []string `json:"show_if,omitempty"`
}
type Output struct {
	Key          string   `json:"key"`
	Label        string   `json:"label"`
	Placeholder  string   `json:"placeholder"`
	Type         string   `json:"type"`
	DefaultValue string   `json:"default_value"`
	Required     bool     `json:"required"`
	EditorType   string   `json:"editor_type"`
	ShowIf       []string `json:"show_if,omitempty"`
}

func ParseAllPlugin() ([]Activity, error) {
	var activities []Activity
	for _, name := range fileList {
		pluginFile, _ := f.Open("config/" + name)
		content, err := io.ReadAll(pluginFile)
		if err != nil {
			return nil, err
		}
		pluginConfig := PluginConfig{}
		err = json.Unmarshal(content, &pluginConfig)
		if err != nil {
			return nil, err
		}
		activities = append(activities, pluginConfig.Activities...)
	}
	return activities, nil
}
