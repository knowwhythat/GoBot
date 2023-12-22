<template>
  <InputGroup>
    <Dropdown
      v-model="content"
      :options="props.input.options"
      optionLabel="label"
      optionValue="value"
      :filter="props.input.options.length > 8"
    ></Dropdown>
  </InputGroup>
</template>
<script setup>
import InputGroup from "primevue/inputgroup";
import Dropdown from "primevue/dropdown";
import { computed } from "vue";
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
const emit = defineEmits(["update"]);

const content = computed({
  get() {
    if (props.value.split(":").length > 1) {
      return props.value.substring(2);
    } else {
      return props.value;
    }
  },
  set(value) {
    emit("update", "0:" + value);
  },
});
</script>
