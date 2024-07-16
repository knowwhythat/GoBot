<template>
  <ActivityBase
    :runnable="true"
    :toggleable="false"
    :id="props.element.id"
    :breakpoint="props.element.breakpoint"
    :collapsed="true"
    :deleteable="props.element.deleteable"
    :commentable="true"
    :label="props.element.label"
    :icon="props.element.icon_path"
    :color="props.element.color"
    :parameter_define="props.element.parameter_define"
    :parameter="props.element.parameter"
    @delete="$emit('delete', { id: props.element.id })"
    @comment="$emit('comment', { ...props.element })"
    @uncomment="$emit('uncomment', { ...props.element })"
    @update="updateData($event)"
    @run="runActivity"
  >
    <template #top>
      <div>
        <InputGroup>
          <InputText
            class="py-0 px-2 h-8"
            :model-value="expression"
            @update:model-value="chageExpression($event)"
            placeholder="输入内容或表达式"
            @keydown="(e) => e.stopPropagation()"
          />
          <Button
            class="h-8"
            :severity="isExpression ? 'Primary' : 'secondary'"
            @click="chageIsExpression"
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
import { ref, watch, inject } from "vue";
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
const emit = defineEmits(["delete", "comment", "uncomment"]);
const { dataChanged, updateDataChanged } = inject("dataChanged");

const isExpression = ref(false);
const expression = ref("");
watch(
  () => props.element.parameter,
  (value, oldValue) => {
    if ("expression" in value && value["expression"].split(":").length > 1) {
      isExpression.value = value["expression"].split(":")[0] === "1";
      expression.value = value["expression"].substring(2);
    } else {
      const defaultValue = props.element.parameter_define.inputs.filter(
        (pd) => pd.key === "expression",
      )[0].default_value;
      isExpression.value = defaultValue.split(":")[0] === "1";
      expression.value = defaultValue.substring(2);
    }
  },
  { immediate: true, deep: true },
);

function chageExpression(data) {
  if (isExpression.value) {
    props.element.parameter["expression"] = "1:" + data;
  } else {
    props.element.parameter["expression"] = "0:" + data;
  }
  updateDataChanged();
}
function chageIsExpression() {
  isExpression.value = !isExpression.value;
  if (isExpression.value) {
    props.element.parameter["expression"] = "1:" + expression.value;
  } else {
    props.element.parameter["expression"] = "0:" + expression.value;
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
