<template>
  <div class="visual-flow-tab flex justify-around">
    <SequenceActivity :element="mainActivity.sequence" @update="update" />
  </div>
</template>

<script setup>
import SequenceActivity from "@/components/activity/SequenceActivity.vue";
import { GetSubFlow, SaveSubFlow } from "@back/go/backend/App";
import { onMounted, provide, ref, watch } from "vue";
import { useToast } from "primevue/usetoast";
import {
  copySelected,
  cutSelected,
  deleteSelected,
  innerPaste,
} from "@/utils/activityUtils.js";
import { useEditorStore } from "@/stores/editor";
import { nanoid } from "nanoid";

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

const emit = defineEmits(["mounted", "dataChanged"]);
const selectedActivity = ref([]);
provide("contextVariable", { contextVariable: [] });
provide("selectedActivity", { selectedActivity });

const dataChanged = ref(false);
provide("dataChanged", { dataChanged, updateDataChanged });

function updateDataChanged() {
  dataChanged.value = true;
}

watch(
  () => dataChanged.value,
  (now, old) => {
    emit("dataChanged", dataChanged.value);
  },
  { deep: true },
);

const mainActivity = ref({
  params: [],
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

onMounted(async () => {
  try {
    const result = await GetSubFlow(props.id, props.subflowId);
    if (result) {
      mainActivity.value = JSON.parse(result);
    } else {
      mainActivity.value.sequence.id = props.subflowId;
      mainActivity.value.sequence.label = props.label ? props.label : "子流程";
    }
    emit("mounted");
  } catch (error) {
    toast.add({
      severity: "error",
      summary: "加载流程异常",
      detail: error,
      life: 3000,
    });
  }
});

function update({ children }) {
  mainActivity.value.sequence.children = children;
  if (!dataChanged.value) {
    dataChanged.value = true;
  }
}

async function save() {
  try {
    await SaveSubFlow(
      props.id,
      props.subflowId,
      JSON.stringify(mainActivity.value),
    );
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

function delBlock() {
  if (selectedActivity.value.length === 0) {
    return;
  }
  selectedActivity.value.forEach((id) => {
    deleteSelected(mainActivity.value.sequence.children, id);
  });
  dataChanged.value = true;
  selectedActivity.value = [];
}

const editorStore = useEditorStore();

function cutBlock() {
  if (selectedActivity.value.length === 0) {
    return;
  }
  let cutBlock = [];
  selectedActivity.value.forEach((id) => {
    cutSelected(mainActivity.value.sequence.children, id, cutBlock);
  });
  dataChanged.value = true;
  selectedActivity.value = [];
  editorStore.addToPasteBlocks(cutBlock);
}

function copyBlock() {
  if (selectedActivity.value.length === 0) {
    return;
  }
  let copyBlock = [];
  selectedActivity.value.forEach((id) => {
    copySelected(mainActivity.value.sequence.children, id, copyBlock);
  });
  editorStore.addToPasteBlocks(copyBlock);
}

function pasteBlock() {
  if (editorStore.copiedBlocks.length > 0) {
    dataChanged.value = true;
    if (selectedActivity.value.length > 0) {
      const id = selectedActivity.value[selectedActivity.value.length - 1];
      innerPaste(
        mainActivity.value.sequence.children,
        id,
        editorStore.copiedBlocks,
      );
    } else {
      editorStore.copiedBlocks.forEach((block) => {
        mainActivity.value.sequence.children.push({ ...block, id: nanoid(16) });
      });
    }
  }
  editorStore.clearCopiedBlocks();
}

function isChanged() {
  return dataChanged.value;
}

function getParams() {
  return mainActivity.value.params;
}

function setParams(params) {
  mainActivity.value.params = params;
}

defineExpose({
  save,
  delBlock,
  cutBlock,
  copyBlock,
  pasteBlock,
  isChanged,
  getParams,
  setParams,
});
</script>
