<template>
  <div class="flex flex-col w-full">
    <InputGroup>
      <Dropdown
        v-model="content"
        :options="subModules"
        optionLabel="label"
        optionValue="key"
        :filter="subModules.length > 8"
      ></Dropdown>
    </InputGroup>
    <div class="flex mt-4">
      <span class="mr-2 w-32 truncate my-auto" title="方法名称">
        方法名称
      </span>
      <InputGroup>
        <Dropdown
          v-model="methodName"
          :options="funcDefine"
          optionLabel="name"
          optionValue="name"
          :filter="funcDefine.length > 8"
        ></Dropdown>
      </InputGroup>
    </div>

    <Divider align="left" type="solid" v-if="inputParams.length > 0">
      <b>方法输入参数</b>
    </Divider>
    <template v-for="input in inputParams" :key="input">
      <div class="flex mb-4">
        <span class="mr-2 w-32 truncate my-auto" :title="input">
          {{ input }}
        </span>
        <ExpressionTextInput
          :input="{}"
          :value="getInputValue(input)"
          @update="updateInput(input, $event)"
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
import ExpressionTextInput from "@/components/editor/ExpressionTextInputOld.vue";
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
    funcDefine.value = JSON.parse(result).funcDefine;
  }
});

const subModules = computed(() => {
  if (projectConfig.value?.children) {
    console.log(projectConfig.value?.children);
    return projectConfig.value?.children.filter(
      (flow) => flow.key !== "main" && flow.nodeType === "code",
    );
  }
  return [];
});

const content = computed({
  get() {
    if (
      props.parameter["module_name"] &&
      props.parameter["module_name"].split(":").length > 1
    ) {
      return props.parameter["module_name"].substring(2);
    } else {
      return props.parameter["module_name"];
    }
  },
  set(value) {
    emit("update", { key: "module_name", value: "0:" + value });
    emit("update", {
      key: "inputs",
      value: JSON.stringify({}),
    });
  },
});

watch(
  () => content.value,
  async (now, old) => {
    const result = await GetSubFlow(projectId, now);
    funcDefine.value = JSON.parse(result).funcDefine;
  },
  { deep: true },
);

const methodName = computed({
  get() {
    if (props.parameter["method_name"]) {
      return props.parameter["method_name"];
    }
    return "";
  },
  set(value) {
    emit("update", { key: "method_name", value: value });
    emit("update", {
      key: "inputs",
      value: JSON.stringify({}),
    });
  },
});
const funcDefine = ref([]);

const inputParams = computed(() => {
  const func = funcDefine.value.find((func) => func.name === methodName.value);
  if (func) {
    return func["arguments"];
  }
  return [];
});

function getInputValue(name) {
  if (props.parameter["inputs"]) {
    return JSON.parse(props.parameter["inputs"])[name];
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
</script>
