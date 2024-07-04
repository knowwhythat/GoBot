<template>
  <div v-if="state.data" class="mb-4">
    <span class="ml-2 mr-2"> 模块名称 </span>
    <InputText type="text" v-model="state.data.label" />
  </div>
  <div
    class="flex items-start rounded-lg bg-green-200 p-4 text-black dark:bg-green-300"
  >
    <v-remixicon name="riInformationLine" />
    <p class="ml-4 flex-1">以下这些规则适用于模块发生错误时</p>
  </div>
  <div v-if="state.data" class="mt-4">
    <div class="inline-flex">
      <InputSwitch v-model:model-value="state.data.errorEnable" />
      <span class="ml-2"> 启用 </span>
    </div>

    <div v-if="state.data.errorEnable" class="mt-4">
      <div class="inline-flex">
        <InputSwitch v-model="state.data.retry" />
        <span class="ml-2"> 重试操作 </span>
      </div>
      <div v-if="state.data.retry" class="mt-2">
        <div class="inline-flex items-center">
          <span> 次数 </span>
          <v-remixicon
            title="重试操作的次数"
            name="riInformationLine"
            size="20"
            class="mr-2"
          />
          <InputNumber
            v-model="state.data.retryCount"
            inputId="retryCount"
            :min="0"
            showButtons
          />
        </div>
        <div class="inline-flex items-center ml-2">
          <span> 间隔(s) </span>
          <v-remixicon
            title="每次尝试等待的时间"
            name="riInformationLine"
            size="20"
            class="mr-2"
          />
          <InputNumber
            v-model="state.data.retryTimeout"
            inputId="retryTimeout"
            :min="0"
            showButtons
          />
        </div>
      </div>
      <div class="mt-4">
        <div class="inline-flex items-center">
          <span class="ml-2 mr-2"> 异常策略 </span>
          <v-remixicon
            title="运行报错时的执行策略"
            name="riInformationLine"
            size="20"
            class="mr-2"
          />
        </div>
        <Dropdown
          v-model="state.data.toDo"
          :options="toDoTypes"
          optionLabel="name"
          option-value="value"
          class="mt-1 w-full"
        />
      </div>
    </div>
  </div>
</template>
<script setup>
import { ref, reactive, onMounted, watch } from "vue";
import InputText from "primevue/inputtext";
import InputSwitch from "primevue/inputswitch";
import InputNumber from "primevue/inputnumber";
import Dropdown from "primevue/dropdown";
const props = defineProps({
  data: {
    type: Object,
    default: () => ({}),
  },
});
const emit = defineEmits(["change"]);
const state = reactive({});
const toDoTypes = [
  { value: "error", name: "抛出错误" },
  { value: "continue", name: "继续执行" },
  { value: "fallback", name: "执行回退" },
];
watch(
  () => state.data,
  (data, newData) => {
    emit("change", newData);
  },
  { deep: true },
);
onMounted(() => {
  Object.assign(state, props.data);
});
</script>
