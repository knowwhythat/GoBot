<template>
  <Panel
    header="变量"
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
        <button class="hover:text-purple-500" @click="$emit('hide')">
          <span class="pi pi-chevron-down"></span>
        </button>
      </div>
    </template>
    <DataTable
      :value="localVariables"
      v-model:selection="selectedLocal"
      showGridlines
      selectionMode="single"
      scrollable
      scrollHeight="flex"
      v-model:filters="filters"
      :globalFilterFields="['name']"
    >
      <Column field="name" header="变量名称"></Column>
    </DataTable>
  </Panel>
</template>
<script setup>
import Panel from "primevue/panel";
import DataTable from "primevue/datatable";
import Column from "primevue/column";
import { FilterMatchMode } from "primevue/api";
import { ref } from "vue";

const emit = defineEmits(["hide"]);
const localVariables = ref([]);
const selectedLocal = ref([]);

const filters = ref({
  global: { value: null, matchMode: FilterMatchMode.CONTAINS },
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
}
:deep(.p-panel-content) {
  height: 100%;
  overflow: auto;
}
:deep(.p-datatable .p-datatable-thead > tr > th) {
  padding: 4px;
  font-weight: 500;
}
:deep(.p-datatable .p-datatable-tbody > tr > td) {
  padding: 8px;
}
:deep(.p-panel-content) {
  padding: 0.75rem;
}
</style>
