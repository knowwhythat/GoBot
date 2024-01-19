<template>
  <div style="height: 100vh; overflow: hidden" class="overflow-hidden">
    <Toolbar class="p-1">
      <template #start>
        <Button @click="router.back()" lable="返回" class="mr-2 px-3 py-2">
          <template #icon>
            <v-remixicon name="riArrowLeftCircleLine" />
          </template>
        </Button>
        <Button
          label="保存"
          v-tooltip="{
            value: shortcuts['editor:save'].readable,
            showDelay: 100,
            hideDelay: 300,
          }"
          class="mr-2 px-3 py-2"
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
          label="停止"
          v-tooltip="{
            value: shortcuts['editor:terminate-flow'].readable,
            showDelay: 100,
            hideDelay: 300,
          }"
          :disabled="!(running || debuging)"
          class="mr-2 px-3 py-2"
          @click="terminate"
        >
          <template #icon>
            <v-remixicon name="riStopLine" />
          </template>
        </Button>
        <Button
          label="运行"
          v-tooltip="{
            value: shortcuts['editor:execute-flow'].readable,
            showDelay: 100,
            hideDelay: 300,
          }"
          :disabled="debuging"
          class="mr-2 px-3 py-2"
          @click="run"
          :loading="running"
        >
          <template #icon>
            <v-remixicon name="riPlayLine" />
          </template>
        </Button>
        <Button
          label="调试"
          v-tooltip="{
            value: shortcuts['editor:debug-flow'].readable,
            showDelay: 100,
            hideDelay: 300,
          }"
          :disabled="running"
          class="mr-2 px-3 py-2"
          @click="debug"
          :loading="debuging"
        >
          <template #icon>
            <v-remixicon name="riBugLine" />
          </template>
        </Button>
        <Button
          label="单行调试"
          class="mr-2 px-3 py-2"
          :disabled="!breakpointHit"
          @click="dealDebug('n')"
          v-if="debuging"
        >
          <template #icon>
            <v-remixicon name="riSpeedMiniLine" />
          </template>
        </Button>
        <Button
          label="下个断点"
          class="mr-2 px-3 py-2"
          :disabled="!breakpointHit"
          @click="dealDebug('c')"
          v-if="debuging"
        >
          <template #icon>
            <v-remixicon name="riSkipForwardMiniLine" />
          </template>
        </Button>
        <Button
          label="重启Repl"
          icon="pi pi-replay"
          class="mr-2 px-3 py-2"
          @click="restartRepl"
        >
        </Button>
      </template>

      <template #center>
        <span class="flex">
          <p class="pt-3 pr-2 text-xl font-serif font-semibold">{{ label }}</p>
        </span>
      </template>

      <template #end>
        <ToggleButton
          v-model="activeSidePanel"
          on-label=""
          off-label=""
          class="px-3 py-2"
          on-icon="pi pi-chevron-right"
          off-icon="pi pi-chevron-left"
        >
        </ToggleButton>
      </template>
    </Toolbar>
    <Splitter id="sequence-designer" class="mb-5" stateStorage="local">
      <SplitterPanel :size="25">
        <TabView>
          <TabPanel header="组件库">
            <Tree
              :value="nodes"
              :filter="true"
              selectionMode="single"
              filterMode="lenient"
              filterPlaceholder="输入组件名称进行搜索"
              class="w-full h-full"
            >
              <template #default="slotProps">
                <div class="flex">
                  <p class="pr-2">
                    <v-remixicon
                      v-bind="getIconPath(slotProps.node.icon_path)"
                    />
                  </p>
                  <div
                    v-if="isLeaf(slotProps)"
                    draggable="true"
                    class="transform select-none cursor-move text-ellipsis max-w-xs"
                    @dragstart="
                      $event.dataTransfer.setData(
                        'activity',
                        JSON.stringify(slotProps.node)
                      )
                    "
                  >
                    {{ slotProps.node.label }}
                  </div>
                  <b v-else>{{ slotProps.node.label }}</b>
                </div>
              </template>
            </Tree>
          </TabPanel>
          <TabPanel header="项目树"></TabPanel>
        </TabView>
      </SplitterPanel>
      <SplitterPanel class="h-full" :size="75">
        <Splitter>
          <SplitterPanel :size="50">
            <div class="flex justify-around p-2">
              <SequenceActivity
                :element="mainActivity.sequence"
                @update="update"
              />
            </div>
          </SplitterPanel>
          <SplitterPanel v-show="activeSidePanel" :size="50">
            <TabView>
              <TabPanel header="运行日志">
                <Listbox
                  :filter="true"
                  emptyMessage="当前没有日志"
                  :options="logs"
                  @contextmenu="onLogRightClick"
                  class="w-full md:w-14rem"
                >
                </Listbox>
                <ContextMenu ref="menu" :model="items" />
              </TabPanel>
              <TabPanel header="参数库"></TabPanel>
              <TabPanel header="元素库"></TabPanel>
            </TabView>
          </SplitterPanel>
        </Splitter>
      </SplitterPanel>
    </Splitter>
  </div>
