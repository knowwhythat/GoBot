<template>
  <Panel
    header="参数库"
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
        <button class="mr-4 hover:text-purple-500" @click="addParam">
          <span class="pi pi-plus-circle"></span>
        </button>
        <button class="hover:text-purple-500" @click="$emit('hide')">
          <span class="pi pi-chevron-down"></span>
        </button>
      </div>
    </template>
    <div ref="tableContainer">
      <DataTable
        :value="flowParams"
        v-model:selection="selectedRow"
        showGridlines
        scrollable
        scrollHeight="flex"
        editMode="cell"
        @cell-edit-complete="onCellEditComplete"
        v-model:filters="filters"
        :globalFilterFields="['name']"
        :pt="{
          column: {
            bodycell: ({ state }) => ({
              class: [{ 'pt-0 pb-0': state['d_editing'] }],
            }),
          },
        }"
      >
        <Column selectionMode="multiple" headerStyle="width: 2rem"></Column>
        <Column sortable field="name" header="变量名" style="width: 25%">
          <template #editor="{ data, field }">
            <InputText v-model="data[field]" autofocus />
          </template>
        </Column>
        <Column field="defaultValue" header="默认值" style="width: 25%">
          <template #editor="{ data, field }">
            <InputText v-model="data[field]" autofocus />
          </template>
        </Column>
        <Column field="remark" header="备注" style="width: 35%">
          <template #editor="{ data, field }">
            <InputText v-model="data[field]" autofocus />
          </template>
        </Column>
        <Column field="remark" header="备注" style="width: 15%"></Column>
      </DataTable>
    </div>
  </Panel>
</template>
<script setup>
import Panel from "primevue/panel";
import InputText from "primevue/inputtext";
import DataTable from "primevue/datatable";
import Column from "primevue/column";
import { FilterMatchMode } from "primevue/api";
import { ref, onMounted, onUnmounted, nextTick } from "vue";
import { nanoid } from "nanoid";
const emit = defineEmits(["hide"]);
const flowParams = ref([]);
const selectedRow = ref([]);

const filters = ref({
  global: { value: null, matchMode: FilterMatchMode.CONTAINS },
});

const tableContainer = ref();
function addParam() {
  flowParams.value.push({
    id: nanoid(16),
    name: "g_variable_" + flowParams.value.length,
    defaultValue: "",
    remark: "",
    reference: 0,
  });
  nextTick(() => {
    const trs = tableContainer.value.getElementsByTagName("tr");
    trs[trs.length - 1].scrollIntoViewIfNeeded();
  });
}
function removeVariable() {
  if (activeTab.value == "global") {
    flowParams.value = flowParams.value.filter(
      (v) => !selectedRow.value?.includes(v)
    );
    selectedRow.value = [];
  }
}
const onCellEditComplete = (event) => {
  let { data, newValue, field } = event;
  data[field] = newValue;
  console.log(event);
};
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
