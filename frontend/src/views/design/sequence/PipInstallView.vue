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
          scrollable
          scrollHeight="350px"
          :value="dependencies"
          size="small"
          showGridlines
          selectionMode="single"
          v-model:selection="selectedRow"
        >
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
        type="button"
        label="升级模块"
        :loading="loading"
        @click="upgradeDependency"
      ></Button>
      <Button
        v-if="activeIndex === 1"
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
import InputText from "primevue/inputtext";
import Button from "primevue/button";
import Checkbox from "primevue/checkbox";
import Dropdown from "primevue/dropdown";
import Textarea from "primevue/textarea";
import TabView from "primevue/tabview";
import TabPanel from "primevue/tabpanel";
import DataTable from "primevue/datatable";
import Column from "primevue/column";
import { ref, inject, markRaw } from "vue";
import {
  PipInstallPackage,
  PipListInstallPackage,
  PipUnInstallPackage,
  SaveProjectDependency,
} from "@back/go/backend/App";
import { EventsOff, EventsOn } from "@back/runtime/runtime";

import { useToast } from "primevue/usetoast";

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
const selectedRow = ref({});
const loading = ref(false);
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
    await saveDepnedency();
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
    await saveDepnedency();
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
    await saveDepnedency();
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

async function refreshTable() {
  const output = await PipListInstallPackage(projectId);
  var lines = output.split(/\r\n|\n/);
  let datas = [];
  for (var i = 2; i < lines.length; i++) {
    var words = lines[i].split(/\s+/);
    if (words.length > 1 && words[0] !== "pip" && words[0] !== "setuptools") {
      datas.push({ name: words[0], version: words[1] });
    }
  }
  dependencies.value = datas;
}

async function saveDepnedency() {
  try {
    let datas = [];
    for (var i = 0; i < dependencies.value.length; i++) {
      datas.push(
        dependencies.value[i].name + "==" + dependencies.value[i].version,
      );
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
