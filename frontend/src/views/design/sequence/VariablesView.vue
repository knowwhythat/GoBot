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
        <span
          class="mt-1"
          :class="{ 'text-purple-500': activeTab == 'global' }"
          @click="activeTab = 'global'"
          >全局变量</span
        >
        <span
          class="mt-1"
          :class="{ 'text-purple-500': activeTab == 'local' }"
          @click="activeTab = 'local'"
          >局部变量</span
        >
        <IconField iconPosition="left">
          <InputIcon>
            <i class="pi pi-search" />
          </InputIcon>
          <InputText
            class="p-1 pl-8"
            v-model="filters.global.value"
            placeholder="搜索变量名"
          />
        </IconField>
      </div>
    </template>
    <template #icons>
      <div class="flex align-items-center gap-2">
        <button class="mr-4 hover:text-purple-500" @click="addVariable">
          <span class="pi pi-plus-circle"></span>
        </button>
        <button class="mr-4 hover:text-purple-500" @click="removeVariable">
          <span class="pi pi-minus-circle"></span>
        </button>
        <button class="hover:text-purple-500" @click="$emit('hide')">
          <span class="pi pi-chevron-down"></span>
        </button>
      </div>
    </template>
    <div ref="tableContainer">
      <DataTable
        v-if="activeTab == 'global'"
        :value="globalVariables"
        v-model:selection="selectedGlobal"
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
        <Column sortable field="reference" header="引用"></Column>
      </DataTable>
      <DataTable
        v-if="activeTab == 'local'"
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
    </div>
  </Panel>
</template>
<script setup>
import Panel from "primevue/panel";
import IconField from "primevue/iconfield";
import InputIcon from "primevue/inputicon";
import InputText from "primevue/inputtext";
import DataTable from "primevue/datatable";
import Column from "primevue/column";
import { FilterMatchMode } from "primevue/api";
import { ref, onMounted, onUnmounted, nextTick } from "vue";
import { nanoid } from "nanoid";
const emit = defineEmits(["hide"]);
const activeTab = ref("global");
const globalVariables = ref([]);
const selectedGlobal = ref([]);
const localVariables = ref([]);
const selectedLocal = ref([]);

const filters = ref({
  global: { value: null, matchMode: FilterMatchMode.CONTAINS },
});

const tableContainer = ref();
function addVariable() {
  if (activeTab.value == "global") {
    globalVariables.value.push({
      id: nanoid(16),
      name: "g_variable_" + globalVariables.value.length,
      defaultValue: "",
      remark: "",
      reference: 0,
    });
    nextTick(() => {
      const trs = tableContainer.value.getElementsByTagName("tr");
      trs[trs.length - 1].scrollIntoViewIfNeeded();
    });
  }
}
function removeVariable() {
  if (activeTab.value == "global") {
    globalVariables.value = globalVariables.value.filter(
      (v) => !selectedGlobal.value?.includes(v)
    );
    selectedGlobal.value = [];
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
