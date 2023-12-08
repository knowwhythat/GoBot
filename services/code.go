package services

import (
	"encoding/json"
	"gobot/log"
	"os"
	"strings"
)

type Flow struct {
	Sequence Activity `json:"sequence,omitempty"`
}

type Input struct {
	Key          string `json:"key,omitempty"`
	Label        string `json:"label,omitempty"`
	Type         string `json:"type,omitempty"`
	DefaultValue string `json:"default_value,omitempty"`
	Required     bool   `json:"required,omitempty"`
	Options      any    `json:"options,omitempty"`
	EditorType   string `json:"editor_type,omitempty"`
}
type Output struct {
	Key          string `json:"key,omitempty"`
	Label        string `json:"label,omitempty"`
	Type         string `json:"type,omitempty"`
	DefaultValue string `json:"default_value,omitempty"`
	Required     bool   `json:"required,omitempty"`
	EditorType   string `json:"editor_type,omitempty"`
}
type ParameterDefine struct {
	Error   bool     `json:"error,omitempty"`
	Inputs  []Input  `json:"inputs,omitempty"`
	Extra   []Input  `json:"extra,omitempty"`
	Outputs []Output `json:"outputs,omitempty"`
}
type Activity struct {
	ID              string            `json:"id,omitempty"`
	Key             string            `json:"key,omitempty"`
	Label           string            `json:"label,omitempty"`
	IconPath        string            `json:"icon_path,omitempty"`
	Description     string            `json:"description,omitempty"`
	Toggleable      bool              `json:"toggleable,omitempty"`
	Deleteable      bool              `json:"deleteable,omitempty"`
	Collapsed       bool              `json:"collapsed,omitempty"`
	Namespace       string            `json:"namespace,omitempty"`
	Method          string            `json:"method,omitempty"`
	ParameterDefine ParameterDefine   `json:"parameter_define,omitempty"`
	Component       string            `json:"component,omitempty"`
	Parameter       map[string]string `json:"parameter,omitempty"`
	IsPseudo        bool              `json:"isPseudo,omitempty"`
	Children        []Activity        `json:"children,omitempty"`
}

func (flow *Flow) GeneratePythonCode(filepath string) error {
	log.Logger.Info("开始生成代码")
	activities := flow.Sequence.Children
	namespace := make(map[string]string)
	code := "def main(args):\n"
	for _, activity := range activities {
		code = code + activity.GeneratePythonCode(namespace, 1)
	}
	importStr := ""
	for key := range namespace {
		importStr += "import " + key + "\n"
	}
	log.Logger.Info(importStr + code)
	err := os.WriteFile(filepath, []byte(importStr+code), 0666)
	return err
}

