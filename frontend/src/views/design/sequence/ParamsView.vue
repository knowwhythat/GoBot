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
        <button
          v-if="selectedRow"
          class="mr-4 hover:text-purple-500"
          v-tooltip.left="{
            value: '删除参数',
            showDelay: 100,
            hideDelay: 300,
          }"
          @click="removeparam"
        >
          <span class="pi pi-delete-left"></span>
        </button>
        <button
          v-if="selectedRow"
          class="mr-4 hover:text-purple-500"
          v-tooltip.left="{
            value: '修改参数',
            showDelay: 100,
            hideDelay: 300,
          }"
          @click="editParam"
        >
          <span class="pi pi-file-edit"></span>
        </button>
        <button
          class="mr-4 hover:text-purple-500"
          v-tooltip.left="{
            value: '新建参数',
            showDelay: 100,
            hideDelay: 300,
          }"
          @click="addParam"
        >
          <span class="pi pi-plus-circle"></span>
        </button>
        <button class="hover:text-purple-500" @click="$emit('hide')">
          <span class="pi pi-chevron-down"></span>
        </button>
      </div>
    </template>
    <DataTable
      class="p-2"
      size="small"
      :value="flowParams"
      v-model:selection="selectedRow"
      showGridlines
      scrollable
      resizableColumns
      columnResizeMode="fit"
      selectionMode="single"
      scrollHeight="flex"
      v-model:filters="filters"
      :globalFilterFields="['name']"
      @rowReorder="flowParams = $event.value"
      :pt="{
        column: {
          bodycell: ({ state }) => ({
            class: [{ 'pt-0 pb-0': state['d_editing'] }],
          }),
        },
      }"
    >
      <Column rowReorder headerStyle="width: 3rem" />
      <Column field="name" header="参数名称"></Column>
      <Column field="direction" header="方向">
        <template #body="slotProps">
          {{ slotProps.data.direction === "In" ? "输入" : "输出" }}
        </template>
      </Column>
      <Column field="type" header="类型">
        <template #body="slotProps">
          {{
            paramTypes.filter((type) => slotProps.data.type === type.code)[0]
              .name
          }}
        </template>
      </Column>
      <Column field="value" header="默认值"></Column>
    </DataTable>
    <Dialog
      v-model:visible="showParamDialog"
      modal
      header="参数"
      :style="{ width: '35rem' }"
    >
      <div class="flex items-center gap-3 mb-3">
        <label for="param_name" class="w-16"> 参数名称 </label>
        <InputText
          id="param_name"
          v-model="newParam.name"
          class="flex-auto"
          autocomplete="off"
        />
      </div>
      <div class="flex items-center gap-3 mb-5">
        <label for="param_direction" class="w-16">方向</label>
        <Dropdown
          v-model="newParam.direction"
          :options="directionTypes"
          optionValue="code"
          optionLabel="name"
          class="flex-auto"
        />
      </div>
      <div class="flex items-center gap-3">
        <label for="param_type" class="w-16">类型</label>
        <Dropdown
          v-model="newParam.type"
          :options="paramTypes"
          optionValue="code"
          optionLabel="name"
          class="flex-auto"
        />
      </div>
      <template
        v-if="newParam.type === 'single' || newParam.type === 'multiple'"
      >
        <div
          class="flex items-center gap-3 mt-5"
          v-for="(item, i) in newParam.options"
          :key="item.id"
        >
          <label for="param_type" class="w-16">选项{{ i + 1 }}</label>
          <InputText
            v-model="item.key"
            class="w-1/2 mx-1"
            autocomplete="off"
            placeholder="选项标签"
          />
          <InputText
            v-model="item.value"
            class="w-1/2 mx-1"
            autocomplete="off"
            placeholder="选项值"
          />
          <button @click="removeOption(item.id)" v-tooltip="'删除选项'">
            <v-remixicon name="riDeleteBin7Line" size="24" />
          </button>
        </div>
        <div class="flex items-center">
          <button @click="addOption" v-tooltip="'添加选项'">
            <v-remixicon name="riAddLine" size="24" />
          </button>
        </div>
      </template>

      <div class="flex relative items-center gap-3 mb-5 mt-5">
        <label for="param_value" class="w-16">默认值</label>
        <Textarea
          id="param_value"
          class="flex-auto"
          v-model="newParam.value"
          rows="2"
        >
        </Textarea>
      </div>
      <div class="flex justify-end gap-2">
        <Button
          type="button"
          label="取消"
          severity="secondary"
          @click="showParamDialog = false"
        ></Button>
        <Button type="button" label="确定" @click="doAddParam"></Button>
      </div>
    </Dialog>
  </Panel>
