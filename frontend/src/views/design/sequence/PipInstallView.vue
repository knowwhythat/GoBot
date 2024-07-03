<template>
  <div class="flex flex-col">
    <TabView
      class="flex-1"
      :activeIndex="activeIndex"
      @update:activeIndex="tabChange"
    >
      <TabPanel header="安装新模块">
        <div class="flex mt-4">
          <span class="mr-2 w-28 my-auto text-right"> Python包名称 </span>
          <InputText class="flex-auto" v-model="dependency.name" />
        </div>
        <div class="flex mt-4">
          <div class="mr-2 w-28"></div>
          <Checkbox
            class="mr-2 my-auto"
            v-model="dependency.upgradePip"
            :binary="true"
          />
          <span class="mr-2 my-auto text-right"> 使用最新版PIP </span>
        </div>
        <div class="flex mt-4">
          <div class="mr-2 w-28"></div>
          <Checkbox
            class="mr-2 my-auto"
            v-model="dependency.inputVersion"
            :binary="true"
          />
          <span class="mr-2 my-auto text-right"> 指定版本 </span>
          <InputText
            v-if="dependency.inputVersion"
            class="ml-2 flex-auto"
            v-model="dependency.version"
          />
        </div>
        <div class="flex mt-4">
          <div class="mr-2 w-28"></div>
          <Checkbox
            class="mr-2 my-auto"
            :binary="true"
            v-model="dependency.inputMirror"
          />
          <span class="mr-2 my-auto text-right"> 指定镜像 </span>
          <Dropdown
            v-if="dependency.inputMirror"
            editable
            :options="mirrors"
            class="ml-2 flex-auto"
            v-model="dependency.mirrorUrl"
          />
        </div>
        <div class="flex mt-4">
          <span class="mr-2 w-28 text-left"> 输出 </span>
          <Textarea v-model="output" class="flex-auto" rows="5" disabled />
        </div>
      </TabPanel>
      <TabPanel header="管理已安装模块">
        <DataTable
          class="mt-2"
          :loading="searchLoading"
          scrollable
          :globalFilterFields="['name']"
          v-model:filters="filters"
          scrollHeight="350px"
          :value="dependencies"
          size="small"
          showGridlines
          selectionMode="single"
          v-model:selection="selectedRow"
        >
          <template #header>
            <div class="flex flex-row justify-between">
              <IconField iconPosition="left">
                <InputIcon>
                  <i class="pi pi-search" />
                </InputIcon>
                <InputText
                  v-model="filters['global'].value"
                  placeholder="搜索"
                />
              </IconField>
              <div class="flex flex-row">
                <label class="mr-2 my-auto"> 是否显示所有依赖 </label>
                <Checkbox
                  class="my-auto"
                  v-model="showAll"
                  :binary="true"
                  @change="refreshTable"
                />
              </div>
            </div>
          </template>
          <Column field="name" header="依赖名称"></Column>
          <Column field="version" header="版本"></Column>
        </DataTable>
      </TabPanel>
    </TabView>
    <div class="mt-4 flex flex-none justify-end gap-2">
      <Button
        v-if="activeIndex === 0"
        type="button"
        label="安装"
        :loading="loading"
        @click="installDependency"
      ></Button>
      <Button
        v-if="activeIndex === 1"
        :disabled="!selectedRow"
        type="button"
        label="升级模块"
        :loading="loading"
        @click="upgradeDependency"
      ></Button>
      <Button
        v-if="activeIndex === 1"
        :disabled="!selectedRow"
        type="button"
        label="删除模块"
        :loading="loading"
        @click="removeDependency"
      ></Button>

      <Button
        type="button"
        label="关闭"
        severity="secondary"
        @click="$emit('close')"
      ></Button>
    </div>
  </div>
</template>

