<template>
  <div style="height: 100vh; overflow: hidden" class="overflow-hidden">
    <Toolbar class="flex-none p-0" style="--wails-draggable: drag">
      <template #start>
        <Button
          text
          @click="confirmQuit"
          lable="返回"
          class="hover:bg-slate-200"
          v-tooltip.bottom="{
            value: '返回',
            showDelay: 100,
            hideDelay: 300,
          }"
        >
          <template #icon>
            <v-remixicon name="riArrowLeftCircleLine" />
          </template>
        </Button>
        <Button
          text
          label="保存"
          v-tooltip.bottom="{
            value: shortcuts['editor:save'].readable,
            showDelay: 100,
            hideDelay: 300,
          }"
          class="py-2 px-4 hover:bg-slate-200"
          @click="save"
        >
          <template #icon>
            <span>
              <span
                v-if="dataChanged"
                class="flex h-3 w-3 absolute top-0 left-0 -ml-1 -mt-1"
              >
                <span
                  class="animate-ping absolute inline-flex h-full w-full rounded-full bg-red-600 opacity-75"
                ></span>
                <span
                  class="relative inline-flex rounded-full h-3 w-3 bg-red-600"
                ></span>
              </span>
              <v-remixicon name="riSaveLine" />
            </span>
          </template>
        </Button>
        <Button
          text
          label="停止"
          v-tooltip.bottom="{
            value: shortcuts['editor:terminate-flow'].readable,
            showDelay: 100,
            hideDelay: 300,
          }"
          :disabled="!(running || debuging)"
          class="py-2 px-4 hover:bg-slate-200"
          @click="terminate"
        >
          <template #icon>
            <v-remixicon name="riStopLine" />
          </template>
        </Button>
        <Button
          text
          label="运行"
          v-tooltip.bottom="{
            value: shortcuts['editor:execute-flow'].readable,
            showDelay: 100,
            hideDelay: 300,
          }"
          :disabled="debuging"
          class="py-2 px-4 hover:bg-slate-200"
          @click="run"
          :loading="running"
        >
          <template #icon>
            <v-remixicon name="riPlayLine" />
          </template>
        </Button>
        <Button
          text
          label="调试"
          v-tooltip.bottom="{
            value: shortcuts['editor:debug-flow'].readable,
            showDelay: 100,
            hideDelay: 300,
          }"
          :disabled="running"
          class="py-2 px-4 hover:bg-slate-200"
          @click="debug"
          :loading="debuging"
        >
          <template #icon>
            <v-remixicon name="riBugLine" />
          </template>
        </Button>
        <Button
          text
          label="单行调试"
          class="py-2 px-4 hover:bg-slate-200"
          :disabled="!breakpointHit"
          @click="dealDebug('n')"
          v-if="debuging"
        >
          <template #icon>
            <v-remixicon name="riSpeedMiniLine" />
          </template>
        </Button>
        <Button
          text
          label="下个断点"
          class="py-2 px-4 hover:bg-slate-200"
          :disabled="!breakpointHit"
          @click="dealDebug('c')"
          v-if="debuging"
        >
          <template #icon>
            <v-remixicon name="riSkipForwardMiniLine" />
          </template>
        </Button>
        <Button
          text
          label="重启Repl"
          icon="pi pi-replay"
          class="py-2 px-4 hover:bg-slate-200"
          @click="restartRepl"
        >
        </Button>
      </template>

      <template #end>
        <SystemOperate @quit="confirmQuit" />
      </template>
    </Toolbar>
    <div class="flex-1">
      <splitpanes
        class="default-theme"
        :dblClickSplitter="false"
        :push-other-panes="false"
      >
        <SplitterPanel
          :size="15"
          v-show="tabs[activeIndex]?.nodeType === 'sequence'"
        >
          <div class="m-1">
            <LeftPaneView />
          </div>
        </SplitterPanel>
        <SplitterPanel
          :size="tabs[activeIndex]?.nodeType === 'sequence' ? 70 : 85"
        >
          <div class="m-1">
            <splitpanes
              class="default-theme"
              :dblClickSplitter="false"
              id="sequence-designer"
              horizontal
              @resized="splitPaneResize"
            >
              <SplitterPanel :size="70" :min-size="30">
                <TabView v-model:activeIndex="activeIndex" :scrollable="true">
                  <TabPanel v-for="(tab, index) in tabs" :key="tab.subflowId">
                    <template #header>
                      <div class="flex align-items-center gap-2">
                        <span class="font-bold text-nowrap">
                          {{ tab.title }}
                        </span>
                        <div
                          v-if="index !== 0"
                          class="hover:bg-slate-200 -mt-px"
                          @click.stop="closeTab(tab.subflowId)"
                        >
                          <v-remixicon
                            size="20"
                            viewBox="0 0 1024 1024"
                            name="closeIcon"
                          />
                        </div>
                      </div>
                    </template>
                    <VisualFlow
                      @mounted="visualFlowMounted"
                      v-if="tab.nodeType === 'sequence'"
                      ref="flowTabs"
                      :id="tab.id"
                      :subflowId="tab.subflowId"
                      :label="tab.title"
                    />
                    <WorkflowView
                      :id="id"
                      v-else-if="tab.nodeType === 'flow'"
                      ref="flowTabs"
                      @edit:block="editBlock"
                      @delete:block="deleteBlock"
                    />
                  </TabPanel>
                </TabView>
              </SplitterPanel>
              <SplitterPanel
                v-if="activeNav !== ''"
                :size="bottomHeight"
                :min-size="10"
              >
                <LogsView
                  v-if="activeNav === 'LogsView'"
                  :logs="logs"
                  @hide="activeNav = ''"
                  @clear="logs = []"
                />
                <ElementsView
                  :id="props.id"
                  :windowsElement="windowsElement"
                  v-else-if="activeNav === 'ElementsView'"
                  @hide="activeNav = ''"
                />
                <ParamsView
                  :id="props.id"
                  v-else-if="activeNav === 'ParamsView'"
                  @hide="activeNav = ''"
                  v-model="params"
                />
                <LocalVariablesView
                  v-else-if="activeNav === 'LocalVariablesView'"
                  @hide="activeNav = ''"
                />
              </SplitterPanel>
            </splitpanes>
            <div class="flex flex-row h-8 bg-white">
              <div
                v-for="nav in bottomNav"
                @click="bottomNavClick(nav.component)"
                class="text-center w-20 my-auto py-1 hover:text-indipurplego-500"
                :class="[
                  nav.component === activeNav
                    ? 'border-b-4 border-purple-500 text-purple-500'
                    : '',
                ]"
              >
                {{ nav.title }}
              </div>
            </div>
          </div>
        </SplitterPanel>
        <SplitterPanel :size="15" id="right-pane">
          <div class="m-1 h-full flex flex-col">
            <Panel class="h-1/2">
              <template #header>
                <div class="flex flex-1">
                  <div class="flex-none">流程</div>
                  <div class="flex-1 flex flex-row-reverse -mt-1 gap-4">
                    <div
                      class="hover:bg-slate-200 px-1"
                      v-tooltip.left="'新建可视化流程'"
                      @click="addNewVisualFlow"
                    >
                      <v-remixicon size="20" name="riNodeTree" />
                    </div>
                    <div
                      class="hover:bg-slate-200 px-1"
                      v-tooltip.left="'新建代码流程'"
                      @click="addNewVisualFlow"
                    >
                      <v-remixicon size="20" name="riCodeBoxLine" />
                    </div>
                    <div
                      class="hover:bg-slate-200 px-1"
                      v-tooltip.left="'依赖管理'"
                      @click="pipDialogVisible = true"
                    >
                      <v-remixicon
                        size="20"
                        viewBox="0 0 1024 1024"
                        name="pythonIcon"
                      />
                    </div>
                  </div>
                </div>
              </template>
              <Tree
                v-model:expandedKeys="expandedKeys"
                filterPlaceholder="搜索"
                :value="projectTreeData"
                selectionMode="single"
                filterMode="lenient"
                class="w-full"
                :pt="{
                  content: (options) => ({
                    oncontextmenu: onContextMenuClick(options),
                  }),
                }"
              >
                <template #default="slotProps">
                  <div
                    class="truncate max-w-48"
                    @dblclick="openTab(slotProps.node)"
                  >
                    {{ slotProps.node.label }}
                  </div>
                </template>
              </Tree>
              <ContextMenu ref="menu" :model="items" />
              <Dialog
                v-model:visible="renameDialogVisible"
                modal
                header="重命名"
                :style="{ width: '35rem' }"
              >
                <div class="flex items-center gap-3 mb-5">
                  <InputText
                    v-model="tempLable"
                    id="email"
                    class="flex-auto"
                    autocomplete="off"
                  />
                </div>
                <div class="flex justify-end gap-2">
                  <Button
                    type="button"
                    label="取消"
                    severity="secondary"
                    @click="renameDialogVisible = false"
                  ></Button>
                  <Button
                    type="button"
                    label="确定"
                    @click="doRenameFlowName"
                  ></Button>
                </div>
              </Dialog>
              <Dialog
                v-model:visible="pipDialogVisible"
                modal
                header="依赖管理"
                :style="{ width: '48rem' }"
              >
                <PipInstallView @close="pipDialogVisible = false" />
              </Dialog>
            </Panel>
            <div class="h-1/2">
              <GlobalVariableView v-model="projectConfig.variables" />
            </div>
          </div>
        </SplitterPanel>
      </splitpanes>
    </div>
  </div>
