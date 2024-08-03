<template>
  <div class="code-editor">
    <MonacoEditor :code="codeData.code" @update="updateCode($event)" />
  </div>
</template>
<script setup>
import MonacoEditor from "@/components/MonacoEditor.vue";
import { onBeforeMount, ref, watch } from "vue";
import { GetSubModule, SaveSubModule } from "@back/go/backend/App";
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
const emit = defineEmits(["mounted", "dataChanged"]);

const codeData = ref({
  params: [],
  code: "",
});

onBeforeMount(async () => {
  const result = await GetSubModule(props.id, props.subflowId);
  if (result) {
    codeData.value = JSON.parse(result);
  } else {
    codeData.value = {
      params: [],
      code: `from package import variables as glv #全局变量,例如glv['test]
from robot_base import log_util
import robot_basic

def main(args):
    #输入参数使用示例
    # if args is None:
    #    输入参数1 = ""
    #else:
    #    输入参数1 = args.get("输入参数1", "")`,
    };
  }
  emit("mounted");
});

const dataChanged = ref(false);

function updateCode(code) {
  if (codeData.value.code !== code) {
    codeData.value.code = code;
    dataChanged.value = true;
    emit("dataChanged", dataChanged.value);
  }
}

async function save() {
  try {
    await SaveSubModule(
      props.id,
      props.subflowId,
      JSON.stringify(codeData.value),
    );
    dataChanged.value = false;
    emit("dataChanged", false);
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

function isChanged() {
  return dataChanged.value;
}

function getParams() {
  return codeData.value.params;
}

function setParams(params) {
  codeData.value.params = params;
}

defineExpose({
  save,
  isChanged,
  getParams,
  setParams,
});
</script>

<style scoped>
.code-editor {
  height: calc(100vh - 144px);
}
</style>