<script setup>
import IconField from "primevue/iconfield";
import InputIcon from "primevue/inputicon";
import InputText from "primevue/inputtext";
import Button from "primevue/button";
import Checkbox from "primevue/checkbox";
import Dropdown from "primevue/dropdown";
import Textarea from "primevue/textarea";
import TabView from "primevue/tabview";
import TabPanel from "primevue/tabpanel";
import DataTable from "primevue/datatable";
import Column from "primevue/column";
import { inject, ref } from "vue";
import {
  PipInstallPackage,
  PipListInstallPackage,
  PipUnInstallPackage,
  SaveProjectDependency,
} from "@back/go/backend/App";
import { EventsOff, EventsOn } from "@back/runtime/runtime";

import { useToast } from "primevue/usetoast";
import { FilterMatchMode } from "primevue/api";

const toast = useToast();

const emit = defineEmits(["close"]);
const { projectId } = inject("projectId");

const dependency = ref({
  name: "",
  upgradePip: false,
  upgrade: false,
  inputVersion: false,
  version: "",
  inputMirror: false,
  mirrorUrl: "https://pypi.tuna.tsinghua.edu.cn/simple",
});

const mirrors = [
  "https://pypi.tuna.tsinghua.edu.cn/simple",
  "https://pypi.org/simple",
  "https://mirrors.aliyun.com/pypi/simple",
  "https://pypi.douban.com",
  "https://pypi.hustunique.com",
  "https://pypi.mirrors.ustc.edu.cn",
];

const output = ref("");

const dependencies = ref([]);
const selectedRow = ref(null);
const loading = ref(false);
const searchLoading = ref(false);

async function installDependency() {
  if (!dependency.value.name) {
    toast.add({
      severity: "warn",
      detail: "Python包名称不能为空",
      life: 1000,
    });
  }
  output.value = "";
  EventsOn("pipOutput", (data) => {
    output.value = output.value + data;
  });
  loading.value = true;
  try {
    await PipInstallPackage(projectId, dependency.value);
    await refreshTable();
    await saveDependency();
    toast.add({
      severity: "success",
      detail: "安装成功",
      group: "tc",
      life: 1000,
    });
  } catch (err) {
    toast.add({
      severity: "error",
      detail: err,
      life: 1000,
    });
  }
  EventsOff("pipOutput");
  loading.value = false;
}

async function upgradeDependency() {
  loading.value = true;
  let pack = {
    name: selectedRow.value.name,
    upgrade: true,
    upgradePip: dependency.value.upgradePip,
    inputVersion: dependency.value.inputVersion,
    version: dependency.value.version,
    inputMirror: dependency.value.inputMirror,
    mirrorUrl: dependency.value.mirrorUrl,
  };
  try {
    await PipInstallPackage(projectId, pack);
    await refreshTable();
    await saveDependency();
  } catch (err) {
    toast.add({
      severity: "error",
      detail: err,
      life: 1000,
    });
  }
  loading.value = false;
}

async function removeDependency() {
  loading.value = true;
  try {
    await PipUnInstallPackage(projectId, selectedRow.value.name);
    await refreshTable();
    await saveDependency();
  } catch (err) {
    toast.add({
      severity: "error",
      detail: err,
      life: 1000,
    });
  }
  loading.value = false;
}

const activeIndex = ref(0);

async function tabChange(index) {
  activeIndex.value = index;
  if (index === 1) {
    await refreshTable();
  }
}

const filters = ref({
  global: { value: null, matchMode: FilterMatchMode.CONTAINS },
});

const showAll = ref(false);

async function refreshTable() {
  searchLoading.value = true;
  dependencies.value = [];
  try {
    dependencies.value = await PipListInstallPackage(projectId, !showAll.value);
  } finally {
    searchLoading.value = false;
    selectedRow.value = null;
  }
}

async function saveDependency() {
  try {
    let datas = [];
    const dependencies = await PipListInstallPackage(projectId, true);
    for (let i = 0; i < dependencies.length; i++) {
      datas.push(dependencies[i].name + "==" + dependencies[i].version);
    }
    await SaveProjectDependency(projectId, datas);
  } catch (err) {
    toast.add({
      severity: "error",
      detail: err,
      life: 1000,
    });
  }
}
</script>
