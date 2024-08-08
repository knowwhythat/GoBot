<template>
  <div class="w-full h-32 border border-slate-300 border-solid rounded">
    <MonacoEditor v-model="content"></MonacoEditor>
  </div>
</template>
<script setup>
import MonacoEditor from "@/components/MonacoEditor.vue";
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
