<template>
  <vue-flow
    :id="props.id"
    :class="{ disabled: isDisabled }"
    @drop="onDropInEditor"
    @dragend="clearHighlightedElements"
    @dragover.prevent="onDragoverEditor"
    :default-edge-options="{
      type: 'custom',
      updatable: true,
      selectable: true,
      markerEnd: settings.arrow ? MarkerType.ArrowClosed : '',
    }"
  >
    <Background />
    <MiniMap
      v-if="minimap"
      :node-class-name="minimapNodeClassName"
      class="hidden md:block"
      pannable
      zoomable
    />
    <div
      class="flex flex-col gap-2 p-2 justify-center absolute z-10 inset-y-28 left-4 w-40 rounded-xl bg-gray-200 dark:bg-gray-800 overflow-auto"
    >
      <div
        v-for="(block, key) in BlockType"
        :key="key"
        draggable="true"
        class="transform select-none cursor-move relative p-4 rounded-lg transition group bg-white"
        @dragstart="$event.dataTransfer.setData('block', JSON.stringify(block))"
      >
        <div
          class="flex items-center absolute right-2 invisible group-hover:visible top-2 text-gray-600 dark:text-gray-300"
        >
          <a :title="block.description" target="_blank" rel="noopener">
            <v-remixicon name="riInformationLine" size="18" />
          </a>
        </div>
        <v-remixicon
          :path="getIconPath(block.icon)"
          :name="block.icon"
          size="24"
          class="mb-2"
        />
        <p class="leading-tight text-overflow capitalize">
          {{ block.name }}
        </p>
      </div>
    </div>
    <div
      v-if="editorControls"
      class="flex items-center absolute w-full p-4 left-0 bottom-0 z-10 md:pr-60"
    >
      <slot name="controls-prepend" />
      <editor-search-blocks :editor="editor" />
      <div class="flex-grow pointer-events-none" />
      <slot name="controls-append" />
      <div class="rounded-lg bg-white dark:bg-gray-800 inline-block">
        <Button
          v-tooltip="'重置'"
          class="p-2 m-1 rounded-lg"
          @click="editor.fitView()"
        >
          <v-remixicon name="riFullscreenLine" />
        </Button>
        <Button
          v-tooltip="'缩小'"
          class="p-2 m-1 rounded-lg relative z-10"
          @click="editor.zoomOut()"
        >
          <v-remixicon name="riSubtractLine" />
        </Button>
        <hr class="h-6 border-r inline-block" />
        <Button
          v-tooltip.top="'放大'"
          class="p-2 m-1 rounded-lg"
          @click="editor.zoomIn()"
        >
          <v-remixicon name="riAddLine" />
        </Button>
      </div>
    </div>
    <template v-for="(node, name) in nodeTypes" :key="name" #[name]="nodeProps">
      <component
        :is="node"
        v-bind="{ ...nodeProps, editor: editor }"
        @delete="deleteBlock(nodeProps)"
        @settings="initEditBlockSettings"
        @edit="editBlock(nodeProps)"
        @update="updateBlockData(nodeProps.id, $event)"
      />
    </template>
    <template #edge-custom="edgeProps">
      <EditorCustomEdge v-bind="edgeProps" />
    </template>
    <Dialog
      v-model:visible="blockSettingsState.show"
      @hide="clearBlockSettings"
      maximizable
      modal
      header="模块设置"
      :style="{ width: '50rem' }"
      :breakpoints="{ '1199px': '75vw', '575px': '90vw' }"
    >
      <editor-block-setting
        :data="blockSettingsState?.data"
        @change="updateBlockSettingsData"
      />
    </Dialog>
  </vue-flow>
</template>
<script setup>
import { computed, onBeforeUnmount, onMounted, reactive, watch } from "vue";
import {
  getConnectedEdges,
  MarkerType,
  useVueFlow,
  VueFlow,
} from "@vue-flow/core";
import { Background, MiniMap } from "@vue-flow/additional-components";
import cloneDeep from "lodash.clonedeep";
import { customAlphabet } from "nanoid";
import { parseJSON } from "@/utils/helper";
import { useAppStore } from "@/stores/app";
import { BlockType } from "@/utils/shared.js";
import EditorCustomEdge from "./editor/EditorCustomEdge.vue";
import EditorSearchBlocks from "./editor/EditorSearchBlocks.vue";
import EditorBlockSetting from "./editor/EditorBlockSetting.vue";
import { useToast } from "primevue/usetoast";
import Button from "primevue/button";
import Dialog from "primevue/dialog";

