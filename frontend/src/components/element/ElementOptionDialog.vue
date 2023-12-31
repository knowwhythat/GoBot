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
          v-model="selectedPath"
          :options="paths"
          class="w-full"
          listStyle="max-height:650px"
        >
          <template #header>
            <div class="text-center">元素节点</div>
          </template>
          <template #option="slotProps">
            <p class="truncate">
              {{ slotProps.option.localizeControlType }}
              -
              {{ slotProps.option.name }}
            </p>
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
            <InputText class="w-full" v-model="slotProps.data.value" />
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
import { ref, watch } from "vue";

const emit = defineEmits(["hide", "update", "check"]);
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
      return path;
    });
  } else {
    paths.value = props.pathOption.paths;
  }
  selectedPath.value = paths.value[0];
}
function checkData() {
  emit("check", paths.value);
}
function updateData() {
  emit("update", paths.value);
}
</script>
<style scoped>
:deep(.p-listbox .p-listbox-list .p-listbox-item.p-highlight) {
  color: #4b5563;
}
</style>
