<template>
  <div style="height: 100vh; overflow: hidden" class="overflow-hidden">
    <Toolbar class="p-1 border-none" style="--wails-draggable: drag">
      <template #start>
        <Button
          @click="confirmQuit"
          lable="返回"
          class="mr-2 px-3 py-2 hover:bg-slate-200"
          text
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
          class="mr-2 px-3 py-2 hover:bg-slate-200"
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
          class="mr-2 px-3 py-2 hover:bg-slate-200"
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
          class="mr-2 px-3 py-2 hover:bg-slate-200"
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
          class="mr-2 px-3 py-2 hover:bg-slate-200"
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
          class="mr-2 px-3 py-2 hover:bg-slate-200"
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
          class="mr-2 px-3 py-2 hover:bg-slate-200"
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
          class="mr-2 px-3 py-2 hover:bg-slate-200"
          @click="restartRepl"
        >
        </Button>
      </template>

      <template #end>
        <SystemOperate @quit="confirmQuit" />
      </template>
    </Toolbar>
    <Splitter class="mb-5" stateStorage="local" :gutter-size="5">
      <SplitterPanel :size="25"> <LeftPaneView /> </SplitterPanel>
      <SplitterPanel class="h-full" :size="75" :min-size="66">
        <div>
          <Splitter id="sequence-designer" layout="vertical" :gutter-size="5">
            <SplitterPanel :size="70" :min-size="30">
              <div class="flex justify-around">
                <SequenceActivity
                  :element="mainActivity.sequence"
                  @update="update"
                />
              </div>
            </SplitterPanel>
            <SplitterPanel v-if="activeNav != ''" :size="30" :min-size="10">
              <LogsView
                v-if="activeNav == 'LogsView'"
                :logs="logs"
                @hide="activeNav = ''"
                @clear="logs = []"
              />
              <VariablesView
                v-else-if="activeNav == 'VariablesView'"
                @hide="activeNav = ''"
              />
              <ElementsView
                :id="props.id"
                v-else-if="activeNav == 'ElementsView'"
                @hide="activeNav = ''"
              />
            </SplitterPanel>
          </Splitter>
          <div class="flex flex-row h-8 bg-slate-100">
            <div
              v-for="nav in bottomNav"
              @click="bottomNavClick(nav.component)"
              class="text-center w-20 my-auto hover:text-indipurplego-500"
              :class="[
                nav.component == activeNav
                  ? 'border-b-4 border-purple-500 text-purple-500'
                  : '',
              ]"
            >
              {{ nav.title }}
            </div>
          </div>
        </div>
      </SplitterPanel>
    </Splitter>
  </div>
</template>
<script setup>
import Button from "primevue/button";
import Splitter from "primevue/splitter";
import SplitterPanel from "primevue/splitterpanel";
import { useRouter } from "vue-router";
import Toolbar from "primevue/toolbar";
import { ref, onMounted, onUnmounted, reactive, provide } from "vue";
import SequenceActivity from "@/components/activity/SequenceActivity.vue";
import SystemOperate from "@/components/SystemOperate.vue";
import LeftPaneView from "@/views/design/sequence/LeftPaneView.vue";
import LogsView from "@/views/design/sequence/LogsView.vue";
import VariablesView from "@/views/design/sequence/VariablesView.vue";
import ElementsView from "@/views/design/sequence/ElementsView.vue";
import { nanoid } from "nanoid";
import {
  GetSubFlow,
  SaveSubFlow,
  RunSubFlow,
  DebugSubFlow,
  TerminateSubFlow,
  DealDebugSignal,
  RestartReplCommand,
} from "@back/go/backend/App";
import {
  EventsOn,
  EventsOff,
  WindowMinimise,
  WindowMaximise,
} from "@back/runtime/runtime";
import { useShortcut, getShortcut } from "@/composable/shortcut";

import { useEditorStore } from "@/stores/editor";
import { useConfirm } from "primevue/useconfirm";
import { useToast } from "primevue/usetoast";
const toast = useToast();
const confirm = useConfirm();

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

const router = useRouter();
const dataChanged = ref(false);
const selectedActivity = ref([]);
provide("dataChanged", { dataChanged, updateDataChanged });
provide("contextVariable", { contextVariable: [] });
provide("projectId", { projectId: props.id });
provide("selectedActivity", { selectedActivity });
const logs = ref([]);

function updateDataChanged() {
  dataChanged.value = true;
}

const mainActivity = reactive({
  sequence: {
    id: props.subflowId,
    key: "main",
    label: props.label,
    toggleable: false,
    deleteable: false,
    icon_path: "riHome5Line",
    children: [],
  },
});

function update({ children }) {
  dataChanged.value = true;
  mainActivity.sequence.children = children;
}

