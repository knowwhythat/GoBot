<template>
  <div class="flex flex-col">
    <Toolbar class="p-2">
      <template #start>
        <Button
          @click="router.back()"
          v-tooltip="'返回'"
          class="mr-2 px-3 py-2"
        >
          <template #icon>
            <v-remixicon name="riArrowLeftCircleLine" />
          </template>
        </Button>
        <Button v-tooltip="'运行'" class="mr-2 px-3 py-2">
          <template #icon>
            <v-remixicon name="riPlayLine" />
          </template>
        </Button>
      </template>

      <template #center>
        <span class="flex">
          <p class="pt-3 pr-2 text-xl font-serif font-semibold">
            {{ projectName }}
          </p>
          <Button v-tooltip="'编辑'" class="px-3 py-2">
            <template #icon>
              <v-remixicon name="riEditBoxLine" />
            </template>
          </Button>
        </span>
      </template>

      <template #end>
        <Button v-tooltip="'发布'" class="mr-2 px-3 py-2">
          <template #icon>
            <v-remixicon name="riSendPlaneLine" />
          </template>
        </Button>
        <Button v-tooltip="'保存'" @click="save" class="px-3 py-2">
          <template #icon>
            <span>
              <span
                v-if="state.dataChanged"
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
    <WorkflowEditor
      v-if="loaded"
      :id="props.id"
      :data="workflow"
      class="workflow-editor"
      @init="onEditorInit"
      @edit="editBlock"
      @update:flow="state.dataChanged = true"
      @update:node="state.dataChanged = true"
      @update:settings="state.dataChanged = true"
      @delete:node="deleteNode"
    />
  </div>
</template>

<script setup>
import { reactive, ref, shallowRef, onBeforeMount } from "vue";
import { useRouter, onBeforeRouteLeave } from "vue-router";
import Toolbar from "primevue/toolbar";
import Button from "primevue/button";
import WorkflowEditor from "@/views/design/flow/WorkflowEditor.vue";
import { useToast } from "primevue/usetoast";
import { GetMainFlow, SaveMainFlow, DeleteSubFlow } from "@back/go/main/App";
import { WindowMaximise } from "@back/runtime/runtime";
import { trim } from "lodash-es";
const toast = useToast();

const props = defineProps({
  id: {
    type: String,
    default: "",
  },
});

const workflow = ref(null);
const loaded = ref(false);
const projectName = ref("");

onBeforeMount(async () => {
  await WindowMaximise();
  GetMainFlow(props.id)
    .then(({ name, data }) => {
      projectName.value = name;
      if (trim(data).length > 0) {
        workflow.value = JSON.parse(data);
      }
      loaded.value = true;
    })
    .catch((error) => {
      loaded.value = true;
      console.error(error);
      toast.add({
        severity: "error",
        summary: "错误",
        detail: error,
        life: 3000,
      });
    });
});
const state = reactive({
  dataChanged: false,
});

onBeforeRouteLeave(onBeforeLeave);

function onBeforeLeave() {
  if (state.dataChanged) {
    const confirm = window.confirm("有修改尚未保存，是否退出");
    if (!confirm) return false;
  }
}

const router = useRouter();
const editor = shallowRef(null);

function onEditorInit(instance) {
  editor.value = instance;
}

async function save() {
  const flow = editor.value.toObject();
  flow.edges = flow.edges.map((edge) => {
    delete edge.sourceNode;
    delete edge.targetNode;

    return edge;
  });

  const triggerBlock = flow.nodes.find((node) => node.label === "Start");
  if (!triggerBlock) {
    toast.add({
      severity: "warn",
      summary: "警告",
      detail: "当前流程图中不存在开始节点",
      life: 3000,
    });
    return;
  }
  await SaveMainFlow(props.id, JSON.stringify(flow));
  state.dataChanged = false;
}
async function editBlock(data) {
  if (data.id === "SubFlow") {
    try {
      await save();
    } catch (ex) {
      toast.add({ severity: "error", summary: "错误", detail: ex, life: 3000 });
    }
    router.push(
      `design/sequence?id=${props.id}&subflowId=${data.blockId}&label=${data.data.label}`
    );
  }
}
async function deleteNode(nodeId) {
  state.dataChanged = true;
  await DeleteSubFlow(props.id, nodeId);
}
</script>
<style scoped>
.workflow-editor {
  height: calc(100vh - 70px);
}
</style>