func (activity *Activity) GeneratePythonCode(namespace map[string]string, indent int) string {
	if activity.Namespace != "" {
		namespace[activity.Namespace] = ""
	}
	code := ""
	needPass := false
	if activity.IsPseudo {
		if activity.Method == "if" || activity.Method == "else if" {
			indent = indent - 1
			code += strings.Repeat(" ", indent*4)
			expression := activity.Parameter["expression"]
			if len(expression) >= 2 {
				code += activity.Method + " " + expression[2:] + " :\n"
			} else {
				code += activity.Method + " " + activity.ParameterDefine.Inputs[0].DefaultValue[2:] + ":\n"
			}
			needPass = true
		} else if activity.Method == "else" {
			indent = indent - 1
			code += strings.Repeat(" ", indent*4)
			code += activity.Key + ":\n"
			needPass = true
		} else if activity.Method == "while" {
			code += strings.Repeat(" ", indent*4)
			expression := activity.Parameter["expression"]
			if len(expression) >= 2 {
				code += activity.Method + " " + expression[2:] + " :\n"
			} else {
				code += activity.Method + " " + activity.ParameterDefine.Inputs[0].DefaultValue[2:] + ":\n"
			}
			needPass = true
		} else if activity.Key == "base.control.loop.foreach_loop_list" {
			code += strings.Repeat(" ", indent*4)
			code += activity.Method + " "
			if loopValue, ok := activity.Parameter["loop_value"]; !ok {
				code += activity.ParameterDefine.Outputs[0].DefaultValue[2:]
			} else {
				code += loopValue[2:]
			}
			code += " in "
			if arrayValue, ok := activity.Parameter["array"]; !ok {
				code += activity.ParameterDefine.Inputs[0].DefaultValue[2:]
			} else {
				code += arrayValue[2:]
			}
			code += ":\n"
			needPass = true
		} else if activity.Key == "base.control.loop.foreach_loop_map" {
			code += strings.Repeat(" ", indent*4)
			code += activity.Method + " "
			if loopKey, ok := activity.Parameter["loop_key"]; !ok {
				code += activity.ParameterDefine.Outputs[0].DefaultValue[2:]
			} else {
				code += loopKey[2:]
			}
			code += ","
			if loopValue, ok := activity.Parameter["loop_value"]; !ok {
				code += activity.ParameterDefine.Outputs[1].DefaultValue[2:]
			} else {
				code += loopValue[2:]
			}
			code += " in "
			if arrayValue, ok := activity.Parameter["map"]; !ok {
				code += activity.ParameterDefine.Inputs[0].DefaultValue[2:]
			} else {
				code += arrayValue[2:]
			}
			code += ":\n"
			needPass = true
		} else if activity.Key == "base.control.loop.for_i_loop" {
			code += strings.Repeat(" ", indent*4)
			code += activity.Method + " "
			if loopIndex, ok := activity.Parameter["loop_index"]; !ok {
				code += activity.ParameterDefine.Outputs[0].DefaultValue[2:]
			} else {
				code += loopIndex[2:]
			}
			code += " in range("
			if start, ok := activity.Parameter["start"]; !ok {
				code += activity.ParameterDefine.Inputs[0].DefaultValue[2:]
			} else {
				code += start[2:]
			}
			code += ", "
			if end, ok := activity.Parameter["end"]; !ok {
				code += activity.ParameterDefine.Inputs[1].DefaultValue[2:]
			} else {
				code += end[2:]
			}
			code += ", "
			if add, ok := activity.Parameter["add"]; !ok {
				code += activity.ParameterDefine.Inputs[2].DefaultValue[2:]
			} else {
				code += add[2:]
			}
			code += "):\n"
			needPass = true
		} else if activity.Method == "continue" || activity.Method == "break" {
			code += strings.Repeat(" ", indent*4)
			code += activity.Method + "\n"
		} else if activity.Method == "try" || activity.Method == "finally" {
			indent = indent - 1
			code += strings.Repeat(" ", indent*4)
			code += activity.Method + ":\n"
			needPass = true
		} else if activity.Method == "except" {
			indent = indent - 1
			code += strings.Repeat(" ", indent*4)
			code += activity.Method + " Exception as "
			expression := activity.Parameter["expression"]
			if len(expression) >= 2 {
				code += expression[2:] + " :\n"
			} else {
				code += activity.ParameterDefine.Inputs[0].DefaultValue[2:] + ":\n"
			}
			needPass = true
		}
	} else {
		code = code + strings.Repeat(" ", indent*4)
		outputValue := ""
		if len(activity.ParameterDefine.Outputs) > 0 {
			for _, o := range activity.ParameterDefine.Outputs {
				if outputName, ok := activity.Parameter[o.Key]; !ok {
					outputValue += o.DefaultValue + ","
				} else if outputName != "" {
					outputValue += outputName + ","
				}
			}
		}
		if len(outputValue) > 0 && outputValue[len(outputValue)-1] == ',' {
			outputValue = strings.TrimRight(outputValue, ",")
		}

		if len(outputValue) > 0 {
			code = code + outputValue + " = "
		}
		code = code + activity.Namespace + "." + activity.Method + "("
		if len(activity.ParameterDefine.Inputs) > 0 {
			for _, i := range activity.ParameterDefine.Inputs {
				if inputValue, ok := activity.Parameter[i.Key]; !ok {
					if i.DefaultValue[0] == '0' {
						code += i.Key + "=\"" + jsonEscape(i.DefaultValue[2:]) + "\","
					} else {
						code += i.Key + "=" + i.DefaultValue[2:] + ","
					}
				} else if len(inputValue) > 1 && inputValue[1] == ':' {
					if inputValue[0] == '0' {
						code += i.Key + "=\"" + jsonEscape(inputValue[2:]) + "\","
					} else {
						code += i.Key + "=" + inputValue[2:] + ","
					}
				}
			}
		}
		if len(activity.ParameterDefine.Extra) > 0 {
			for _, i := range activity.ParameterDefine.Extra {
				if inputValue, ok := activity.Parameter[i.Key]; !ok {
					if i.DefaultValue[0] == '0' {
						code += i.Key + "=\"" + jsonEscape(i.DefaultValue[2:]) + "\","
					} else {
						code += i.Key + "=" + i.DefaultValue[2:] + ","
					}
				} else if len(inputValue) > 1 && inputValue[1] == ':' {
					if inputValue[0] == '0' {
						code += i.Key + "=\"" + jsonEscape(inputValue[2:]) + "\","
					} else {
						code += i.Key + "=" + inputValue[2:] + ","
					}
				}
			}
		}
		code += "code_map_id=\"" + activity.ID + "\"" + ")\n"
	}
	if len(activity.Children) > 0 {
		for _, child := range activity.Children {
			code = code + child.GeneratePythonCode(namespace, indent+1)
		}
	} else if needPass {
		code += strings.Repeat(" ", (indent+1)*4) + "pass\n"
	}
	return code
}

func jsonEscape(s string) string {
	b, _ := json.Marshal(s)
	return string(b[1 : len(b)-1])
}
