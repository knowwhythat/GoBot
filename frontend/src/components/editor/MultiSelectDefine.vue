<template>
  <InputGroup>
    <MultiSelect
      v-model="content"
      display="chip"
      :options="props.input.options"
      optionLabel="label"
      optionValue="value"
      :filter="props.input.options.length > 8"
    ></MultiSelect>
  </InputGroup>
</template>
<script setup>
import InputGroup from "primevue/inputgroup";
import MultiSelect from "primevue/multiselect";
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
      if (props.value.substring(2)) {
        return props.value.substring(2).split(";");
      } else {
        return [];
      }
    } else {
      if (props.value) {
        return props.value.split(";");
      } else {
        return [];
      }
    }
  },
  set(value) {
    emit("update", "0:" + value.join(";"));
  },
});
</script>
