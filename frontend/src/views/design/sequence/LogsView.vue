<template>
  <Panel
    header="日志"
    :pt="{
      header: (options) => ({
        class: ['bg-slate-100', 'p-2'],
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
        <button class="mr-4 hover:text-purple-500" @click="$emit('clear')">
          <span class="pi pi-delete-left"></span>
        </button>
        <button class="mr-4 hover:text-purple-500" @click="filter = !filter">
          <span class="pi pi-search"></span>
        </button>
        <button class="hover:text-purple-500" @click="$emit('hide')">
          <span class="pi pi-chevron-down"></span>
        </button>
      </div>
    </template>
    <div ref="container">
      <Listbox
        :options="logs"
        :filter="filter"
        emptyMessage="没有日志"
        listStyle="height: calc(100% - 30px);"
      />
    </div>
  </Panel>
</template>
<script setup>
import Panel from "primevue/panel";
import Listbox from "primevue/listbox";
import { ref, watch, nextTick } from "vue";
const props = defineProps({
  logs: {
    type: Array,
    default: [],
  },
});
const emit = defineEmits(["hide", "clear"]);
const container = ref(null);
watch(
  () => props.logs,
  (o, n) => {
    nextTick(() => {
      if (container && container.value) {
        const trs = container.value.getElementsByTagName("li");
        trs[trs.length - 1].scrollIntoViewIfNeeded();
      }
    });
  },
  { deep: true }
);
const filter = ref(false);
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
:deep(.p-panel-content) {
  padding: 0.75rem;
}
</style>
