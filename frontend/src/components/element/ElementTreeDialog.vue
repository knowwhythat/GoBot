<template>
  <Dialog
    :visible="visible"
    v-on:update:visible="$emit('hide')"
    maximizable
    modal
    header="元素结构"
    @show="show"
    :style="{ width: '80rem' }"
    :breakpoints="{ '1199px': '75vw', '575px': '90vw' }"
    :pt="{
      header: (options) => ({
        style: {
          padding: '0.75rem',
        },
      }),
    }"
  >
    <Splitter>
      <SplitterPanel :size="60">
        <Tree
          :value="nodes"
          selectionMode="single"
          @node-expand="onNodeExpand"
          @node-select="onNodeSelect"
          :loading="loading"
          class="w-full md:w-30rem"
        ></Tree>
      </SplitterPanel>
      <SplitterPanel :size="40">
        <DataTable :value="seletedNode.attrs" showGridlines>
          <Column field="key" header="属性名称"> </Column>
          <Column field="value" header="属性值"></Column>
        </DataTable>
      </SplitterPanel>
    </Splitter>
    <template #footer>
      <Button label="取消" icon="pi pi-times" @click="$emit('hide')" text />
      <Button label="确认" icon="pi pi-check" @click="updateData" />
    </template>
  </Dialog>
</template>

<script setup>
import Dialog from "primevue/dialog";
import Splitter from "primevue/splitter";
import SplitterPanel from "primevue/splitterpanel";
import Tree from "primevue/tree";
import DataTable from "primevue/datatable";
import Column from "primevue/column";
import Button from "primevue/button";
import { ref, watch } from "vue";
import { GetWindowsElementList } from "@back/go/main/App";
import { useToast } from "primevue/usetoast";
const toast = useToast();

const emit = defineEmits(["hide", "confirm"]);
const props = defineProps({
  dialogShow: {
    type: Boolean,
    default: false,
  },
});

const visible = ref(false);
watch(
  () => props.dialogShow,
  (value, oldValue) => {
    visible.value = value;
  },
  { immediate: true }
);

const loading = ref(false);
const nodes = ref(null);
const seletedNode = ref({
  id: "",
  attrs: [],
});
const onNodeSelect = (node) => {
  let attrs = [];
  for (let [key, value] of Object.entries(node.data)) {
    if (key != "id") {
      attrs.push({ key: key, value: value });
    }
  }
  seletedNode.id = node.key;
  seletedNode.value.attrs = attrs;
};
const onNodeExpand = async (node) => {
  if (!node.leaf && node.children.length == 0) {
    loading.value = true;
    try {
      const resp = await GetWindowsElementList(node.key);
      const result = JSON.parse(resp);
      if (result.result === "ok") {
        const elements = JSON.parse(result.data);
        if (elements && elements.length > 0) {
          elements.forEach((element) => {
            node.children.push({
              key: element.id,
              label: element.localizeControlType + "-" + element.name,
              children: [],
              data: element,
              leaf: false,
            });
          });
        } else {
          node.leaf = true;
        }
      } else {
        toast.add({
          severity: "error",
          summary: "拾取失败",
          detail: result.reason,
          life: 3000,
        });
      }
    } catch (err) {
      toast.add({
        severity: "error",
        summary: "拾取失败",
        detail: err,
        life: 3000,
      });
    }
    loading.value = false;
  }
};
async function show() {
  loading.value = true;
  try {
    const resp = await GetWindowsElementList("");
    const result = JSON.parse(resp);
    if (result.result === "ok") {
      const elements = JSON.parse(result.data);
      nodes.value = elements.map((element) => {
        return {
          key: element.id,
          label: element.localizeControlType + "-" + element.name,
          children: [],
          data: element,
          leaf: false,
        };
      });
    } else {
      toast.add({
        severity: "error",
        summary: "拾取失败",
        detail: result.reason,
        life: 3000,
      });
    }
  } catch (err) {
    toast.add({
      severity: "error",
      summary: "拾取失败",
      detail: err,
      life: 3000,
    });
  }
  seletedNode.value = {
    id: "",
    attrs: [],
  };
  loading.value = false;
}

function updateData() {
  if (seletedNode.id) {
    emit("confirm", seletedNode.id);
  } else {
    toast.add({
      severity: "warn",
      summary: "请先选择左侧树元素",
      life: 3000,
    });
  }
}
</script>
<style scoped>
:deep(.p-treenode-label) {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
:deep(.p-tree-container) {
  max-height: 650px;
}
</style>
