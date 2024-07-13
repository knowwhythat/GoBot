<template>
  <Panel
    header="元素库"
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
          <v-remixicon name="riFocus2Line" size="16" />
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
          <v-remixicon name="riNodeTree" size="16" />
        </button>
        <button class="hover:text-purple-500" @click="$emit('hide')">
          <span class="pi pi-chevron-down"></span>
        </button>
      </div>
    </template>
    <div class="flex relative h-full">
      <Tree
        :value="windowsElement"
        :filter="false"
        filterMode="lenient"
        selectionMode="single"
        v-model:expandedKeys="expandedKeys"
        v-model:selectionKeys="selectedKey"
        @node-select="nodeSelect"
        :pt="{
          root: (options) => ({
            class: ['flex-1'],
          }),
          content: (options) => ({
            oncontextmenu: onContextMenuClick(options),
          }),
        }"
      >
        <template #default="slotProps">
          <div>{{ slotProps.node.label }}</div>
        </template>
      </Tree>
      <ContextMenu ref="menu" :model="items" />
      <div class="w-64 h-full">
        <Image
          v-if="imagePath"
          :src="imagePath"
          alt="Image"
          preview
          class="flex"
          image-class="flex justify-center items-center object-scale-down"
        />
      </div>
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
import Image from "primevue/image";
import ContextMenu from "primevue/contextmenu";
import Panel from "primevue/panel";
import Tree from "primevue/tree";
import ElementOptionDialog from "@/components/element/ElementOptionDialog.vue";
import ElementTreeDialog from "@/components/element/ElementTreeDialog.vue";
import { ref, toRaw } from "vue";
import {
  GetElementImage,
  GetSelectedWindowsElement,
  RemoveWindowsElement,
  SaveWindowsElement,
  StartPickWindowsElement,
} from "@back/go/backend/App";
import {
  EventsOnce,
  WindowMaximise,
  WindowMinimise,
} from "@back/runtime/runtime.js";
import { useToast } from "primevue/usetoast";
import { useConfirm } from "primevue/useconfirm";

const toast = useToast();
const confirm = useConfirm();
const emit = defineEmits(["hide"]);
const props = defineProps({
  id: {
    type: String,
    default: "",
  },
  windowsElement: {
    type: Array,
    default: () => [],
  },
});
const selectedKey = ref(null);
const expandedKeys = ref({});

const treeDialogShow = ref(false);
const dialogShow = ref(false);
const pathOption = ref({});
const needInit = ref(false);

async function pickElement() {
  try {
    await StartPickWindowsElement();
    await WindowMinimise();
    EventsOnce("windowsEvent", async (resp) => {
      await WindowMaximise();
      const result = JSON.parse(resp);
      console.log(result);
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
    });
  } catch (err) {
    toast.add({
      severity: "error",
      summary: "拾取失败",
      detail: err,
      life: 3000,
    });
  }
}

function deepPickElement() {
  treeDialogShow.value = true;
}

async function selectElement(id) {
  try {
    await GetSelectedWindowsElement(id);
    await WindowMinimise();
    EventsOnce("windowsEvent", async (resp) => {
      await WindowMaximise();
      const result = JSON.parse(resp);
      if (result.result === "ok") {
        needInit.value = true;
        pathOption.value = JSON.parse(result.data);
        treeDialogShow.value = false;
        dialogShow.value = true;
      } else {
        toast.add({
          severity: "error",
          summary: "拾取失败",
          detail: result.reason,
          life: 3000,
        });
      }
    });
  } catch (err) {
    toast.add({
      severity: "error",
      summary: "拾取失败",
      detail: err,
      life: 3000,
    });
  }
}

async function saveElement(element) {
  const image = element.image;
  delete element.image;

  await SaveWindowsElement(
    props.id,
    element.id,
    image,
    JSON.stringify(element),
  );
  dialogShow.value = false;
}

const imagePath = ref(null);

async function nodeSelect(node) {
  if (!node["children"]) {
    const image = await GetElementImage(props.id, node.key);
    imagePath.value = "data:image/png;base64," + image;
  } else {
    imagePath.value = null;
  }
}

const selectedNode = ref(null);
const menu = ref();
const items = ref([
  {
    label: "删除",
    icon: "pi pi-times-circle",
    command: () => {
      confirm.require({
        header: "提示",
        message: "确定要删除这条记录吗?",
        icon: "pi pi-info-circle",
        rejectClass: "p-button-secondary p-button-outlined p-button-sm",
        acceptClass: "p-button-danger p-button-sm",
        rejectLabel: "取消",
        acceptLabel: "删除",
        accept: () => {
          RemoveWindowsElement(props.id, selectedNode.value.key).catch(
            (error) => {
              console.error(error);
            },
          );
        },
        reject: () => {},
      });
    },
  },
  {
    label: "编辑",
    icon: "pi pi-file-edit",
    command: () => {
      needInit.value = false;
      console.log(selectedNode.value);
      pathOption.value = toRaw(selectedNode.value.data);
      dialogShow.value = true;
    },
  },
]);

function onContextMenuClick(options) {
  return async function (event) {
    if (!options.props.node["children"]) {
      selectedNode.value = options.props.node;
      menu.value.show(event);
    }
  };
}
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

:deep(.p-tree-wrapper) {
  height: 100%;
}

:deep(.p-panel-content) {
  padding: 0.75rem;
}
</style>
