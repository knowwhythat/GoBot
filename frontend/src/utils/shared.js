export const tasks = {
  '开始': {
    name: '开始',
    description: 'Block where workflow will start executing',
    icon: 'riPlayCircleLine',
    component: 'BlockBasic',
    editComponent: 'EditTrigger',
    category: 'general',
    inputs: 0,
    outputs: 1,
    allowedInputs: true,
    maxConnection: 1,
    refDataKeys: ['url'],
    data: {
      disableBlock: false,
      description: '',
      type: 'manual',
      interval: 60,
      delay: 5,
      date: '',
      time: '00:00',
      url: '',
      shortcut: '',
      activeInInput: false,
      isUrlRegex: false,
      days: [],
      contextMenuName: '',
      contextTypes: [],
      parameters: [],
      preferParamsInTab: false,
      observeElement: {
        selector: '',
        baseSelector: '',
        matchPattern: '',
        targetOptions: {
          subtree: false,
          childList: true,
          attributes: false,
          attributeFilter: [],
          characterData: false,
        },
        baseElOptions: {
          subtree: false,
          childList: true,
          attributes: false,
          attributeFilter: [],
          characterData: false,
        },
      },
    },
  },
  '子流程': {
    name: '子流程',
    description: '',
    icon: 'riFlowChart',
    component: 'BlockBasic',
    category: 'general',
    editComponent: 'EditExecuteWorkflow',
    inputs: 1,
    outputs: 1,
    allowedInputs: true,
    maxConnection: 1,
    refDataKeys: ['globalData'],
    data: {
      disableBlock: false,
      executeId: '',
      workflowId: '',
      globalData: '',
      description: '',
      insertAllVars: false,
    },
  },
  '条件': {
    name: '条件',
    description: 'Conditional block',
    icon: 'riAB',
    component: 'BlockConditions',
    editComponent: 'EditConditions',
    category: 'conditions',
    inputs: 1,
    outputs: 0,
    allowedInputs: true,
    maxConnection: 1,
    data: {
      description: '',
      disableBlock: false,
      conditions: [],
      retryConditions: false,
      retryCount: 10,
      retryTimeout: 1000,
    },
  },
  '注释': {
    name: '注释',
    description: '',
    icon: 'riFileEditLine',
    component: 'BlockNote',
    category: 'general',
    disableEdit: true,
    inputs: 1,
    outputs: 1,
    allowedInputs: true,
    maxConnection: 1,
    data: {
      disableBlock: false,
      note: '',
      drawing: false,
      width: 280,
      height: 168,
      color: 'white',
      fontSize: 'regular',
    },
  },
};

export const categories = {
  interaction: {
    name: 'Web interaction',
    border: 'border-green-200 dark:border-green-300',
    color: 'bg-green-200 dark:bg-green-300 fill-green-200 dark:fill-green-300',
  },
  browser: {
    name: 'Browser',
    border: 'border-orange-200 dark:border-orange-300',
    color:
      'bg-orange-200 dark:bg-orange-300 fill-orange-200 dark:fill-orange-300',
  },
  general: {
    name: 'General',
    border: 'border-yellow-200 dark:border-yellow-300',
    color:
      'bg-yellow-200 dark:bg-yellow-300 fill-yellow-200 dark:fill-yellow-300',
  },
  onlineServices: {
    name: 'Online services',
    border: 'border-red-200 dark:border-red-300',
    color: 'bg-red-200 dark:bg-red-300 fill-red-200 dark:fill-red-300',
  },
  data: {
    name: 'Data',
    border: 'border-lime-200 dark:border-lime-300',
    color: 'bg-lime-200 dark:bg-lime-300 fill-lime-200 dark:fill-lime-300',
  },
  conditions: {
    name: 'Control flow',
    border: 'border-blue-200 dark:border-blue-300',
    color: 'bg-blue-200 dark:bg-blue-300 fill-blue-200 dark:fill-blue-300',
  },
  package: {
    name: 'Packages',
    border: 'border-cyan-200 dark:border-cyan-300',
    color: 'bg-cyan-200 dark:bg-cyan-300 fill-cyan-200 dark:fill-cyan-300',
  },
};

