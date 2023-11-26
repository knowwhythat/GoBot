<template>
  <ActivityBase
    :collapsed="props.element.collapsed"
    :toggleable="props.element.toggleable"
    :deleteable="props.element.deleteable"
    :label="props.element.label"
    :icon="props.element.icon_path"
    :color="props.element.color"
    @delete="$emit('delete', { id: props.element.id })"
    @update="updateData($event)"
  >
    <div class="flex flex-col">
      <SequenceActivity
        v-for="element in props.element.children"
        :element="element"
        @update="update"
      />
    </div>
  </ActivityBase>
</template>
<script setup>
import {
  ref,
  onMounted,
  computed,
  shallowReactive,
  watch,
  reactive,
} from "vue";
import ActivityBase from "./ActivityBase.vue";
import SequenceActivity from "@/components/activity/SequenceActivity.vue";
import { nanoid } from "nanoid";

const props = defineProps({
  element: {
    type: Object,
    default: {},
  },
});
defineEmits(["delete"]);
onMounted(() => {
  if (!props.element.children || props.element.children.length < 2) {
    props.element.children = reactive([]);
    props.element.children.push({
      id: nanoid(16),
      label: "Try",
      toggleable: true,
      deleteable: false,
      icon_path: "riHome5Line",
      children: [],
    });
    props.element.children.push({
      id: nanoid(16),
      label: "Catch",
      toggleable: true,
      deleteable: false,
      icon_path: "riHome5Line",
      children: [],
    });
    props.element.children.push({
      id: nanoid(16),
      label: "Finally",
      toggleable: true,
      deleteable: false,
      icon_path: "riHome5Line",
      children: [],
    });
  }
});

function update({ id, children }) {
  if (id) {
    console.log(id);
    let target = props.element.children.filter((element) => element.id === id);
    target[0].children = children;
  } else {
    props.element.children = children;
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
