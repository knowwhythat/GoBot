<template>
  <div style="height: 100vh; overflow: hidden" class="overflow-hidden">
    <Toolbar class="p-0" style="--wails-draggable: drag">
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
    <splitpanes class="default-theme" :dblClickSplitter="false">
      <SplitterPanel :size="15">
        <div class="m-1">
          <LeftPaneView />
        </div>
      </SplitterPanel>
      <SplitterPanel :size="70" :min-size="50">
        <div class="m-1">
          <splitpanes
            class="default-theme"
            :dblClickSplitter="false"
            id="sequence-designer"
            horizontal
            @resized="splitPaneResize"
          >
            <SplitterPanel :size="70" :min-size="30">
              <TabView
                :scrollable="true"
                :activeIndex="activeIndex"
                @update:activeIndex="tabChange"
              >
                <TabPanel
                  v-for="tab in tabs"
                  :key="tab.subflowId"
                  :header="tab.title"
                >
                  <VisualFlow
                    ref="flowTabs"
                    :id="tab.id"
                    :subflowId="tab.subflowId"
                    :label="tab.title"
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
              <VariablesView
                v-else-if="activeNav === 'VariablesView'"
                @hide="activeNav = ''"
              />
              <ElementsView
                :id="props.id"
                :windowsElement="windowsElement"
                v-else-if="activeNav === 'ElementsView'"
                @hide="activeNav = ''"
              />
            </SplitterPanel>
          </splitpanes>
          <div class="flex flex-row h-8 bg-white">
            <div
              v-for="nav in bottomNav"
              @click="bottomNavClick(nav.component)"
              class="text-center w-20 my-auto hover:text-indipurplego-500"
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
      <SplitterPanel :size="15" max-size="30">
        <div class="m-1">
          <Panel header="流程"></Panel>
        </div>
      </SplitterPanel>
    </splitpanes>
  </div>
</template>
<script setup>
import Button from "primevue/button";
import TabView from "primevue/tabview";
import TabPanel from "primevue/tabpanel";
import { Pane as SplitterPanel, Splitpanes } from "splitpanes";
import "splitpanes/dist/splitpanes.css";
import { useRouter } from "vue-router";
import Toolbar from "primevue/toolbar";
import { onMounted, onUnmounted, provide, ref } from "vue";
import SystemOperate from "@/components/SystemOperate.vue";
import LeftPaneView from "@/views/design/sequence/LeftPaneView.vue";
import LogsView from "@/views/design/sequence/LogsView.vue";
import VariablesView from "@/views/design/sequence/VariablesView.vue";
import ElementsView from "@/views/design/sequence/ElementsView.vue";
import VisualFlow from "@/views/design/sequence/VisualFlow.vue";
import { nanoid } from "nanoid";
import {
  DealDebugSignal,
  DebugSubFlow,
  GetProjectWindowsElements,
  RestartReplCommand,
  RunSubFlow,
  TerminateSubFlow,
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

const tabs = ref([]);
const activeIndex = ref(0);
const flowTabs = ref();

function tabChange(index) {
  dataChanged.value = flowTabs.value[index].isChanged();
  //Todo  参数，变量
}

onMounted(async () => {
  await WindowMaximise();
  tabs.value.push({
    id: props.id,
    subflowId: props.subflowId,
    title: props.label,
  });
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
        (node) => node.label === processName,
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
  flowTabs.value[activeIndex.value].delBlock();
}

function cutBlock() {
  flowTabs.value[activeIndex.value].cutBlock();
}

function copyBlock() {
  flowTabs.value[activeIndex.value].copyBlock();
}

function pasteBlock() {
  flowTabs.value[activeIndex.value].pasteBlock();
}

function save() {
  flowTabs.value[activeIndex.value].save();
}

//底部工具栏逻辑
const bottomNav = [
  { title: "日志", component: "LogsView" },
  { title: "元素库", component: "ElementsView" },
  { title: "变量库", component: "VariablesView" },
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
  await save();
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
  await save();
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
  if (!dataChanged.value && !debuging.value && !running.value) {
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
  height: calc(100vh - 88px);

  .splitpanes__pane {
    overflow-y: auto;
  }
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
</style>
