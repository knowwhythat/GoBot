<template>
  <InputGroup>
    <InputText
      placeholder="点击右侧按钮选择路径"
      :model-value="modelValue"
      @update:model-value="emit('update:modelValue', $event)"
    />
    <Button icon="pi pi-folder-open" @click="openDialog"></Button>
  </InputGroup>
</template>
<script setup>
import InputGroup from "primevue/inputgroup";
import InputText from "primevue/inputtext";
import Button from "primevue/button";
import { OpenDialog } from "@back/go/backend/App";

const props = defineProps({
  modelValue: {
    type: String,
    default: "",
  },
  type: {
    type: String,
    default: "file",
  },
});

const emit = defineEmits(["update:modelValue"]);

async function openDialog() {
  let path = await OpenDialog({ type: props.type, display: "", pattern: "" });
  if (path) {
    emit("update:modelValue", path);
  }
}
</script>