</template>
<script setup>
import ToggleButton from "primevue/togglebutton";
import TabView from "primevue/tabview";
import TabPanel from "primevue/tabpanel";
import ContextMenu from "primevue/contextmenu";
import Listbox from "primevue/listbox";
import Button from "primevue/button";
import Splitter from "primevue/splitter";
import SplitterPanel from "primevue/splitterpanel";
import { useRouter, onBeforeRouteLeave } from "vue-router";
import Toolbar from "primevue/toolbar";
import Tree from "primevue/tree";
import { ref, onMounted, reactive, provide, shallowReactive } from "vue";
import SequenceActivity from "@/components/activity/SequenceActivity.vue";
import { nanoid } from "nanoid";
import {
  ParseAllPlugin,
  GetSubFlow,
  SaveSubFlow,
  RunSubFlow,
  DebugSubFlow,
  TerminateSubFlow,
  DealDebugSignal,
  RestartReplCommand,
} from "@back/go/main/App";
import {
  EventsOn,
  EventsOff,
  WindowMinimise,
  WindowMaximise,
} from "@back/runtime/runtime";
import { useShortcut, getShortcut } from "@/composable/shortcut";
import { getIconPath } from "@/utils/helper";
import { useEditorStore } from "@/stores/editor";
import { useToast } from "primevue/usetoast";
const toast = useToast();
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
const nodes = ref(null);
const router = useRouter();
const dataChanged = ref(false);
const selectedActivity = ref([]);
provide("dataChanged", { dataChanged, updateDataChanged });
provide("contextVariable", { contextVariable: [] });
provide("projectId", { projectId: props.id });
provide("selectedActivity", { selectedActivity });

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

onMounted(() => {
  ParseAllPlugin()
    .then((result) => {
      nodes.value = result;
    })
    .catch((err) => {
      console.error(err);
      toast.add({
        severity: "error",
        summary: "加载组件树异常",
        detail: err,
        life: 3000,
      });
    });
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
});

function isLeaf(data) {
  return !(data.node.children && data.node.children.length > 0);
}

const activeSidePanel = ref(false);

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
    cutBlock.push(shallowReactive(children.splice(index, 1)[0]));
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
    copyBlock.push(shallowReactive({ ...block, id: nanoid(16) }));
  } else {
    children.forEach((ele) => {
      if (ele.children) {
        innerCopy(ele.children, id, copyBlock);
      }
    });
  }
}

function pasteBlock() {
  console.log(editorStore.copiedBlocks);
  if (editorStore.copiedBlocks.length > 0) {
    if (selectedActivity.value.length > 0) {
      const id = selectedActivity.value[selectedActivity.value.length - 1];
      innerPaste(mainActivity.sequence.children, id);
    } else {
      editorStore.copiedBlocks.forEach((block) => {
        mainActivity.sequence.children.push(
          shallowReactive({ ...block, id: nanoid(16) })
        );
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
      children.splice(
        index + innerIndex + 1,
        0,
        shallowReactive({ ...block, id: nanoid(16) })
      );
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
const running = ref(false);
const logs = ref([]);
var reg = /code_map_id="([^"]*)/;
async function run() {
  await save();
  running.value = true;
  EventsOn("log_event", (data) => {
    logs.value.push(data);
  });
  try {
    WindowMinimise();
    await RunSubFlow(props.id, props.subflowId);
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
  setTimeout(function () {
    EventsOff("log_event");
  }, 1000);
  WindowMaximise();
}
const debuging = ref(false);
const debugingId = ref(nanoid(16));
const breakpointHit = ref(false);
provide("debugingId", { debugingId });
async function debug() {
  await save();
  debuging.value = true;
  EventsOn("log_event", (data) => {
    logs.value.push(data);
  });
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
  EventsOff("log_event");
  EventsOff("debug");
  WindowMaximise();
}
function dealDebug(command) {
  breakpointHit.value = false;
  DealDebugSignal(command);
}
async function terminate() {
  await TerminateSubFlow(props.id, props.subflowId);
  debuging.value = false;
  running.value = false;
}

onBeforeRouteLeave(onBeforeLeave);

async function onBeforeLeave() {
  if (debuging.value || running.value) {
    const confirm = window.confirm("退出将停止正在运行的流程，是否退出");
    if (!confirm) {
      return false;
    } else {
      await terminate();
    }
  }
  if (dataChanged.value) {
    const confirm = window.confirm("有修改尚未保存，是否退出");
    if (!confirm) return false;
  }
}
const menu = ref();
const items = ref([
  {
    label: "清空",
    icon: "pi pi-trash",
    command: () => {
      logs.value.splice(0, logs.value.length);
    },
  },
]);
const onLogRightClick = (event) => {
  menu.value.show(event);
};

async function restartRepl() {
  await RestartReplCommand(props.id);
}
</script>
<style scoped>
#sequence-designer {
  height: calc(100vh - 56px);
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
:deep(.p-listbox-list) {
  height: calc(100vh - 186px);
}
:deep(.p-tree-wrapper) {
  height: calc(100vh - 168px);
}
:deep(.p-tree) {
  padding: 0.25rem;
}
:deep(.p-splitter) {
  margin: 0;
}
:deep(.p-treenode-icon) {
  display: none;
}
:deep(.p-treenode-label) {
  margin-left: -0.75rem;
}
</style>
