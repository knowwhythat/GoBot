<template>
  <ActivityBase
    :toggleable="false"
    :collapsed="true"
    :deleteable="props.element.deleteable"
    :label="props.element.label"
    :icon="props.element.icon_path"
    :color="props.element.color"
    @delete="$emit('delete', { id: props.element.id })"
    @update="updateData($event)"
  />
</template>
<script setup>
import { onMounted } from "vue";
import ActivityBase from "./ActivityBase.vue";
const props = defineProps({
  element: {
    type: Object,
    default: {},
  },
});
const emit = defineEmits(["delete"]);
onMounted(() => {
  console.log(props);
});
function updateData(data) {
  for (const key in data) {
    if (Object.hasOwnProperty.call(data, key)) {
      props.element[key] = data[key];
    }
  }
}
</script>