</template>
<script setup>
import Button from "primevue/button";
import TabView from "primevue/tabview";
import TabPanel from "primevue/tabpanel";
import Tree from "primevue/tree";
import ContextMenu from "primevue/contextmenu";
import Dialog from "primevue/dialog";
import { Pane as SplitterPanel, Splitpanes } from "splitpanes";
import "splitpanes/dist/splitpanes.css";
import { useRouter } from "vue-router";
import Toolbar from "primevue/toolbar";
import {
  computed,
  nextTick,
  onBeforeMount,
  onMounted,
  onUnmounted,
  provide,
  ref,
  watch,
} from "vue";
import SystemOperate from "@/components/SystemOperate.vue";
import LeftPaneView from "@/views/design/sequence/LeftPaneView.vue";
import LogsView from "@/views/design/sequence/LogsView.vue";
import LocalVariablesView from "@/views/design/sequence/LocalVariablesView.vue";
import ElementsView from "@/views/design/sequence/ElementsView.vue";
import ParamsView from "@/views/design/sequence/ParamsView.vue";
import VisualFlow from "@/views/design/sequence/VisualFlow.vue";
import GlobalVariableView from "@/views/design/sequence/GlobalVariableView.vue";
import PipInstallView from "@/views/design/sequence/PipInstallView.vue";
import WorkflowView from "@/views/design/flow/WorkflowView.vue";
import { customAlphabet, nanoid } from "nanoid";
import {
  ReadProjectConfig,
  DealDebugSignal,
  DebugSubFlow,
  GetProjectWindowsElements,
  RestartReplCommand,
  RunSubFlow,
  TerminateSubFlow,
  SaveProjectConfig,
  DeleteSubFlow,
} from "@back/go/backend/App";
import {
  EventsOff,
  EventsOn,
  WindowMaximise,
  WindowMinimise,
} from "@back/runtime/runtime";
import { getShortcut, useShortcut } from "@/composable/shortcut";

