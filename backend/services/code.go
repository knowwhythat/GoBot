package services

import (
	"gobot/backend/log"
	"gobot/backend/models"
	"os"
	"strconv"
	"strings"

	"github.com/duke-git/lancet/v2/slice"
)

type Flow struct {
	Params   []Param  `json:"params,omitempty"`
	Sequence Activity `json:"sequence,omitempty"`
}

type Param struct {
	Name      string   `json:"name,omitempty"`
	Direction string   `json:"direction,omitempty"`
	Type      string   `json:"type,omitempty"`
	Options   []Option `json:"options,omitempty"`
	Value     string   `json:"value,omitempty"`
}

type Option struct {
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
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
	Breakpoint      bool              `json:"breakpoint,omitempty"`
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

func GenerateGlobalVariable(variables []*models.Variable, path string) error {
	code := "variables = {} \n"
	for _, v := range variables {
		if v.Kind == "expression" {
			code += "variables['" + v.Name + "'] = " + v.Value + "\n"
		} else {
			code += "variables['" + v.Name + "'] = " + strconv.Quote(v.Value) + "\n"
		}
	}
	err := os.WriteFile(path, []byte(code), 0666)
	return err
}

func (flow *Flow) GeneratePythonCode(filepath string) error {
	log.Logger.Info("开始生成代码")
	activities := flow.Sequence.Children
	namespace := make(map[string]string)
	code := "def main(args):\n"
	inParams := slice.Filter(flow.Params, func(index int, item Param) bool { return item.Direction == "In" })
	if len(inParams) > 0 {
		code += generateParams(inParams)
	}
	if len(activities) > 0 {
		for _, activity := range activities {
			code = code + activity.GeneratePythonCode(namespace, 1)
		}
	} else {
		code += strings.Repeat(" ", 1*4) + "pass"
	}

	importStr := "from package import variables as glv\n"
	for key := range namespace {
		importStr += "import " + key + "\n"
	}
	log.Logger.Info(importStr + code)
	err := os.WriteFile(filepath, []byte(importStr+code), 0666)
	return err
}

func generateParams(params []Param) string {
	ifExpression := strings.Repeat(" ", 1*4) + "if args is None:\n"
	elseExpression := strings.Repeat(" ", 1*4) + "else:\n"
	for _, param := range params {
		defaultValue := "None"
		if param.Type == "str" || param.Type == "filePath" || param.Type == "dirPath" || param.Type == "password" || param.Type == "single" || param.Type == "multiple" {
			defaultValue = strconv.Quote(param.Value)
		} else {
			defaultValue = param.Value
		}
		if len(defaultValue) == 0 {
			defaultValue = "\"\""
		}
		ifExpression += strings.Repeat(" ", 2*4) + param.Name + " = " + defaultValue + "\n"
		elseExpression += strings.Repeat(" ", 2*4) + param.Name + " = args.get(\"" + param.Name + "\", " + defaultValue + ")" + "\n"
	}
	return ifExpression + elseExpression
}

func (activity *Activity) GeneratePythonCode(namespace map[string]string, indent int) string {
	if activity.Namespace != "" {
		namespace[activity.Namespace] = ""
	}
	code := ""
	needPass := false
	if activity.IsPseudo {
		if activity.Method == "if" || activity.Method == "elif" {
			indent, code, needPass = generateIfExpression(indent, *activity)
		} else if activity.Key == "base.control.group" {
			indent, code, needPass = generateCommentExpression(indent, *activity)
		} else if activity.Method == "else" {
			indent, code, needPass = generateElseExpression(indent, *activity)
		} else if activity.Method == "while" {
			indent, code, needPass = generateWhileExpression(indent, *activity)
		} else if activity.Key == "base.control.loop.foreach_loop_list" {
			indent, code, needPass = generateForEachListExpression(indent, *activity)
		} else if activity.Key == "base.control.loop.foreach_loop_map" {
			indent, code, needPass = generateForEachMapExpression(indent, *activity)
		} else if activity.Key == "base.control.loop.for_i_loop" {
			indent, code, needPass = generateForLoopIndexExpression(indent, *activity)
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
				code += activity.ParameterDefine.Outputs[0].DefaultValue[2:] + ":\n"
			}
			needPass = true
		} else if activity.Method == "raise" {
			code += strings.Repeat(" ", indent*4)
			code += activity.Method + " Exception("
			expression := ""
			if _, ok := activity.Parameter["expression"]; ok {
				expression = activity.Parameter["expression"]
			} else {
				expression = activity.ParameterDefine.Inputs[0].DefaultValue
			}
			if expression[0] == '0' {
				code += "" + strconv.Quote(expression[2:]) + ")\n"
			} else {
				if len(expression[2:]) > 0 {
					code += expression[2:] + ")\n"
				} else {
					code += "None)\n"
				}
			}
			needPass = false
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
				if ok, param := generateParameter(activity.Parameter, i); ok {
					code += param
				}
			}
		}
		if len(activity.ParameterDefine.Extra) > 0 {
			for _, i := range activity.ParameterDefine.Extra {
				if ok, param := generateParameter(activity.Parameter, i); ok {
					code += param
				}
			}
		}
		if activity.ParameterDefine.Error {
			if retry, ok := activity.Parameter["retry"]; ok {
				code += "retry=" + retry + ","
			}
			if exception, ok := activity.Parameter["exception"]; ok {
				code += "exception=\"" + exception + "\","
			}
			if tryCount, ok := activity.Parameter["retry_count"]; ok {
				code += "retry_count=" + tryCount + ","
			}
			if tryInterval, ok := activity.Parameter["retry_interval"]; ok {
				code += "retry_interval=" + tryInterval + ","
			}
		}
		code += "local_data=locals(),"
		code += "code_map_id=\"" + activity.ID + "\"" + ")\n"
	}
	if len(activity.Children) > 0 {
		for _, child := range activity.Children {
			code += child.GeneratePythonCode(namespace, indent+1)
		}
	} else if needPass {
		code += strings.Repeat(" ", (indent+1)*4) + "pass\n"
	}
	return code
}

