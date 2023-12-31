<template>
  <Panel
    :id="props.id"
    :toggleable="props.toggleable"
    :collapsed="nodeData.collapsed"
    class="activity-node min-w-[400px] m-2"
    style="cursor: grab"
    @update:collapsed="panelToggle"
    @dblclick.stop="dialogShow = true"
    :pt="{
      header: (options) => ({
        style: {
          padding: '0.75rem',
        },
        class: ['handle', props.id === debugingId ? 'bg-red-200' : ''],
      }),
    }"
  >
    <template #header>
      <div class="flex align-items-center gap-2">
        <span
          :class="[props.color, props.breakpoint ? 'bg-red-600' : '']"
          class="inline-block rounded-lg dark:text-black p-1 hover:bg-red-600 cursor-auto"
          @click="$emit('update', { breakpoint: !props.breakpoint })"
        >
          <v-remixicon v-bind="getIconPath(props?.icon)" />
        </span>
        <span class="max-w-xs truncate p-1">
          {{ label }}
        </span>
      </div>
    </template>
    <template #icons>
      <slot name="operate"></slot>
      <button
        v-if="props.deleteable"
        class="p-panel-header-icon p-link mr-2"
        @click.stop="$emit('delete')"
        v-tooltip="{ value: '删除', showDelay: 100, hideDelay: 300 }"
      >
        <span class="pi pi-trash"></span>
      </button>
      <button
        class="p-panel-header-icon p-link mr-2"
        @click.stop="dialogShow = true"
        v-tooltip="{ value: '编辑', showDelay: 100, hideDelay: 300 }"
      >
        <span class="pi pi-pencil"></span>
      </button>
    </template>
    <div class="flex-row content-center">
      <slot></slot>
    </div>
    <ParamEditorDialog
      :runnable="props.runnable"
      :icon="props.icon"
      :color="props.color"
      :dialogShow="dialogShow"
      :parameter_define="props.parameter_define"
      :nodeData="nodeData"
      @hide="dialogShow = false"
      @run="$emit('run')"
      @update="updateData($event)"
    ></ParamEditorDialog>
  </Panel>
</template>
<script setup>
import { getIconPath } from "@/utils/helper";
import Button from "primevue/button";
import Panel from "primevue/panel";
import ParamEditorDialog from "@/components/editor/ParamEditorDialog.vue";
import { reactive, onMounted, inject, ref } from "vue";

const props = defineProps({
  id: {
    type: String,
    default: "",
  },
  label: {
    type: String,
    default: "",
  },
  icon: {
    type: String,
    default: "",
  },
  color: {
    type: String,
    default:
      "bg-green-200 dark:bg-green-300 fill-green-200 dark:fill-green-300",
  },
  collapsed: {
    type: Boolean,
    default: false,
  },
  toggleable: {
    type: Boolean,
    default: true,
  },
  deleteable: {
    type: Boolean,
    default: true,
  },
  runnable: {
    type: Boolean,
    default: false,
  },
  breakpoint: {
    type: Boolean,
    default: false,
  },
  parameter_define: {
    type: Object,
    default: () => ({}),
  },
  parameter: {
    type: Object,
    default: () => ({}),
  },
});

const emit = defineEmits(["delete", "update", "run"]);
const { dataChanged, updateDataChanged } = inject("dataChanged");
const { debugingId } = inject("debugingId");

const nodeData = reactive({
  collapsed: false,
  label: "",
  parameter: {},
});

onMounted(() => {
  nodeData.collapsed = props.collapsed;
  nodeData.label = props.label;
  nodeData.parameter = props.parameter;
});

const dialogShow = ref(false);

function panelToggle(event) {
  nodeData.collapsed = event;
  emit("update", { collapsed: event });
}

function updateData(data) {
  dialogShow.value = false;
  updateDataChanged();
  nodeData.label = data.label;
  nodeData.parameter = data.parameter;
  emit("update", { ...data });
}
</script>