import { useConfirm } from "primevue/useconfirm";
import { useToast } from "primevue/usetoast";

const toast = useToast();
const confirm = useConfirm();
const router = useRouter();

const props = defineProps({
  id: {
    type: String,
    default: "",
  },
  subflowId: {
    type: String,
    default: "",
  },
  label: {
    type: String,
    default: "",
  },
});

provide("projectId", { projectId: props.id });

const dataChanged = ref(false);

const logs = ref([]);

const expandedKeys = ref({});
const activeIndex = ref(0);
const flowTabs = ref();

const generateSubflowId = customAlphabet("abcdefghijklmnopqrstuvwxyz", 12);

const projectConfig = ref({});

const tabs = computed(() => {
  let openFlows = [];
  if (projectConfig.value.nodeType === "flow") {
    openFlows.push({
      id: props.id,
      subflowId: projectConfig.value.key,
      title: "主流程",
      nodeType: projectConfig.value.nodeType,
    });
  }
  if (projectConfig.value && projectConfig.value.children) {
    projectConfig.value.children.forEach((config) => {
      if (config.opened) {
        openFlows.push({
          id: props.id,
          subflowId: config.key,
          title: config.label,
          nodeType: config.nodeType,
        });
      }
    });
  }

  return openFlows;
});

const projectTreeData = computed(() => {
  return [
    {
      key: projectConfig.value.key,
      label: projectConfig.value.label,
      children: projectConfig.value.children,
    },
  ];
});

