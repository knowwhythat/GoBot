<template>
  <ActivityBase
    :toggleable="true"
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
    <div class="flex mx-2">
      <InputText
        :model-value="variableName"
        @update:model-value="chageName($event)"
        class="w-32"
        placeholder="变量名称"
      />
      <span class="mx-1 my-auto"> = </span>
      <InputGroup>
        <InputText
          :model-value="variableValue"
          @update:model-value="chageValue($event)"
          class="w-32"
          placeholder="变量值"
        />
        <Button
          :severity="isExpression ? 'Primary' : 'secondary'"
          @click="changeValueType"
        >
          <template #icon>
            <v-remixicon name="riReactjsFill" size="24" />
          </template>
        </Button>
      </InputGroup>
    </div>
  </ActivityBase>
</template>
<script setup>
import { onMounted, ref, watch } from "vue";
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
      isExpression.value = false;
      variableValue.value = value["variable_value"];
    }
  },
  { immediate: true, deep: true }
);
function chageName(data) {
  props.element.parameter["variable_name"] = data;
}
function chageValue(data) {
  if (isExpression.value) {
    props.element.parameter["variable_value"] = "1:" + data;
  } else {
    props.element.parameter["variable_value"] = "0:" + data;
  }
}
function changeValueType() {
  isExpression.value = !isExpression.value;
  if (isExpression.value) {
    props.element.parameter["variable_value"] = "1:" + variableValue.value;
  } else {
    props.element.parameter["variable_value"] = "0:" + variableValue.value;
  }
}
function updateData(data) {
  for (const key in data) {
    if (Object.hasOwnProperty.call(data, key)) {
      props.element[key] = data[key];
    }
  }
}
</script>
