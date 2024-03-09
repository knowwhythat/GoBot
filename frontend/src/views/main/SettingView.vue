<template>
  <div class="flex flex-col gap-2">
    <div class="flex justify-content-center">
      <label for="value" class="mt-3 w-32">数据库存储目录</label>
      <InputText
        id="value"
        v-model="setting.kvdbPath"
        type="text"
        autocomplete="off"
        class="w-full ml-6"
      />
    </div>
    <div class="flex justify-content-center">
      <label for="value" class="mt-3 w-32">日志存储目录</label>
      <InputText
        id="value"
        v-model="setting.logPath"
        type="text"
        autocomplete="off"
        class="w-full ml-6"
      />
    </div>
    <div class="flex justify-content-center">
      <label for="value" class="mt-3 w-32">日志级别</label>
      <InputText
        id="value"
        v-model="setting.logLevel"
        type="text"
        autocomplete="off"
        class="w-full ml-6"
      />
    </div>
    <div class="flex justify-content-center">
      <label for="value" class="mt-3 w-32">项目数据存储目录</label>
      <InputText
        id="value"
        v-model="setting.dataPath"
        type="text"
        autocomplete="off"
        class="w-full ml-6"
      />
    </div>
    <div class="flex justify-content-center">
      <label for="value" class="mt-3 w-32">Python可执行文件目录</label>
      <InputText
        id="value"
        v-model="setting.pythonPath"
        type="text"
        autocomplete="off"
        class="w-full ml-6"
      />
    </div>
    <div class="flex justify-content-center flex-row-reverse">
      <Button>保存</Button>
    </div>
  </div>
</template>
<script setup>
import InputText from "primevue/inputtext";
import Button from "primevue/button";
import { reactive } from "vue";

const setting = reactive({});
</script>
