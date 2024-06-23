<template>
  <Panel class="h-full">
    <template #header>
      <div class="flex flex-1">
        <div class="flex-none">全局变量</div>
        <div class="flex-1 flex flex-row-reverse -mt-1">
          <div
            class="hover:bg-slate-200 px-1"
            v-tooltip="'新建全局变量'"
            @click="addGlobalVariable"
          >
            <v-remixicon size="20" name="riPlayListAddFill" />
          </div>
        </div>
      </div>
    </template>
    <DataTable :value="variables" size="small" showGridlines>
      <Column field="name" header="变量名称"></Column>
      <Column field="type" header="类型">
        <template #body="slotProps">
          {{ formatVariableType(slotProps.data.type) }}
        </template>
      </Column>
      <Column field="value" header="值"></Column>
    </DataTable>
    <Dialog
      v-model:visible="showGlobalVariableDialog"
      modal
      header="新建全局变量"
      :style="{ width: '35rem' }"
    >
      <div class="flex items-center gap-3 mb-3">
        <label for="variable_name" class="w-16"> 变量名称 </label>
        <InputText
          id="variable_name"
          v-model="newVariable.name"
          class="flex-auto"
          autocomplete="off"
        />
      </div>
      <div class="flex items-center gap-3 mb-5">
        <label for="variable_type" class="w-16">类型</label>
        <Dropdown
          v-model="newVariable.type"
          :options="variableTypes"
          optionValue="code"
          optionLabel="name"
          class="flex-auto"
        />
      </div>
      <div class="flex relative items-center gap-3 mb-5">
        <label for="variable_value" class="w-16">默认值</label>
        <Textarea
          id="variable_value"
          class="flex-auto indent-5"
          v-model="newVariable.value"
          rows="2"
        >
        </Textarea>
        <button class="absolute bottom-10 left-20" @click="chageVariableKind">
          <v-remixicon
            name="riReactjsFill"
            :fill="
              newVariable.kind === 'expression' ? 'var(--primary-color)' : ''
            "
            size="24"
          />
        </button>
      </div>
      <div class="flex justify-end gap-2">
        <Button
          type="button"
          label="取消"
          severity="secondary"
          @click="showGlobalVariableDialog = false"
        ></Button>
        <Button
          type="button"
          label="确定"
          @click="doAddGlobalVariable"
        ></Button>
      </div>
    </Dialog>
  </Panel>
</template>
<script setup>
import DataTable from "primevue/datatable";
import Column from "primevue/column";
import Dropdown from "primevue/dropdown";
import Dialog from "primevue/dialog";
import Panel from "primevue/panel";
import InputText from "primevue/inputtext";
import Button from "primevue/button";
import Textarea from "primevue/textarea";

import { ref, watch } from "vue";

const variables = defineModel();

const variableTypes = [
  { name: "字符串", code: "str" },
  { name: "数字", code: "number" },
  { name: "布尔", code: "bool" },
  { name: "列表", code: "list" },
  { name: "字典", code: "dict" },
  { name: "任意", code: "any" },
];

function formatVariableType(code) {
  return variableTypes.filter((type) => type.code === code)[0].name;
}

const showGlobalVariableDialog = ref(false);
const newVariable = ref({
  name: "",
  type: "",
  value: "",
  kind: "",
});

function chageVariableKind() {
  if (
    newVariable.value.type === "str" &&
    newVariable.value.kind === "expression"
  ) {
    newVariable.value.kind = "text";
  } else if (
    newVariable.value.type === "str" &&
    newVariable.value.kind === "text"
  ) {
    newVariable.value.kind = "expression";
  }
}

watch(
  () => newVariable.value.type,
  (type, prevType) => {
    if (type !== "str") {
      newVariable.value.kind = "expression";
    }
  }
);
function addGlobalVariable() {
  showGlobalVariableDialog.value = true;
  newVariable.value = {
    name: "",
    type: "str",
    value: "",
    kind: "text",
  };
}

function doAddGlobalVariable() {
  showGlobalVariableDialog.value = false;
  variables.value.push(newVariable.value);
}
</script>