func generateIfExpression(indent int, activity Activity) (int, string, bool) {
	indent = indent - 1
	code := strings.Repeat(" ", indent*4)
	expression := activity.Parameter["expression"]
	if len(expression) >= 2 {
		code += activity.Method + " " + expression[2:] + " :\n"
	} else {
		code += activity.Method + " " + activity.ParameterDefine.Inputs[0].DefaultValue[2:] + ":\n"
	}
	return indent, code, true
}

func generateCommentExpression(indent int, activity Activity) (int, string, bool) {
	code := strings.Repeat(" ", indent*4)
	indent = indent - 1
	code += "#" + activity.Label + "\n"
	return indent, code, false
}

func generateElseExpression(indent int, activity Activity) (int, string, bool) {
	indent = indent - 1
	code := strings.Repeat(" ", indent*4)
	code += activity.Key + ":\n"
	return indent, code, true
}

func generateWhileExpression(indent int, activity Activity) (int, string, bool) {
	code := strings.Repeat(" ", indent*4)
	expression := activity.Parameter["expression"]
	if len(expression) >= 2 {
		code += activity.Method + " " + expression[2:] + " :\n"
	} else {
		code += activity.Method + " " + activity.ParameterDefine.Inputs[0].DefaultValue[2:] + ":\n"
	}
	return indent, code, true
}

func generateForEachListExpression(indent int, activity Activity) (int, string, bool) {
	code := strings.Repeat(" ", indent*4)
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
	return indent, code, true
}

func generateForEachMapExpression(indent int, activity Activity) (int, string, bool) {
	code := strings.Repeat(" ", indent*4)
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
	return indent, code, true
}

func generateForLoopIndexExpression(indent int, activity Activity) (int, string, bool) {
	code := strings.Repeat(" ", indent*4)
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
	return indent, code, true
}

// generateParameter
// 0:字符串
// 1:表达式
// 2:变量
// 3:混合模式
func generateParameter(parameter map[string]string, input Input) (bool, string) {
	code := ""
	if inputValue, ok := parameter[input.Key]; !ok {
		if input.DefaultValue[0] == '0' {
			code += input.Key + "=" + strconv.Quote(input.DefaultValue[2:]) + ","
			return true, code
		} else {
			if len(input.DefaultValue[2:]) > 0 {
				code += input.Key + "=" + input.DefaultValue[2:] + ","
				return true, code
			} else {
				code += input.Key + "=None,"
				return true, code
			}
		}
	} else if len(inputValue) > 1 && inputValue[1] == ':' {
		if inputValue[0] == '0' {
			code += input.Key + "=" + strconv.Quote(inputValue[2:]) + ","
			return true, code
		} else {
			if len(inputValue[2:]) > 0 {
				code += input.Key + "=" + inputValue[2:] + ","
				return true, code
			} else {
				code += input.Key + "=None,"
				return true, code
			}
		}
	}
	return false, code
}