onMounted(async () => {
  await WindowMaximise();
  GetSubFlow(props.id, props.subflowId)
    .then((result) => {
      if (result) {
        const data = JSON.parse(result);
        mainActivity.sequence = reactive(data.sequence);
      } else {
        mainActivity.sequence.id = props.subflowId;
        mainActivity.sequence.label = props.label ? props.label : "子流程";
      }
    })
    .catch((err) => {
      toast.add({
        severity: "error",
        summary: "加载流程异常",
        detail: err,
        life: 3000,
      });
    });
  EventsOn("log_event", (data) => {
    logs.value.push(data);
  });
});

onUnmounted(() => {
  EventsOff("log_event");
});

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
  selectedActivity.value.forEach((id) => {
    innerDelete(mainActivity.sequence.children, id);
  });
  selectedActivity.value = [];
}

function innerDelete(children, id) {
  const index = children.findIndex(
    (item) => item.id === id && !item.hasOwnProperty("deleteable")
  );
  if (index != -1) {
    children.splice(index, 1);
    dataChanged.value = true;
  } else {
    children.forEach((ele) => {
      if (ele.children) {
        innerDelete(ele.children, id);
      }
    });
  }
}
const editorStore = useEditorStore();
function cutBlock() {
  let cutBlock = [];
  selectedActivity.value.forEach((id) => {
    innerCut(mainActivity.sequence.children, id, cutBlock);
  });
  selectedActivity.value = [];
  editorStore.addToPasteBlocks(cutBlock);
}

function innerCut(children, id, cutBlock) {
  const index = children.findIndex(
    (item) => item.id === id && !item.hasOwnProperty("deleteable")
  );
  if (index != -1) {
    cutBlock.push(children.splice(index, 1)[0]);
    dataChanged.value = true;
  } else {
    children.forEach((ele) => {
      if (ele.children) {
        innerCut(ele.children, id, cutBlock);
      }
    });
  }
}

function copyBlock() {
  let copyBlock = [];
  selectedActivity.value.forEach((id) => {
    innerCopy(mainActivity.sequence.children, id, copyBlock);
  });
  editorStore.addToPasteBlocks(copyBlock);
}

function innerCopy(children, id, copyBlock) {
  const block = children.find(
    (item) => item.id === id && !item.hasOwnProperty("deleteable")
  );
  if (block) {
    copyBlock.push({ ...block, id: nanoid(16) });
  } else {
    children.forEach((ele) => {
      if (ele.children) {
        innerCopy(ele.children, id, copyBlock);
      }
    });
  }
}

function pasteBlock() {
  if (editorStore.copiedBlocks.length > 0) {
    if (selectedActivity.value.length > 0) {
      const id = selectedActivity.value[selectedActivity.value.length - 1];
      innerPaste(mainActivity.sequence.children, id);
    } else {
      editorStore.copiedBlocks.forEach((block) => {
        mainActivity.sequence.children.push({ ...block, id: nanoid(16) });
      });
      dataChanged.value = true;
    }
  }
}

function innerPaste(children, id) {
  const index = children.findIndex(
    (item) => item.id === id && !item.hasOwnProperty("deleteable")
  );
  if (index != -1) {
    editorStore.copiedBlocks.forEach((block, innerIndex) => {
      children.splice(index + innerIndex + 1, 0, { ...block, id: nanoid(16) });
    });
    dataChanged.value = true;
  } else {
    children.forEach((ele) => {
      if (ele.children) {
        innerPaste(ele.children, id);
      }
    });
  }
}

async function save() {
  try {
    await SaveSubFlow(props.id, props.subflowId, JSON.stringify(mainActivity));
    dataChanged.value = false;
    toast.add({
      severity: "success",
      summary: "提示",
      detail: "保存成功",
      life: 3000,
    });
  } catch (err) {
    toast.add({
      severity: "error",
      summary: "保存失败",
      detail: err,
      life: 3000,
    });
  }
}

//底部工具栏逻辑
const bottomNav = [
  { title: "日志", component: "LogsView" },
  { title: "变量库", component: "VariablesView" },
  { title: "元素库", component: "ElementsView" },
];

const activeNav = ref("");

function bottomNavClick(nav) {
  if (activeNav.value != nav) {
    activeNav.value = nav;
  } else {
    activeNav.value = "";
  }
}

const running = ref(false);

var reg = /code_map_id="([^"]*)/;
async function run() {
  await save();
  await restartRepl();
  logs.value = [];
  running.value = true;
  try {
    WindowMinimise();
    await RunSubFlow(props.id, props.subflowId);
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
      document.getElementById(match[1]).scrollIntoView();
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
    await DebugSubFlow(props.id, props.subflowId);
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
<style scoped>
#sequence-designer {
  height: calc(100vh - 86px);
}
:deep(.p-splitter .p-splitter-gutter) {
  background: #e5e7eb;
}
:deep(.p-splitter-panel) {
  overflow: auto;
}
:deep(.p-tabview-panels) {
  padding: 0rem;
}
:deep(.p-tabview-nav-link) {
  padding: 0.75rem;
}
:deep(.p-splitter) {
  margin: 0;
}
:deep(.p-panel .p-panel-content) {
  padding: 0.25rem;
}
</style>
