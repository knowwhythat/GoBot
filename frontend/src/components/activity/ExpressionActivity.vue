<template>
  <ActivityBase
    :toggleable="true"
    :id="props.element.id"
    :breakpoint="props.element.breakpoint"
    :collapsed="props.element.collapsed"
    :deleteable="props.element.deleteable"
    :label="props.element.label"
    :icon="props.element.icon_path"
    :color="props.element.color"
    :parameter_define="props.element.parameter_define"
    :parameter="props.element.parameter"
    @delete="$emit('delete', { id: props.element.id })"
    @update="updateData($event)"
  >
    <InputGroup class="mb-2">
      <InputText
        :model-value="expression"
        @update:model-value="chageExpression($event)"
        placeholder="输入内容或表达式"
      />
      <Button
        :severity="isExpression ? 'Primary' : 'secondary'"
        @click="chageIsExpression"
      >
        <template #icon>
          <v-remixicon name="riReactjsFill" size="24" />
        </template>
      </Button>
    </InputGroup>
  </ActivityBase>
</template>
<script setup>
import { ref, watch, inject } from "vue";
import InputGroup from "primevue/inputgroup";
import InputText from "primevue/inputtext";
import Button from "primevue/button";
import ActivityBase from "./ActivityBase.vue";
const props = defineProps({
  element: {
    type: Object,
    default: {},
  },
});
const emit = defineEmits(["delete"]);
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
        (pd) => pd.key === "expression"
      )[0].default_value;
      isExpression.value = defaultValue.split(":")[0] === "1";
      expression.value = defaultValue.substring(2);
    }
  },
  { immediate: true, deep: true }
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
</script>