export const tagColors = {
  stage: 'bg-yellow-200 dark:bg-yellow-300',
  production: 'bg-green-200 dark:bg-green-300',
};

export const eventList = [
  { id: 'click', name: 'Click', type: 'mouse-event' },
  { id: 'dblclick', name: 'Double Click', type: 'mouse-event' },
  { id: 'mouseup', name: 'Mouseup', type: 'mouse-event' },
  { id: 'mousedown', name: 'Mousedown', type: 'mouse-event' },
  { id: 'mouseenter', name: 'Mouseenter', type: 'mouse-event' },
  { id: 'mouseleave', name: 'Mouseleave', type: 'mouse-event' },
  { id: 'mouseover', name: 'Mouseover', type: 'mouse-event' },
  { id: 'mouseout', name: 'Mouseout', type: 'mouse-event' },
  { id: 'mousemove', name: 'Mousemove', type: 'mouse-event' },
  { id: 'focus', name: 'Focus', type: 'focus-event' },
  { id: 'blur', name: 'Blur', type: 'focus-event' },
  { id: 'input', name: 'Input', type: 'input-event' },
  { id: 'change', name: 'Change', type: 'event' },
  { id: 'touchstart', name: 'Touch start', type: 'touch-event' },
  { id: 'touchend', name: 'Touch end', type: 'touch-event' },
  { id: 'touchmove', name: 'Touch move', type: 'touch-event' },
  { id: 'touchcancel', name: 'Touch cancel', type: 'touch-event' },
  { id: 'keydown', name: 'Keydown', type: 'keyboard-event' },
  { id: 'keyup', name: 'Keyup', type: 'keyboard-event' },
  { id: 'submit', name: 'Submit', type: 'submit-event' },
  { id: 'wheel', name: 'Wheel', type: 'wheel-event' },
];

export const dataExportTypes = [
  { name: 'JSON', id: 'json' },
  { name: 'CSV', id: 'csv' },
  { name: 'Plain text', id: 'plain-text' },
];

export const workflowCategories = {
  scrape: 'Scraping',
  automation: 'Automation',
  productivity: 'Productivity',
};

export const excludeOnError = [
  'note',
  'delay',
  'webhook',
  'trigger',
  'while-loop',
  'conditions',
  'block-package',
  'element-exists',
];

export const contentTypes = [
  { name: 'text/plain', value: 'text' },
  { name: 'application/json', value: 'json' },
  { name: 'multipart/form-data', value: 'form-data' },
  { name: 'application/x-www-form-urlencoded', value: 'form' },
];

export const supportLocales = [
  { id: 'en', name: 'English' },
  { id: 'fr', name: 'Français' },
  { id: 'it', name: 'Italiano' },
  { id: 'vi', name: 'Tiếng Việt' },
  { id: 'zh', name: '简体中文' },
  { id: 'zh-tw', name: '繁體中文' },
];

export const communities = [
  {
    name: 'GitHub',
    icon: 'riGithubFill',
    url: 'https://github.com/kholid060/automa',
  },
  {
    name: 'Twitter',
    icon: 'riTwitterLine',
    url: 'https://twitter.com/AutomaApp',
  },
  {
    name: 'Discord',
    icon: 'riDiscordLine',
    url: 'https://discord.gg/C6khwwTE84',
  },
  {
    name: 'YouTube',
    icon: 'riYoutubeLine',
    url: 'https://www.youtube.com/channel/UCL3qU64hW0fsIj2vOayOQUQ',
  },
];