</template>
<script setup>
import Panel from "primevue/panel";
import InputText from "primevue/inputtext";
import Dropdown from "primevue/dropdown";
import Dialog from "primevue/dialog";
import Button from "primevue/button";
import DataTable from "primevue/datatable";
import Column from "primevue/column";
import { paramTypes } from "@/utils/shared.js";
import { FilterMatchMode } from "primevue/api";
import { ref } from "vue";
import { nanoid } from "nanoid";
import { useToast } from "primevue/usetoast";
import { useConfirm } from "primevue/useconfirm";

const toast = useToast();

const confirm = useConfirm();
const emit = defineEmits(["hide"]);

const flowParams = defineModel({
  type: Array,
  default: () => [],
});

const selectedRow = ref();

const filters = ref({
  global: { value: null, matchMode: FilterMatchMode.CONTAINS },
});

const showParamDialog = ref(false);
const newParam = ref({});
const directionTypes = [
  { name: "输入", code: "In" },
  { name: "输出", code: "Out" },
];

function addParam() {
  newParam.value = {
    name: "",
    direction: "In",
    type: "str",
    options: [],
    value: "",
  };
  showParamDialog.value = true;
}

let editIndex = -1;

function editParam() {
  editIndex = flowParams.value.findIndex(
    (v) => v.name === selectedRow.value.name,
  );
  newParam.value = {
    name: selectedRow.value.name,
    direction: selectedRow.value.direction,
    type: selectedRow.value.type,
    options: selectedRow.value.options,
    value: selectedRow.value.value,
  };

  showParamDialog.value = true;
}

function addOption() {
  newParam.value.options.push({ id: nanoid(8), key: "", value: "" });
}

function removeOption(id) {}

function doAddParam() {
  if (!newParam.value.name) {
    toast.add({
      severity: "warn",
      detail: "参数名称不能为空",
      life: 1000,
    });
    return;
  }
  if (editIndex > -1) {
    let index = flowParams.value.findIndex(
      (v, currentIndex) =>
        v.name === newParam.value.name && currentIndex !== editIndex,
    );
    if (index === -1) {
      flowParams.value[editIndex] = newParam.value;
    } else {
      toast.add({
        severity: "warn",
        detail: "变量名称重复",
        life: 1000,
      });
      return;
    }
  } else {
    let index = flowParams.value.findIndex(
      (v) => v.name === newParam.value.name,
    );
    if (index === -1) {
      flowParams.value.push(newParam.value);
    } else {
      toast.add({
        severity: "warn",
        detail: "参数名称重复",
        life: 1000,
      });
      return;
    }
  }
  selectedRow.value = null;
  showParamDialog.value = false;
}

function removeparam() {
  confirm.require({
    header: "提示",
    message: "确定要删除这条记录吗?",
    icon: "pi pi-info-circle",
    rejectClass: "p-button-secondary p-button-outlined p-button-sm",
    acceptClass: "p-button-danger p-button-sm",
    rejectLabel: "取消",
    acceptLabel: "删除",
    accept: () => {
      flowParams.value = flowParams.value.filter(
        (v) => v !== selectedRow.value,
      );
      selectedRow.value = null;
    },
    reject: () => {},
  });
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
}

:deep(.p-datatable .p-datatable-thead > tr > th) {
  padding: 4px;
  font-weight: 500;
}

:deep(.p-panel-content) {
  padding: 0.75rem;
}
</style>
