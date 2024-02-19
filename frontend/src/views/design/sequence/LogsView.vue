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
      <div>
        <button class="mr-4 hover:text-purple-500" @click="logs = []">
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
    <Listbox
      :options="logs"
      :filter="filter"
      emptyMessage="没有日志"
      listStyle="height: calc(100% - 30px);"
    />
  </Panel>
</template>
<script setup>
import Panel from "primevue/panel";
import Listbox from "primevue/listbox";
import { ref, onMounted, onUnmounted } from "vue";
import { EventsOn, EventsOff } from "@back/runtime/runtime";
const emit = defineEmits(["hide"]);
const filter = ref(false);
const logs = ref([]);
onMounted(() => {
  EventsOn("log_event", (data) => {
    logs.value.push(data);
  });
});

onUnmounted(() => {
  EventsOff("log_event");
});
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
  overflow: auto;
}
</style>