function addNewVisualFlow() {
  if (projectConfig.value.children) {
    projectConfig.value.children.push({
      key: generateSubflowId(8),
      label: `子流程${tabs.value.length}`,
      nodeType: "sequence",
      opened: true,
      children: [],
    });
    activeIndex.value = tabs.value.length - 1;
  }
}

function openTab(node) {
  node.opened = true;
  activeIndex.value = tabs.value.findIndex((tab) => tab.subflowId === node.key);
}

function editBlock(data) {
  const index = projectConfig.value.children.findIndex(
    (v) => v.key === data.blockId
  );
  if (index === -1) {
    projectConfig.value.children.push({
      key: data.blockId,
      label: data.data.label ? data.data.label : `子流程${tabs.value.length}`,
      nodeType: "sequence",
      opened: true,
      children: [],
    });
  } else {
    projectConfig.value.children[index].opened = true;
    if (data.data.label) {
      projectConfig.value.children[index].label = data.data.label;
    }
  }
  activeIndex.value = tabs.value.findIndex(
    (tab) => tab.subflowId === data.blockId
  );
}

async function deleteBlock(blockId) {
  let current = projectConfig.value.children.findIndex(
    (child) => child.key === blockId
  );
  if (current > -1) {
    projectConfig.value.children = projectConfig.value.children.filter(
      (child) => child.key !== blockId
    );
    await SaveProjectConfig(props.id, JSON.stringify(projectConfig.value));
    await DeleteSubFlow(props.id, selectedNode.value.key);
  }
}

function closeTab(id) {
  let current = 0;
  projectConfig.value.children.forEach((config, index) => {
    if (config.key === id) {
      config.opened = false;
      current = index;
    }
  });
  activeIndex.value = current - 1;
}
const selectedNode = ref(null);
const renameDialogVisible = ref(false);
const tempLable = ref("");
const menu = ref();
const items = ref([
  {
    label: "编辑",
    icon: "pi pi-file-edit",
    command: () => {
      renameDialogVisible.value = true;
      tempLable.value = selectedNode.value.label;
    },
  },
  {
    label: "删除",
    icon: "pi pi-times-circle",
    command: () => {
      confirm.require({
        header: "提示",
        message: "确定要删除这条记录吗?",
        icon: "pi pi-info-circle",
        rejectClass: "p-button-secondary p-button-outlined p-button-sm",
        acceptClass: "p-button-danger p-button-sm",
        rejectLabel: "取消",
        acceptLabel: "删除",
        accept: () => {
          if (selectedNode.value.key === "main") {
            toast.add({
              severity: "warn",
              detail: "主流程不能删除",
              life: 1000,
            });
            return;
          }
          let current = projectConfig.value.children.findIndex(
            (child) => child.key === selectedNode.value.key
          );
          projectConfig.value.children = projectConfig.value.children.filter(
            (child) => child.key !== selectedNode.value.key
          );
          SaveProjectConfig(props.id, JSON.stringify(projectConfig.value));
          DeleteSubFlow(props.id, selectedNode.value.key);
          activeIndex.value = current - 1;
        },
        reject: () => {},
      });
    },
  },
]);

