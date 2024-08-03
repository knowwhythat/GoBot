<template>
  <div class="flex flex-col w-full">
    <InputGroup>
      <Dropdown
        v-model="content"
        :options="subFlows"
        optionLabel="label"
        optionValue="key"
        :filter="subFlows.length > 8"
      ></Dropdown>
    </InputGroup>
    <Divider align="left" type="solid" v-if="inputParams.length > 0">
      <b>流程输入参数</b>
    </Divider>
    <template v-for="input in inputParams" :key="input.name">
      <div class="flex mb-4">
        <span class="mr-2 w-32 truncate my-auto" :title="input.name">
          {{ input.name }}
        </span>
        <ExpressionTextInput
          :input="{}"
          :value="getParamValue('inputs', input.name)"
          @update="updateInput(input.name, $event)"
        />
      </div>
    </template>
    <Divider align="left" type="solid" v-if="outputParams.length > 0">
      <b>流程输出参数</b>
    </Divider>
    <template v-for="output in outputParams" :key="output.name">
      <div class="flex mb-4">
        <span class="mr-2 w-32 truncate my-auto" :title="output.name">
          {{ output.name }}
        </span>
        <SelectOrInput
          :input="{}"
          :value="getParamValue('outputs', output.name)"
          @update="updateOutput(output.name, $event)"
        />
      </div>
    </template>
  </div>
</template>
<script setup>
import InputGroup from "primevue/inputgroup";
import Dropdown from "primevue/dropdown";
import { computed, inject, onBeforeMount, ref, watch } from "vue";
import { GetSubFlow, ReadProjectConfig } from "@back/go/backend/App";
import ExpressionTextInput from "@/components/editor/ExpressionTextInput.vue";
import SelectOrInput from "@/components/editor/SelectOrInput.vue";
import Divider from "primevue/divider";
import { throttle } from "@/utils/helper.js";

const props = defineProps({
  parameter: {
    type: Object,
    default: () => ({}),
  },
});
const emit = defineEmits(["update"]);
const projectConfig = ref({});
const { projectId } = inject("projectId");
onBeforeMount(async () => {
  projectConfig.value = await ReadProjectConfig(projectId);
  if (content.value) {
    const result = await GetSubFlow(projectId, content.value);
    flowParams.value = JSON.parse(result).params;
  }
});

const subFlows = computed(() => {
  if (projectConfig.value?.children) {
    return projectConfig.value?.children.filter((flow) => flow.key !== "main");
  }
  return [];
});

const content = computed({
  get() {
    if (
      props.parameter["flow_name"] &&
      props.parameter["flow_name"].split(":").length > 1
    ) {
      return props.parameter["flow_name"].substring(2);
    } else {
      return props.parameter["flow_name"];
    }
  },
  set(value) {
    emit("update", { key: "flow_name", value: "0:" + value });
  },
});

const flowParams = ref([]);

const inputParams = computed(() => {
  return flowParams.value.filter((param) => param.direction === "In");
});

const outputParams = computed(() => {
  return flowParams.value.filter((param) => param.direction === "Out");
});

watch(
  () => content.value,
  async (now, old) => {
    const result = await GetSubFlow(projectId, now);
    flowParams.value = JSON.parse(result).params;
  },
  { deep: true },
);

function getParamValue(category, name) {
  if (props.parameter[category]) {
    return JSON.parse(props.parameter[category])[name];
  }
  return "";
}

function updateInput(key, value) {
  throttle(function () {
    if (props.parameter["inputs"]) {
      const inputs = JSON.parse(props.parameter["inputs"]);
      emit("update", {
        key: "inputs",
        value: JSON.stringify({ ...inputs, [key]: value }),
      });
    } else {
      emit("update", {
        key: "inputs",
        value: JSON.stringify({ [key]: value }),
      });
    }
  }, 500)();
}

function updateOutput(key, value) {
  throttle(function () {
    if (props.parameter["outputs"]) {
      const outputs = JSON.parse(props.parameter["outputs"]);
      emit("update", {
        key: "outputs",
        value: JSON.stringify({ ...outputs, [key]: value }),
      });
    } else {
      emit("update", {
        key: "outputs",
        value: JSON.stringify({ [key]: value }),
      });
    }
  }, 500)();
}
</script>
