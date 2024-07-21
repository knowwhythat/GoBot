export const BlockType = {
  Start: {
    blockId: "Start",
    name: "开始",
    description: "流程执行的起点,每个工作流中有且只能有一个开始节点",
    icon: "riPlayCircleLine",
    component: "BlockBasic",
    disableEdit: true,
    disableDelete: true,
    disableSetting: true,
    color: "bg-green-200 dark:bg-green-300 fill-green-200 dark:fill-green-300",
    data: {
      nodeType: "Start",
      label: "开始",
    },
  },
  SubFlow: {
    blockId: "SubFlow",
    name: "子流程",
    description:
      "一个相对独立的功能块,可以编写具体的业务逻辑,一个工作流由多个子流程组合而成",
    icon: "riFlowChart",
    component: "BlockBasic",
    color:
      "bg-orange-200 dark:bg-orange-300 fill-orange-200 dark:fill-orange-300",
    data: {
      nodeType: "SubFlow",
      label: "子流程",
      errorHandler: {
        errorEnable: false,
        retry: false,
        retryCount: 2,
        retryTimeout: 1000,
        toDo: "error",
      },
      params: {},
    },
  },
  Condition: {
    blockId: "Condition",
    name: "条件",
    description: "控制流程执行逻辑,根据不同的条件,控制流程执行不同的逻辑",
    icon: "riAB",
    component: "BlockConditions",
    inputs: 1,
    outputs: 0,
    color: "bg-lime-200 dark:bg-lime-300 fill-lime-200 dark:fill-lime-300",
    allowedInputs: true,
    maxConnection: 1,
    disableEdit: true,
    data: {
      nodeType: "Condition",
      label: "",
      description: "",
      conditions: [],
    },
  },
  Comment: {
    blockId: "Comment",
    name: "注释",
    description: "注释模块,在流程中编写注释,解释流程运行逻辑",
    icon: "riFileEditLine",
    component: "BlockNote",
    inputs: 1,
    outputs: 1,
    color: "bg-cyan-200 dark:bg-cyan-300 fill-cyan-200 dark:fill-cyan-300",
    maxConnection: 1,
    disableEdit: true,
    data: {
      nodeType: "Comment",
      label: "注释",
      note: "",
      width: 280,
      height: 168,
      color: "white",
      fontSize: "regular",
    },
  },
};

export const conditionBuilder = {
  valueTypes: [
    {
      id: "value",
      category: "value",
      name: "Value",
      compareable: true,
      data: { value: "" },
    },
    {
      id: "code",
      category: "value",
      name: "Code",
      compareable: false,
      data: { code: "\nreturn true;", context: "background" },
    },
    {
      id: "data#exists",
      category: "value",
      name: "Data exists",
      compareable: false,
      valueKey: "dataPath",
      data: { dataPath: "" },
    },
    {
      id: "element#text",
      category: "element",
      name: "Element text",
      compareable: true,
      data: { selector: "" },
    },
    {
      id: "element#exists",
      category: "element",
      name: "Element exists",
      compareable: false,
      data: { selector: "" },
    },
    {
      id: "element#notExists",
      category: "element",
      name: "Element not exists",
      compareable: false,
      data: { selector: "" },
    },
    {
      id: "element#visible",
      category: "element",
      name: "Element visible",
      compareable: false,
      data: { selector: "" },
    },
    {
      id: "element#visibleScreen",
      category: "element",
      name: "Element visible in screen",
      compareable: false,
      data: { selector: "" },
    },
    {
      id: "element#invisible",
      category: "element",
      name: "Element hidden in screen",
      compareable: false,
      data: { selector: "" },
    },
    {
      id: "element#attribute",
      category: "element",
      name: "Element attribute value",
      compareable: true,
      data: { selector: "", attrName: "" },
    },
  ],
  compareTypes: [
    { id: "eq", name: "Equals", needValue: true, category: "basic" },
    {
      id: "eqi",
      name: "Equals (case insensitive)",
      needValue: true,
      category: "basic",
    },
    { id: "nq", name: "Not equals", needValue: true, category: "basic" },
    { id: "gt", name: "Greater than", needValue: true, category: "number" },
    {
      id: "gte",
      name: "Greater than or equal",
      needValue: true,
      category: "number",
    },
    { id: "lt", name: "Less than", needValue: true, category: "number" },
    {
      id: "lte",
      name: "Less than or equal",
      needValue: true,
      category: "number",
    },
    { id: "cnt", name: "Contains", needValue: true, category: "text" },
    {
      id: "cni",
      name: "Contains (case insensitive)",
      needValue: true,
      category: "text",
    },
    { id: "nct", name: "Not contains", needValue: true, category: "text" },
    {
      id: "nci",
      name: "Not contains (case insensitive)",
      needValue: true,
      category: "text",
    },
    { id: "stw", name: "Starts with", needValue: true, category: "text" },
    { id: "enw", name: "Ends with", needValue: true, category: "text" },
    { id: "rgx", name: "Match with RegEx", needValue: true, category: "text" },
    { id: "itr", name: "Is truthy", needValue: false, category: "boolean" },
    { id: "ifl", name: "Is falsy", needValue: false, category: "boolean" },
  ],
  inputTypes: {
    selector: {
      placeholder: ".class",
      label: "CSS selector or XPath",
    },
    value: {
      label: "Value",
      placeholder: "abc123",
    },
    attrName: {
      label: "Attribute name",
      placeholder: "name",
    },
    dataPath: {
      label: "variables@variableName",
      placeholder: "",
    },
  },
};

export const typeBuilder = {
  str: {
    icon: "pi pi-dollar",
    property: [{ label: "trim", type: "method" }],
  },
  number: {
    icon: "pi pi-plus-circle",
    property: [{ label: "trim", type: "method" }],
  },
};

export const paramTypes = [
  { name: "普通文本", code: "str" },
  { name: "文件路径", code: "filePath" },
  { name: "文件目录", code: "dirPath" },
  { name: "密码", code: "password" },
  { name: "数字", code: "number" },
  { name: "布尔", code: "bool" },
  { name: "单选", code: "single" },
  { name: "多选", code: "multiple" },
  { name: "任意", code: "any" },
];

export const variableType = [
  { name: "Python表达式", code: "1:" },
  { name: "普通文本", code: "0:" },
];
export const matchTypes = [
  { name: "模板匹配", value: "template" },
  { name: "特征点匹配", value: "sift" },
];