function doRenameFlowName() {
  selectedNode.value.label = tempLable.value;
  renameDialogVisible.value = false;
}
function onContextMenuClick(options) {
  return async function (event) {
    if (options.props.level !== 1) {
      selectedNode.value = options.props.node;
      menu.value.show(event);
    }
  };
}

const pipDialogVisible = ref(false);

const params = ref([]);

function visualFlowMounted() {
  dataChanged.value = flowTabs.value[activeIndex.value]?.isChanged();
  params.value = flowTabs.value[activeIndex.value]?.getParams();
}

watch(
  () => activeIndex.value,
  (now, old) => {
    dataChanged.value = flowTabs.value[now]?.isChanged();
    params.value = flowTabs.value[now]?.getParams();
  },
  { deep: true }
);

watch(
  () => params.value,
  (now, old) => {
    flowTabs.value[activeIndex.value]?.setParams(now);
  },
  { deep: true }
);

onBeforeMount(async () => {
  WindowMaximise();
  projectConfig.value = await ReadProjectConfig(props.id);
  expandedKeys.value[props.id] = true;
  EventsOn("log_event", (data) => {
    logs.value.push(data);
  });
  await loadElements();
  EventsOn("windows_element_change", async () => {
    await loadElements();
  });
});

onUnmounted(() => {
  EventsOff("log_event");
  EventsOff("windows_element_change");
});

const windowsElement = ref([]);
provide("windowsElement", { id: props.id, windowsElement: windowsElement });

async function loadElements() {
  const result = await GetProjectWindowsElements(props.id);
  windowsElement.value = [];
  if (result) {
    const elements = JSON.parse(result);
    for (let [key, value] of Object.entries(elements)) {
      const processName = value["processName"];
      const process = windowsElement.value.find(
        (node) => node.label === processName
      );
      if (process) {
        process["children"].push({
          key: key,
          label: value["displayName"],
          data: value,
        });
      } else {
        windowsElement.value.push({
          key: nanoid(16),
          label: value["processName"],
          data: [],
          children: [
            {
              key: key,
              label: value["displayName"],
              data: value,
            },
          ],
        });
      }
    }
  }
}

const dragBlockId = ref("");
const dropBlockId = ref("");

provide("dragBlockId", { dragBlockId });
provide("dropBlockId", { dropBlockId });

const shortcuts = useShortcut([
  getShortcut("editor:save", save),
  getShortcut("editor:execute-flow", run),
  getShortcut("editor:debug-flow", debug),
  getShortcut("editor:terminate-flow", terminate),
  getShortcut("editor:del-block", delBlock),
  getShortcut("editor:cut-block", cutBlock),
  getShortcut("editor:copy-block", copyBlock),
  getShortcut("editor:paste-block", pasteBlock),
]);

function delBlock() {
  if (tabs.value[activeIndex.value].nodeType === "sequence") {
    flowTabs.value[activeIndex.value].delBlock();
  }
}

function cutBlock() {
  if (tabs.value[activeIndex.value].nodeType === "sequence") {
    flowTabs.value[activeIndex.value].cutBlock();
  }
}

function copyBlock() {
  if (tabs.value[activeIndex.value].nodeType === "sequence") {
    flowTabs.value[activeIndex.value].copyBlock();
  }
}

function pasteBlock() {
  if (tabs.value[activeIndex.value].nodeType === "sequence") {
    flowTabs.value[activeIndex.value].pasteBlock();
  }
}

function save() {
  flowTabs.value[activeIndex.value].save();
  SaveProjectConfig(props.id, JSON.stringify(projectConfig.value));
}