const toast = useToast();
const nanoid = customAlphabet("1234567890abcdefghijklmnopqrstuvwxyz", 7);

const props = defineProps({
  id: {
    type: String,
    default: "editor",
  },
  data: {
    type: Object,
    default: () => ({
      x: 0,
      y: 0,
      zoom: 0,
      nodes: [],
      edges: [],
      params: [],
    }),
  },
  options: {
    type: Object,
    default: () => ({}),
  },
  editorControls: {
    type: Boolean,
    default: true,
  },
  minimap: {
    type: Boolean,
    default: true,
  },
  disabled: Boolean,
});

const emit = defineEmits([
  "edit",
  "init",
  "update:node",
  "delete:node",
  "update:settings",
  "update:flow",
]);

const blockComponents = import.meta.glob("@/components/block/*.vue", {
  eager: true,
});
const nodeTypes = {};
for (let each in blockComponents) {
  const name = blockComponents[each].default.__name;
  nodeTypes[`node-${name}`] = blockComponents[each].default;
}

const getPosition = (position) => (Array.isArray(position) ? position : [0, 0]);
const setMinValue = (num, min) => (num < min ? min : num);

const appStore = useAppStore();
const isMac = appStore.settings.operateSystem.indexOf("Mac") !== -1;

const editor = useVueFlow({
  id: props.id,
  edgeUpdaterRadius: 20,
  deleteKeyCode: "Delete",
  elevateEdgesOnSelect: true,
  defaultZoom: props.data?.zoom ?? 1,
  minZoom: setMinValue(+appStore.settings.editor.minZoom || 0.5, 0.1),
  maxZoom: setMinValue(
    +appStore.settings.editor.maxZoom || 1.2,
    +appStore.settings.editor.minZoom + 0.1,
  ),
  multiSelectionKeyCode: isMac ? "Meta" : "Control",
  defaultPosition: getPosition(props.data?.position),
  ...props.options,
});

editor.onConnect((params) => {
  if (params.source === params.target) {
    toast.add({
      severity: "warn",
      summary: "警告",
      detail: "起点和终点不能为同一个",
      life: 3000,
    });
    return;
  }
  if (
    params.sourceHandle.includes("input") ||
    params.targetHandle.includes("output")
  ) {
    toast.add({
      severity: "warn",
      summary: "警告",
      detail: "两个端点不允许都为输入或者输出",
      life: 3000,
    });
    return;
  }
  const oldEdge = editor.getEdges.value.find(
    (edge) =>
      edge.source === params.source &&
      edge.sourceHandle === params.sourceHandle,
  );
  if (oldEdge) {
    toast.add({
      severity: "warn",
      summary: "警告",
      detail: "只允许有一个输出方向",
      life: 3000,
    });
    return;
  }
  if (params.sourceHandle.includes("fallback")) {
    params.label = "异常";
    params.animated = true;
  }
  params.class = `source-${params.sourceHandle} target-${params.targetHandle}`;
  params.updatable = true;
  editor.addEdges([params]);
  emit("update:flow");
});

editor.onEdgeUpdate(({ edge, connection }) => {
  if (connection.source === connection.target) {
    toast.add({
      severity: "warn",
      summary: "警告",
      detail: "起点和终点不能为同一个",
      life: 3000,
    });
    return;
  }
  if (
    connection.sourceHandle.includes("input") ||
    connection.targetHandle.includes("output")
  ) {
    toast.add({
      severity: "warn",
      summary: "警告",
      detail: "两个端点不允许都为输入或者输出",
      life: 3000,
    });
    return;
  }
  if (connection.sourceHandle.includes("fallback")) {
    connection.label = "异常";
    connection.animated = true;
  }
  Object.assign(edge, connection);
  emit("update:flow");
});

function getIconPath(path) {
  if (path && path.startsWith("path")) {
    const { 1: iconPath } = path.split(":");
    return iconPath;
  }

  return "";
}

function toggleHighlightElement({ target, elClass, classes }) {
  const targetEl = target.closest(elClass);

  if (targetEl) {
    targetEl.classList.add(classes);
  } else {
    const elements = document.querySelectorAll(`.${classes}`);
    elements.forEach((element) => {
      element.classList.remove(classes);
    });
  }
}

