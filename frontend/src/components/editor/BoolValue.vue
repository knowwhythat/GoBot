<template>
  <Checkbox
    class="pt-3"
    v-model="content"
    trueValue="True"
    falseValue="False"
    :binary="true"
  ></Checkbox>
</template>
<script setup>
import Checkbox from "primevue/checkbox";
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
    emit("update", "1:" + value);
  },
});
</script>
