<template>
  <BasicActivity :element="element"></BasicActivity>
</template>
<script setup>
import { onMounted } from "vue";
import BasicActivity from "./activity/BasicActivity.vue";
const props = defineProps({
  component: {
    type: String,
    default: "",
  },
  element: {
    type: Object,
    default: {},
  },
});
const blockComponents = import.meta.glob("@/components/activity/*.vue", {
  eager: true,
});
onMounted(() => {
  console.log(props);
});
</script>
