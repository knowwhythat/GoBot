<template>
  <Panel
    :id="props.id"
    :toggleable="false"
    :collapsed="nodeData.collapsed"
    class="rounded m-1 min-w-[400px]"
    :class="[{ 'ring-2 ring-offset-0 ring-purple-700': isSelected }]"
    style="cursor: grab; width: calc(100% - 8px)"
    @dblclick.stop="dialogShow = true"
    @click.stop="handleClick"
    :pt="{
      header: (options) => ({
        style: {
          padding: '0.25rem',
        },
        class: [
          'handle',
          'hover:bg-slate-200',
          'group/item',
          props.id === debugingId ? 'bg-red-200' : '',
        ],
      }),
    }"
  >
    <template #header>
      <div class="flex align-items-center gap-2">
        <button
          v-if="props.toggleable"
          @click.stop="panelToggle(!nodeData.collapsed)"
        >
          <span
            :class="[
              'pi',
              nodeData.collapsed ? 'pi-chevron-right' : 'pi-chevron-down',
            ]"
          ></span>
        </button>
        <p v-else class="pl-4"></p>
        <span
          :class="[props.color, props.breakpoint ? 'bg-red-600' : '']"
          class="inline-block rounded-lg dark:text-black p-1 cursor-auto hover:bg-red-600"
          @click="$emit('update', { breakpoint: !props.breakpoint })"
        >
          <v-remixicon v-bind="getIconPath(props?.icon)" />
        </span>
        <span class="w-48 truncate p-1">
          {{ label }}
        </span>
        <slot name="top"></slot>
      </div>
    </template>
    <template #icons>
      <div :class="['invisible', 'group-hover/item:visible']">
        <slot name="operate"></slot>
        <button
          v-if="props.deleteable"
          class="p-panel-header-icon p-link mr-2"
          @click.stop="$emit('delete')"
        >
          <span class="pi pi-trash"></span>
        </button>
        <button
          class="p-panel-header-icon p-link mr-2"
          @click.stop="dialogShow = true"
        >
          <span class="pi pi-pencil"></span>
        </button>
      </div>
    </template>
    <div class="flex-row content-center">
      <slot name="default"></slot>
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
import { reactive, onMounted, inject, ref, computed } from "vue";

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
const { selectedActivity } = inject("selectedActivity");

const isSelected = computed(() => selectedActivity.value.includes(props.id));

function handleClick(event) {
  if (event.ctrlKey) {
    selectedActivity.value.push(props.id);
  } else if (event.shiftKey) {
    selectedActivity.value.push(props.id);
  } else {
    selectedActivity.value.splice(0);
    selectedActivity.value.push(props.id);
  }
}

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