export const elementsHighlightData = {
  selectedElements: {
    stroke: '#2563EB',
    activeStroke: '#f87171',
    fill: 'rgba(37, 99, 235, 0.1)',
    activeFill: 'rgba(248, 113, 113, 0.1)',
  },
  hoveredElements: {
    stroke: '#fbbf24',
    fill: 'rgba(251, 191, 36, 0.1)',
  },
};

export const excludeGroupBlocks = [
  'trigger',
  'repeat-task',
  'loop-data',
  'loop-breakpoint',
  'blocks-group',
  'conditions',
  'webhook',
  'element-exists',
  'while-loop',
  'block-package',
];

export const conditionBuilder = {
  valueTypes: [
    {
      id: 'value',
      category: 'value',
      name: 'Value',
      compareable: true,
      data: { value: '' },
    },
    {
      id: 'code',
      category: 'value',
      name: 'Code',
      compareable: false,
      data: { code: '\nreturn true;', context: 'background' },
    },
    {
      id: 'data#exists',
      category: 'value',
      name: 'Data exists',
      compareable: false,
      valueKey: 'dataPath',
      data: { dataPath: '' },
    },
    {
      id: 'element#text',
      category: 'element',
      name: 'Element text',
      compareable: true,
      data: { selector: '' },
    },
    {
      id: 'element#exists',
      category: 'element',
      name: 'Element exists',
      compareable: false,
      data: { selector: '' },
    },
    {
      id: 'element#notExists',
      category: 'element',
      name: 'Element not exists',
      compareable: false,
      data: { selector: '' },
    },
    {
      id: 'element#visible',
      category: 'element',
      name: 'Element visible',
      compareable: false,
      data: { selector: '' },
    },
    {
      id: 'element#visibleScreen',
      category: 'element',
      name: 'Element visible in screen',
      compareable: false,
      data: { selector: '' },
    },
    {
      id: 'element#invisible',
      category: 'element',
      name: 'Element hidden in screen',
      compareable: false,
      data: { selector: '' },
    },
    {
      id: 'element#attribute',
      category: 'element',
      name: 'Element attribute value',
      compareable: true,
      data: { selector: '', attrName: '' },
    },
  ],
  compareTypes: [
    { id: 'eq', name: 'Equals', needValue: true, category: 'basic' },
    {
      id: 'eqi',
      name: 'Equals (case insensitive)',
      needValue: true,
      category: 'basic',
    },
    { id: 'nq', name: 'Not equals', needValue: true, category: 'basic' },
    { id: 'gt', name: 'Greater than', needValue: true, category: 'number' },
    {
      id: 'gte',
      name: 'Greater than or equal',
      needValue: true,
      category: 'number',
    },
    { id: 'lt', name: 'Less than', needValue: true, category: 'number' },
    {
      id: 'lte',
      name: 'Less than or equal',
      needValue: true,
      category: 'number',
    },
    { id: 'cnt', name: 'Contains', needValue: true, category: 'text' },
    {
      id: 'cni',
      name: 'Contains (case insensitive)',
      needValue: true,
      category: 'text',
    },
    { id: 'nct', name: 'Not contains', needValue: true, category: 'text' },
    {
      id: 'nci',
      name: 'Not contains (case insensitive)',
      needValue: true,
      category: 'text',
    },
    { id: 'stw', name: 'Starts with', needValue: true, category: 'text' },
    { id: 'enw', name: 'Ends with', needValue: true, category: 'text' },
    { id: 'rgx', name: 'Match with RegEx', needValue: true, category: 'text' },
    { id: 'itr', name: 'Is truthy', needValue: false, category: 'boolean' },
    { id: 'ifl', name: 'Is falsy', needValue: false, category: 'boolean' },
  ],
  inputTypes: {
    selector: {
      placeholder: '.class',
      label: 'CSS selector or XPath',
    },
    value: {
      label: 'Value',
      placeholder: 'abc123',
    },
    attrName: {
      label: 'Attribute name',
      placeholder: 'name',
    },
    dataPath: {
      label: 'variables@variableName',
      placeholder: '',
    },
  },
};
