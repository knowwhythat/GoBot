<template>
  <TabView>
    <TabPanel header="基础设置">
      <div class="flex flex-col gap-2">
        <div class="flex justify-content-center">
          <label class="mt-3 w-64">数据库存储目录</label>
          <FilePathInput
            v-model="baseConfig.db_path"
            type="directory"
            class="w-full ml-6"
          />
        </div>
        <div class="flex justify-content-center">
          <label class="mt-3 w-64">日志存储目录</label>
          <FilePathInput
            v-model="baseConfig.log_path"
            type="directory"
            class="w-full ml-6"
          />
        </div>
        <div class="flex justify-content-center">
          <label class="mt-3 w-64">日志级别</label>
          <Dropdown
            v-model="baseConfig.log_level"
            :options="logLevels"
            checkmark
            :highlightOnSelect="false"
            class="w-full ml-6"
          />
        </div>
        <div class="flex justify-content-center">
          <label class="mt-3 w-64">项目数据存储目录</label>
          <FilePathInput
            v-model="baseConfig.data_path"
            type="directory"
            class="w-full ml-6"
          />
        </div>
        <div class="flex justify-content-center">
          <label class="mt-3 w-64">Python可执行文件目录</label>
          <FilePathInput
            v-model="baseConfig.python_path"
            type="file"
            class="w-full ml-6"
          />
        </div>
        <div class="flex justify-content-center flex-row-reverse">
          <Button @click="saveBaseConfig">保存</Button>
        </div>
      </div>
    </TabPanel>
    <TabPanel header="快捷键设置"></TabPanel>
  </TabView>
</template>
<script setup>
import TabView from "primevue/tabview";
import TabPanel from "primevue/tabpanel";
import Button from "primevue/button";
import Dropdown from "primevue/dropdown";
import { GetBasicConfigData, SaveBasicConfigData } from "@back/go/backend/App";
import FilePathInput from "@/components/common/FilePathInput.vue";

import { onMounted, ref } from "vue";

let baseConfig = ref({});

onMounted(async () => {
  baseConfig.value = await GetBasicConfigData();
});

const logLevels = ref([
  "Trace",
  "Debug",
  "Info",
  "Warn",
  "Error",
  "Fatal",
  "Panic",
  "Disabled",
]);

async function saveBaseConfig() {
  await SaveBasicConfigData(baseConfig.value);
}
</script>
