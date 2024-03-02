<template>
  <Panel
    :pt="{
      header: (options) => ({
        class: ['bg-slate-100', 'p-1'],
      }),
      root: (options) => ({
        style: {
          height: '100%',
          overflow: 'hidden',
        },
      }),
    }"
    :toggleable="false"
  >
    <template #header>
      <div class="flex align-items-center gap-4">
        <span class="mt-1">元素库</span>
      </div>
    </template>
    <template #icons>
      <div class="flex align-items-center gap-2">
        <button
          class="mr-4 hover:text-purple-500"
          @click="pickElement"
          v-tooltip.left="{
            value: '拾取元素',
            showDelay: 100,
            hideDelay: 300,
          }"
        >
          <v-remixicon name="riFocus2Line" size="24" />
        </button>
        <button
          class="mr-4 hover:text-purple-500"
          @click="deepPickElement"
          v-tooltip.left="{
            value: '深度拾取',
            showDelay: 100,
            hideDelay: 300,
          }"
        >
          <v-remixicon name="riNodeTree" size="24" />
        </button>
        <button class="hover:text-purple-500" @click="$emit('hide')">
          <span class="pi pi-chevron-down"></span>
        </button>
      </div>
    </template>
    <div class="flex relative">
      <div ref="container">
        <Tree :value="nodes" :filter="false" filterMode="lenient"></Tree>
      </div>
      <div class="w-48 h-48 absolute top-0 right-0"></div>
    </div>
    <ElementOptionDialog
      :dialogShow="dialogShow"
      :pathOption="pathOption"
      :needInit="needInit"
      @hide="dialogShow = false"
      @update="saveElement($event)"
    />
    <ElementTreeDialog
      :dialogShow="treeDialogShow"
      @hide="treeDialogShow = false"
      @confirm="selectElement($event)"
    />
  </Panel>
</template>
<script setup>
import Panel from "primevue/panel";
import Tree from "primevue/tree";
import ElementOptionDialog from "@/components/element/ElementOptionDialog.vue";
import ElementTreeDialog from "@/components/element/ElementTreeDialog.vue";
import { ref, nextTick } from "vue";
import {
  StartPickWindowsElement,
  StartCheckWindowsElement,
  RunActivity,
  SaveWindowsElement,
  GetWindowsElement,
  GetElementImage,
  GetSelectedWindowsElement,
} from "@back/go/main/App";
import { useToast } from "primevue/usetoast";
const toast = useToast();
const emit = defineEmits(["hide"]);
const props = defineProps({
  id: {
    type: String,
    default: "",
  },
});
const nodes = ref([
  {
    key: "0",
    label: "Documents",
    data: "Documents Folder",
    icon: "pi pi-fw pi-inbox",
    children: [
      {
        key: "0-0",
        label: "Work",
        data: "Work Folder",
        icon: "pi pi-fw pi-cog",
        children: [
          {
            key: "0-0-0",
            label: "Expenses.doc",
            icon: "pi pi-fw pi-file",
            data: "Expenses Document",
          },
          {
            key: "0-0-1",
            label: "Resume.doc",
            icon: "pi pi-fw pi-file",
            data: "Resume Document",
          },
        ],
      },
      {
        key: "0-1",
        label: "Home",
        data: "Home Folder",
        icon: "pi pi-fw pi-home",
        children: [
          {
            key: "0-1-0",
            label: "Invoices.txt",
            icon: "pi pi-fw pi-file",
            data: "Invoices for this month",
          },
        ],
      },
    ],
  },
]);

const container = ref();

const treeDialogShow = ref(false);
const dialogShow = ref(false);
const pathOption = ref({});
const needInit = ref(false);

async function pickElement() {
  try {
    const resp = await StartPickWindowsElement(props.id);
    const result = JSON.parse(resp);
    if (result.result === "ok") {
      needInit.value = true;
      pathOption.value = JSON.parse(result.data);
      dialogShow.value = true;
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
}

function deepPickElement() {}
</script>
<style scoped>
:deep(.p-splitter-panel) {
  overflow: hidden;
}
:deep(.p-component) {
  height: 100%;
  overflow: hidden;
}
:deep(.p-toggleable-content) {
  height: calc(100% - 40px);
}
:deep(.p-panel-content) {
  height: 100%;
  overflow: auto;
}
:deep(.p-tree) {
  padding: 0;
  border: none;
}
</style>
