<template>
  <TabView>
    <TabPanel header="常规">
      <div class="flex flex-row mb-4">
        <span class="w-28 my-auto mx-2"> 流程块名称 </span>
        <InputText class="w-full" type="text" v-model="blockLabel" />
      </div>
      <Divider align="left" type="solid">
        <b>输入参数</b>
      </Divider>
      <template v-for="input in inputParams" :key="input.name">
        <div class="flex mb-4">
          <span class="mr-2 w-32 truncate my-auto" :title="input.name">
            {{ input.name }}
          </span>
          <VariableInput class="w-full" v-model="params[input.name]" />
        </div>
      </template>
      <Divider align="left" type="solid">
        <b>输出参数</b>
      </Divider>
      <template v-for="output in outputParams" :key="output.name">
        <div class="flex mb-4">
          <span class="mr-2 w-32 truncate my-auto" :title="output.name">
            {{ output.name }}
          </span>
        </div>
      </template>
    </TabPanel>
    <TabPanel header="错误处理">
      <div
        class="flex items-start rounded-lg bg-green-200 p-4 text-black dark:bg-green-300"
      >
        <v-remixicon name="riInformationLine" />
        <p class="ml-4 flex-1">以下这些规则适用于模块发生错误时</p>
      </div>
      <div class="mt-4">
        <div class="inline-flex">
          <InputSwitch v-model="errorHandler.errorEnable" />
          <span class="ml-2"> 启用 </span>
        </div>

        <div v-if="errorHandler.errorEnable" class="mt-4">
          <div class="inline-flex">
            <InputSwitch v-model="errorHandler.retry" />
            <span class="ml-2"> 重试操作 </span>
          </div>
          <div v-if="errorHandler.retry" class="flex flex-row mt-2">
            <div class="items-center">
              <span> 次数 </span>
              <v-remixicon
                title="重试操作的次数"
                name="riInformationLine"
                size="20"
                class="mr-2"
              />
              <InputNumber
                v-model="errorHandler.retryCount"
                inputId="retryCount"
                :min="0"
                showButtons
              />
            </div>
            <div class="items-center ml-2">
              <span> 间隔(s) </span>
              <v-remixicon
                title="每次尝试等待的时间"
                name="riInformationLine"
                size="20"
                class="mr-2"
              />
              <InputNumber
                v-model="errorHandler.retryTimeout"
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
              v-model="errorHandler.toDo"
              :options="toDoTypes"
              optionLabel="name"
              option-value="value"
              class="mt-1 w-full"
            />
          </div>
        </div>
      </div>
    </TabPanel>
  </TabView>
</template>
<script setup>
import TabView from "primevue/tabview";
import TabPanel from "primevue/tabpanel";
import InputText from "primevue/inputtext";
import InputSwitch from "primevue/inputswitch";
import InputNumber from "primevue/inputnumber";
import Dropdown from "primevue/dropdown";
import VariableInput from "@/views/design/flow/editor/VariableInput.vue";
import { computed, inject, onMounted, ref } from "vue";
import { GetSubFlow } from "@back/go/backend/App";
import Divider from "primevue/divider";

const props = defineProps({
  blockId: {
    type: String,
    default: "",
  },
});

const errorHandler = defineModel("errorHandler", {
  type: Object,
});
const params = defineModel("params", {
  type: Object,
});
const blockLabel = defineModel("blockLabel", { type: String });
const { projectId } = inject("projectId");
const flowParams = ref([]);

const inputParams = computed(() => {
  return flowParams.value.filter((item) => item.direction === "In");
});

const outputParams = computed(() => {
  return flowParams.value.filter((item) => item.direction === "Oun");
});

onMounted(async () => {
  try {
    const result = await GetSubFlow(projectId, props.blockId);
    if (result) {
      flowParams.value = JSON.parse(result).params;
      for (let i = 0; i < flowParams.value.length; i++) {
        const item = flowParams.value[i];
        if (!(item.name in params.value)) {
          params.value[item.name] = item.value;
        }
      }
    }
    console.log(flowParams.value);
    console.log(params.value);
  } catch (error) {
    toast.add({
      severity: "error",
      summary: "加载流程异常",
      detail: error,
      life: 3000,
    });
  }
});

const toDoTypes = [
  { value: "error", name: "抛出错误" },
  { value: "continue", name: "继续执行" },
  { value: "fallback", name: "执行回退" },
];
</script>
