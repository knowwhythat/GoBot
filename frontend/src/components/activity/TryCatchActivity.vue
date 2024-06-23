<template>
  <ActivityBase
    :id="props.element.id"
    :collapsed="props.element.collapsed"
    :toggleable="props.element.toggleable"
    :deleteable="props.element.deleteable"
    :label="props.element.label"
    :icon="props.element.icon_path"
    :color="props.element.color"
    :parameter_define="props.element.parameter_define"
    :parameter="props.element.parameter"
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
import { onMounted, reactive } from "vue";
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
      key: "try",
      method: "try",
      isPseudo: true,
      toggleable: true,
      deleteable: false,
      icon_path: "riHome5Line",
      children: [],
    });
    props.element.children.push({
      id: nanoid(16),
      label: "Catch",
      key: "except",
      method: "except",
      isPseudo: true,
      toggleable: true,
      deleteable: false,
      icon_path: "riHome5Line",
      parameter_define: {
        inputs: [],
        outputs: [
          {
            category: null,
            key: "expression",
            label: "异常变量",
            type: "any",
            default_value: "1:exception",
            required: true,
            options: null,
            editor_type: "text",
            show_if: [],
          },
        ],
      },
      children: [],
    });
    props.element.children.push({
      id: nanoid(16),
      label: "Finally",
      key: "finally",
      method: "finally",
      isPseudo: true,
      toggleable: true,
      deleteable: false,
      icon_path: "riHome5Line",
      children: [],
    });
  }
});

function update({ id, children }) {
  if (id) {
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