//底部工具栏逻辑
const bottomNav = [
  { title: "日志", component: "LogsView" },
  { title: "元素库", component: "ElementsView" },
  { title: "参数库", component: "ParamsView" },
  { title: "变量库", component: "LocalVariablesView" },
];

const activeNav = ref("");

function bottomNavClick(nav) {
  if (activeNav.value !== nav) {
    activeNav.value = nav;
  } else {
    activeNav.value = "";
  }
}

const bottomHeight = ref(30);

function splitPaneResize(e) {
  if (e.length === 2) {
    bottomHeight.value = e[1].size;
  }
}

const running = ref(false);

const reg = /code_map_id="([^"]*)/;

async function run() {
  save();
  await restartRepl();
  logs.value = [];
  running.value = true;
  try {
    WindowMinimise();
    await RunSubFlow(props.id, tabs.value[activeIndex.value].subflowId);
    debugingId.value = "";
    toast.add({
      severity: "success",
      summary: "提示",
      detail: "运行成功",
      life: 3000,
    });
  } catch (err) {
    const match = err.match(reg);
    if (match && match.length > 1) {
      debugingId.value = match[1];
      // document.getElementById(match[1]).scrollIntoView();
    }
    toast.add({
      severity: "error",
      summary: "运行失败",
      detail: err,
      life: 80000,
    });
  }
  running.value = false;
  WindowMaximise();
}

const debuging = ref(false);
const debugingId = ref(nanoid(16));
const breakpointHit = ref(false);

provide("debugingId", { debugingId });

async function debug() {
  await restartRepl();
  save();
  debuging.value = true;
  EventsOn("debug", (data) => {
    debugingId.value = data;
    breakpointHit.value = true;
  });
  try {
    WindowMinimise();
    await DebugSubFlow(props.id, tabs.value[activeIndex.value].subflowId);
    debugingId.value = nanoid(16);
    toast.add({
      severity: "success",
      summary: "提示",
      detail: "运行成功",
      life: 3000,
    });
  } catch (err) {
    toast.add({
      severity: "error",
      summary: "运行失败",
      detail: err,
      life: 80000,
    });
  }
  debuging.value = false;
  EventsOff("debug");
  WindowMaximise();
}

function dealDebug(command) {
  breakpointHit.value = false;
  DealDebugSignal(command);
}

async function terminate() {
  await TerminateSubFlow();
  debuging.value = false;
  running.value = false;
}

const confirmQuit = () => {
  let isChanged = flowTabs.value.some((tab) => tab.isChanged());
  if (!isChanged && !debuging.value && !running.value) {
    terminate();
    router.back();
  } else if (debuging.value || running.value) {
    confirm.require({
      message: "退出将停止正在运行的流程，是否退出?",
      header: "确认",
      icon: "pi pi-exclamation-triangle",
      rejectClass: "p-button-secondary p-button-outlined",
      rejectLabel: "取消",
      acceptLabel: "确认",
      accept: async () => {
        await terminate();
        router.back();
      },
      reject: () => {},
    });
  } else {
    confirm.require({
      message: "工作空间的修改尚未保存,确认将丢弃所有修改,是否确认?",
      header: "确认",
      icon: "pi pi-exclamation-triangle",
      rejectClass: "p-button-secondary p-button-outlined",
      rejectLabel: "取消",
      acceptLabel: "确认",
      accept: () => {
        router.back();
      },
      reject: () => {},
    });
  }
};

async function restartRepl() {
  await RestartReplCommand(props.id);
}
</script>
<style scoped lang="scss">
#sequence-designer {
  height: calc(100vh - 96px);

  .splitpanes__pane {
    overflow-y: auto;
  }
}
#right-pane {
  height: calc(100vh - 54px);
}
:deep(.p-tabview-panels) {
  padding: 0;
}

:deep(.p-tabview-nav-link) {
  padding: 0.75rem;
}

:deep(.p-panel-header) {
  padding: 0.75rem;
}
:deep(.p-panel-content) {
  padding: 0;
}
:deep(.p-tree) {
  padding: 0;
  border: none;
}
</style>