function onDragoverEditor({ target }) {
  toggleHighlightElement({
    target,
    elClass: ".vue-flow__handle.source",
    classes: "dropable-area__handle",
  });

  if (!target.closest(".vue-flow__handle")) {
    toggleHighlightElement({
      target,
      elClass: ".vue-flow__node:not(.vue-flow__node-BlockGroup)",
      classes: "dropable-area__node",
    });
  }
}

function onDropInEditor({ dataTransfer, clientX, clientY, target }) {
  const block = parseJSON(dataTransfer.getData("block"), null);
  clearHighlightedElements();

  const isTriggerExists =
    block.blockId === "Start" &&
    editor.getNodes.value.some((node) => node.data.nodeType === "Start");
  if (isTriggerExists) {
    toast.add({
      severity: "warn",
      summary: "警告",
      detail: "只允许有一个开始节点",
      life: 3000,
    });
    return;
  }
  const position = editor.project({ x: clientX - 60, y: clientY - 80 });
  const nodeId = nanoid();
  const newNode = {
    position,
    data: block.data,
    type: block.component,
    id: nodeId,
  };
  editor.addNodes([newNode]);
  emit("update:flow");
}

function clearHighlightedElements() {
  const elements = document.querySelectorAll(
    ".dropable-area__node, .dropable-area__handle",
  );
  elements.forEach((element) => {
    element.classList.remove("dropable-area__node");
    element.classList.remove("dropable-area__handle");
  });
}

const settings = appStore.settings.editor;
const isDisabled = computed(() => props.options.disabled ?? props.disabled);

const blockSettingsState = reactive({
  show: false,
  data: {},
});

function initEditBlockSettings({ blockId, data }) {
  blockSettingsState.data = {
    blockId,
    data: data,
  };
  blockSettingsState.show = true;
}

function clearBlockSettings() {
  Object.assign(blockSettingsState, {
    data: null,
    show: false,
  });
}

function minimapNodeClassName({ data }) {
  return BlockType[data.nodeType].color;
}

function updateBlockData(nodeId, data = {}) {
  if (isDisabled.value) return;

  const node = editor.findNode(nodeId);
  node.data = { ...node.data, ...data };

  emit("update:node", node);
}

function updateBlockSettingsData(newSettings) {
  if (isDisabled.value) return;
  const nodeId = blockSettingsState.data.blockId;
  const node = editor.findNode(nodeId);
  node.data = { ...node.data, ...newSettings };

  emit("update:settings", {
    settings: newSettings,
    blockId: blockSettingsState.data.blockId,
  });
}

function editBlock(node) {
  if (isDisabled.value) return;
  emit("edit", node);
}

function deleteBlock(node) {
  if (isDisabled.value) return;
  editor.removeNodes([node.id]);
  emit("delete:node", node);
}

function onMousedown(event) {
  if (isDisabled.value && event.shiftKey) {
    event.stopPropagation();
    event.preventDefault();
  }
}

function applyFlowData() {
  if (settings.snapToGrid) {
    editor.snapToGrid.value = true;
    editor.snapGrid.value = Object.values(settings.snapGrid);
  }
  editor.setNodes(
    props.data?.nodes?.map((node) => ({ ...node, events: {} })) || [],
  );
  editor.setEdges(props.data?.edges || []);
  editor.setViewport({
    x: props.data?.x || 0,
    y: props.data?.y || 0,
    zoom: props.data?.zoom || 1,
  });
}

watch(
  () => props.disabled,
  (value) => {
    const keys = [
      "nodesDraggable",
      "edgesUpdatable",
      "nodesConnectable",
      "elementsSelectable",
    ];

    keys.forEach((key) => {
      editor[key].value = !value;
    });
  },
  { immediate: true },
);
watch(editor.getSelectedNodes, (nodes, _, cleanup) => {
  const connectedEdges = getConnectedEdges(nodes, editor.getEdges.value);

  connectedEdges.forEach((edge) => {
    edge.class = "connected-edges";
  });

  cleanup(() => {
    connectedEdges.forEach((edge) => {
      edge.class = undefined;
    });
  });
});

onMounted(() => {
  applyFlowData();
  window.addEventListener("mousedown", onMousedown, true);
  emit("init", editor);
});
onBeforeUnmount(() => {
  window.removeEventListener("mousedown", onMousedown, true);
});
</script>
<style>
@import "@vue-flow/core/dist/style.css";
@import "@vue-flow/core/dist/theme-default.css";

.control-button {
  @apply p-2 rounded-lg bg-white dark:bg-gray-800 transition-colors;
}
</style>
