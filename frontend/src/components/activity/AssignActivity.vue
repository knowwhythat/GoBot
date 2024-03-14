<template>
  <ActivityBase
    :toggleable="false"
    :runnable="true"
    :id="props.element.id"
    :breakpoint="props.element.breakpoint"
    :collapsed="true"
    :deleteable="props.element.deleteable"
    :label="props.element.label"
    :icon="props.element.icon_path"
    :color="props.element.color"
    :parameter_define="props.element.parameter_define"
    :parameter="props.element.parameter"
    @delete="$emit('delete', { id: props.element.id })"
    @update="updateData($event)"
    @run="runActivity"
  >
    <template #top>
      <div class="flex">
        <InputText
          :model-value="variableName"
          @update:model-value="chageName($event)"
          class="py-0 px-2 h-8"
          placeholder="变量名称"
          @keydown="(e) => e.stopPropagation()"
        />
        <span class="mx-1 my-auto"> = </span>
        <InputGroup>
          <InputText
            :model-value="variableValue"
            @update:model-value="chageValue($event)"
            class="py-0 px-2 h-8"
            placeholder="变量值"
            @keydown="(e) => e.stopPropagation()"
          />
          <Button
            class="h-8"
            :severity="isExpression ? 'Primary' : 'secondary'"
            @click="changeValueType"
          >
            <template #icon>
              <v-remixicon name="riReactjsFill" size="24" />
            </template>
          </Button>
        </InputGroup>
      </div>
    </template>
  </ActivityBase>
</template>
<script setup>
import { inject, ref, watch } from "vue";
import InputGroup from "primevue/inputgroup";
import InputText from "primevue/inputtext";
import Button from "primevue/button";
import ActivityBase from "./ActivityBase.vue";
import { RunActivity } from "@back/go/backend/App";
import { EventsOn, EventsOff } from "@back/runtime/runtime.js";
import { useToast } from "primevue/usetoast";
const toast = useToast();
const props = defineProps({
  element: {
    type: Object,
    default: {},
  },
});
const emit = defineEmits(["delete"]);
const isExpression = ref(false);
const variableName = ref("");
const variableValue = ref("");
watch(
  () => props.element.parameter,
  (value, oldValue) => {
    variableName.value = value["variable_name"];
    if (
      "variable_value" in value &&
      value["variable_value"].split(":").length > 1
    ) {
      isExpression.value = value["variable_value"].split(":")[0] === "1";
      variableValue.value = value["variable_value"].substring(2);
    } else {
      const defaultValue = props.element.parameter_define.inputs.filter(
        (pd) => pd.key === "variable_value"
      )[0].default_value;
      isExpression.value = defaultValue.split(":")[0] === "1";
      variableValue.value = defaultValue.substring(2);
    }
  },
  { immediate: true, deep: true }
);
const { dataChanged, updateDataChanged } = inject("dataChanged");
function chageName(data) {
  props.element.parameter["variable_name"] = data;
  updateDataChanged();
}
function chageValue(data) {
  if (isExpression.value) {
    props.element.parameter["variable_value"] = "1:" + data;
  } else {
    props.element.parameter["variable_value"] = "0:" + data;
  }
  updateDataChanged();
}
function changeValueType() {
  isExpression.value = !isExpression.value;
  if (isExpression.value) {
    props.element.parameter["variable_value"] = "1:" + variableValue.value;
  } else {
    props.element.parameter["variable_value"] = "0:" + variableValue.value;
  }
  updateDataChanged();
}
function updateData(data) {
  for (const key in data) {
    if (Object.hasOwnProperty.call(data, key)) {
      props.element[key] = data[key];
    }
  }
}

const { projectId } = inject("projectId");

async function runActivity() {
  EventsOn("repl", (data, out) => {
    if (out) {
      toast.add({
        severity: "error",
        summary: "执行失败",
        detail: out,
        life: 3000,
      });
    } else {
      toast.add({
        severity: "success",
        summary: "执行成功",
        life: 3000,
      });
    }
    EventsOff("repl");
  });
  await RunActivity(projectId, JSON.stringify(props.element));
}
</script>
