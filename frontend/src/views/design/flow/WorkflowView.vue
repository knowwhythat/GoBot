<template>
  <WorkflowEditor
    v-if="loaded"
    :id="props.id"
    :data="workflow"
    class="workflow-editor"
    @init="onEditorInit"
    @edit="editBlock"
    @update:flow="dataChanged = true"
    @update:node="dataChanged = true"
    @update:settings="dataChanged = true"
    @delete:node="deleteNode"
  />
</template>

<script setup>
import { ref, shallowRef, onBeforeMount } from "vue";
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

const emit = defineEmits(["edit:block", "delete:block"]);

const workflow = ref({});
const loaded = ref(false);

onBeforeMount(async () => {
  try {
    const data = await GetMainFlow(props.id);
    if (trim(data).length > 0) {
      workflow.value = JSON.parse(data);
    }
    loaded.value = true;
  } catch (err) {
    loaded.value = true;
    toast.add({
      severity: "error",
      summary: "错误",
      detail: err,
      life: 1000,
    });
  }
});

const dataChanged = ref(false);

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

  const triggerBlock = flow.nodes.find(
    (node) => node.data.nodeType === "Start",
  );
  if (!triggerBlock) {
    toast.add({
      severity: "warn",
      summary: "警告",
      detail: "当前流程图中不存在开始节点",
      life: 3000,
    });
    return;
  }
  flow.params = workflow.value.params;
  try {
    await SaveMainFlow(props.id, JSON.stringify(flow));
    dataChanged.value = false;
  } catch (err) {
    toast.add({
      severity: "error",
      summary: "错误",
      detail: err,
      life: 1000,
    });
  }
}
async function editBlock(node) {
  if (node.data.nodeType === "SubFlow") {
    try {
      await save();
      emit("edit:block", node);
    } catch (ex) {
      toast.add({ severity: "error", summary: "错误", detail: ex, life: 3000 });
    }
  }
}

function deleteNode(node) {
  dataChanged.value = true;
  if (node.data.nodeType === "SubFlow") {
    emit("delete:block", node.id);
  }
}

function isChanged() {
  return dataChanged.value;
}

function getParams() {
  return workflow.value.params ? workflow.value.params : [];
}

function setParams(params) {
  workflow.value.params = params;
}
defineExpose({
  save,
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
