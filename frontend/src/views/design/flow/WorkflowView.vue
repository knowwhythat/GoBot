<template>
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
</template>

<script setup>
import { reactive, ref, shallowRef, onBeforeMount } from "vue";
import { useRouter } from "vue-router";
import WorkflowEditor from "@/views/design/flow/WorkflowEditor.vue";
import { GetMainFlow, SaveMainFlow, DeleteSubFlow } from "@back/go/backend/App";
import { trim } from "lodash-es";
import { useToast } from "primevue/usetoast";
const toast = useToast();

const props = defineProps({
  id: {
    type: String,
    default: "",
  },
});

const workflow = ref(null);
const loaded = ref(false);

onBeforeMount(async () => {
  GetMainFlow(props.id)
    .then(({ name, data }) => {
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
        life: 1000,
      });
    });
});

const state = reactive({
  dataChanged: false,
});

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
    //打开子流程的Tab
  }
}
async function deleteNode(nodeId) {
  state.dataChanged = true;
  await DeleteSubFlow(props.id, nodeId);
}
function isChanged() {
  return state.dataChanged;
}

function getParams() {
  return [];
}

function setParams(params) {}
defineExpose({
  save,
  // delBlock,
  // cutBlock,
  // copyBlock,
  // pasteBlock,
  isChanged,
  getParams,
  setParams,
});
</script>
<style scoped>
.workflow-editor {
  height: calc(100vh - 150px);
}
</style>
