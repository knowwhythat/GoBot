<template>
  <Dialog
    :visible="visible"
    v-on:update:visible="$emit('hide')"
    maximizable
    modal
    header="元素设置"
    @show="show"
    :style="{ width: '80rem' }"
    :breakpoints="{ '1199px': '75vw', '575px': '90vw' }"
    :pt="{
      header: (options) => ({
        style: {
          padding: '0.75rem',
        },
      }),
    }"
  >
    <Splitter>
      <SplitterPanel :size="30">
        <Listbox
          :model-value="selectedPath"
          @update:modelValue="updateSelected"
          :options="paths"
          class="w-full"
          listStyle="max-height:650px"
        >
          <template #header>
            <div class="text-center mt-1">
              <span class="mr-2">元素名称:</span>
              <InputText v-model="pathOption.displayName" type="text" />
            </div>
          </template>
          <template #option="slotProps">
            <div class="flex align-items-center">
              <Checkbox
                class="mr-2"
                v-model="slotProps.option.enable"
                :binary="true"
              />
              <div class="truncate">
                {{ slotProps.option.localizeControlType }}
                -
                {{ slotProps.option.name }}
              </div>
            </div>
          </template>
        </Listbox>
      </SplitterPanel>
      <SplitterPanel :size="70">
        <DataTable :value="selectedPath.attrs" tableStyle="min-width: 50rem">
          <Column field="enable" header="启用" #body="slotProps">
            <Checkbox v-model="slotProps.data.enable" :binary="true" />
          </Column>
          <Column field="key" header="属性名称"> </Column>
          <Column field="enable" header="匹配方式" #body="slotProps">
            <Dropdown
              v-model="slotProps.data.match"
              :options="matchType"
              optionLabel="name"
              optionValue="value"
              class="w-full md:w-14rem"
            />
          </Column>
          <Column field="value" header="属性值" #body="slotProps">
            <ElementValue class="w-full" v-model="slotProps.data.value" />
          </Column>
        </DataTable>
      </SplitterPanel>
    </Splitter>
    <template #footer>
      <Button label="取消" icon="pi pi-times" @click="$emit('hide')" text />
      <Button label="校验" icon="pi pi-sun" @click="checkData" />
      <Button label="确认" icon="pi pi-check" @click="updateData" />
    </template>
  </Dialog>
</template>

<script setup>
import { StartCheckWindowsElement } from "@back/go/main/App";
import Dialog from "primevue/dialog";
import Splitter from "primevue/splitter";
import SplitterPanel from "primevue/splitterpanel";
import Listbox from "primevue/listbox";
import DataTable from "primevue/datatable";
import Column from "primevue/column";
import Checkbox from "primevue/checkbox";
import InputText from "primevue/inputtext";
import Dropdown from "primevue/dropdown";
import Button from "primevue/button";
import ElementValue from "@/components/element/ElementValue.vue";
import { ref, toRaw, watch } from "vue";
import { useToast } from "primevue/usetoast";
const toast = useToast();

const emit = defineEmits(["hide", "update"]);
const props = defineProps({
  dialogShow: {
    type: Boolean,
    default: false,
  },
  needInit: {
    type: Boolean,
    default: false,
  },
  pathOption: {
    type: Object,
    default: () => ({
      id: "",
      processNmae: "",
      displayName: "",
      image: "",
      paths: [],
    }),
  },
});

const visible = ref(false);
watch(
  () => props.dialogShow,
  (value, oldValue) => {
    visible.value = value;
  },
  { immediate: true }
);
const matchType = ref([
  { name: "相等", value: "equal" },
  { name: "包含", value: "include" },
  { name: "正则", value: "regex" },
]);
const paths = ref([]);
const selectedPath = ref({
  attrs: [],
});

function updateSelected(value) {
  if (value) {
    selectedPath.value = value;
  }
}

function show() {
  if (props.needInit) {
    paths.value = props.pathOption.paths.map((path, index) => {
      let attrs = [];
      for (let [key, value] of Object.entries(path)) {
        if (key !== "id" && key !== "localizeControlType") {
          let enable = false;
          if (index == 0 && key === "processName") {
            enable = true;
          } else if (key != "index" && key !== "processName" && value) {
            enable = true;
          }
          if (key === "index" && path["name"] === "") {
            enable = true;
          }
          attrs.push({
            enable: enable,
            key: key,
            match: "equal",
            value: value,
          });
        }
      }
      path.attrs = attrs;
      path.enable = true;
      return path;
    });
  } else {
    paths.value = props.pathOption.paths;
  }
  selectedPath.value = paths.value[paths.value.length - 1];
}

async function checkData() {
  try {
    const resp = await StartCheckWindowsElement(JSON.stringify(paths.value));
    const result = JSON.parse(resp);
    if (result["result"] === "ok") {
      toast.add({
        severity: "success",
        summary: "校验成功",
        detail: "找到元素",
        life: 3000,
      });
    } else {
      toast.add({
        severity: "error",
        summary: "校验失败",
        detail: "元素未找到",
        life: 3000,
      });
    }
  } catch (err) {
    toast.add({
      severity: "error",
      summary: "校验失败",
      detail: err,
      life: 3000,
    });
  }
}

function updateData() {
  emit("update", {
    id: props.pathOption.id,
    processName: props.pathOption.processName,
    displayName: props.pathOption.displayName,
    image: props.pathOption.image,
    paths: toRaw(paths.value),
  });
}
</script>
<style scoped>
:deep(.p-listbox .p-listbox-list .p-listbox-item.p-highlight) {
  color: #4b5563;
}
</style>
