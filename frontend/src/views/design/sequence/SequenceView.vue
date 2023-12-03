<template>
  <div style="height: 100vh; overflow: hidden" class="overflow-hidden">
    <Toolbar class="p-2">
      <template #start>
        <Button @click="router.back()" v-tooltip="'返回'" class="mr-2">
          <template #icon>
            <v-remixicon name="riArrowLeftCircleLine" />
          </template>
        </Button>
      </template>

      <template #center>
        <span class="flex">
          <p class="pt-3 pr-2 text-xl font-serif font-semibold">{{ label }}</p>
        </span>
      </template>

      <template #end>
        <Button v-tooltip="'运行'" class="mr-2" @click="run" :loading="running">
          <template #icon>
            <v-remixicon name="riPlayLine" />
          </template>
        </Button>
        <Button v-tooltip="'保存'" @click="save">
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
      </template>
    </Toolbar>
    <Splitter id="sequence-designer" class="mb-5" stateStorage="local">
      <SplitterPanel :size="25">
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
                <v-remixicon v-bind="getIconPath(slotProps.node.icon_path)" />
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
      </SplitterPanel>
      <SplitterPanel class="h-full" :size="75">
        <div class="flex justify-around">
          <SequenceActivity :element="mainActivity.sequence" @update="update" />
        </div>
      </SplitterPanel>
    </Splitter>
  </div>
</template>
<script setup>
import Button from "primevue/button";
import Splitter from "primevue/splitter";
import SplitterPanel from "primevue/splitterpanel";
import { useRouter, onBeforeRouteLeave } from "vue-router";
import Toolbar from "primevue/toolbar";
import Tree from "primevue/tree";
import { ref, onMounted, reactive, provide } from "vue";
import SequenceActivity from "@/components/activity/SequenceActivity.vue";
import { nanoid } from "nanoid";
import {
  ParseAllPlugin,
  GetSubFlow,
  SaveSubFlow,
  RunSubFlow,
} from "@back/go/main/App";
import { EventsOn, EventsOff } from "@back/runtime/runtime";
import { getIconPath } from "@/utils/helper";
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
provide("dataChanged", { dataChanged, updateDataChanged });

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
async function run() {
  await save();
  running.value = true;
  EventsOn("log_event", (data) => {
    console.log(data);
  });
  try {
    await RunSubFlow(props.id, props.subflowId);
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
    });
  }
  running.value = false;
  EventsOff("log_event");
}
onBeforeRouteLeave(onBeforeLeave);

function onBeforeLeave() {
  if (dataChanged.value) {
    const confirm = window.confirm("有修改尚未保存，是否退出");
    if (!confirm) return false;
  }
}
</script>
<style scoped>
#sequence-designer {
  height: calc(100vh - 68px);
}
:deep(.p-splitter-panel) {
  overflow: auto;
}
</style>
