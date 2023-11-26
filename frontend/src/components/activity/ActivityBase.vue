<template>
  <Panel
    :toggleable="props.toggleable"
    :collapsed="nodeData.collapsed"
    class="activity-node min-w-[400px] m-2"
    style="cursor: grab"
    @update:collapsed="panelToggle"
    @dblclick.stop="dialogShow = true"
  >
    <template #header>
      <div class="flex align-items-center gap-2">
        <span
          :class="props.color"
          class="inline-block rounded-lg dark:text-black p-1"
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
    <Dialog
      v-model:visible="dialogShow"
      @hide="dialogShow = false"
      maximizable
      modal
      header="参数设置"
      :style="{ width: '50rem' }"
      :breakpoints="{ '1199px': '75vw', '575px': '90vw' }"
    >
      <template #header>
        <div class="flex align-items-center gap-2">
          <span
            :class="props.color"
            class="inline-block rounded-lg dark:text-black p-3"
          >
            <v-remixicon v-bind="getIconPath(props?.icon)" />
          </span>
          <InputText type="text" v-model="nodeData.label" />
        </div>
      </template>
      <div class="mb-4">
        <span class="ml-2 mr-2"> 模块名称 </span>
        <ExpressionTextInput value="" />
      </div>
      <template #footer>
        <Button
          label="取消"
          icon="pi pi-times"
          @click="dialogShow = false"
          text
        />
        <Button label="确认" icon="pi pi-check" @click="updateData" />
      </template>
    </Dialog>
  </Panel>
</template>
<script setup>
import { getIconPath } from "@/utils/helper";
import Button from "primevue/button";
import Panel from "primevue/panel";
import Dialog from "primevue/dialog";
import InputText from "primevue/inputtext";
import { reactive, ref, markRaw, inject } from "vue";
import ExpressionTextInput from "@/components/editor/ExpressionTextInput.vue";
const props = defineProps({
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
});
const emit = defineEmits(["delete", "update"]);
const { dataChanged, updateDataChanged } = inject("dataChanged");
const dialogShow = ref(false);
const nodeData = reactive({
  collapsed: props.collapsed,
  label: props.label,
});

function panelToggle(event) {
  nodeData.collapsed = event;
  emit("update", { collapsed: event });
}

function updateData() {
  console.log(dataChanged);
  dialogShow.value = false;
  updateDataChanged();
  emit("update", { ...nodeData });
}
</script>
