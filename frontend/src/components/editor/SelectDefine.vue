<template>
  <InputGroup>
    <Dropdown
      v-model="content"
      @update:modelValue="updateValue"
      :options="props.input.options"
      optionLabel="label"
      optionValue="value"
    ></Dropdown>
  </InputGroup>
</template>
<script setup>
import InputGroup from "primevue/inputgroup";
import Dropdown from "primevue/dropdown";
import { reactive, onMounted, ref, watch } from "vue";
const props = defineProps({
  value: {
    type: String,
    default: "",
  },
  input: {
    type: Object,
    default: () => ({}),
  },
});

const content = ref("");
const emit = defineEmits(["update"]);

watch(
  () => props.value,
  (value, oldValue) => {
    if (value.split(":").length > 1) {
      content.value = value.substring(2);
    } else {
      content.value = value;
    }
    console.log(content.value);
  },
  { immediate: true }
);

function updateValue(data) {
  content.value = data;
  emit("update", "0:" + data);
}
</script>
