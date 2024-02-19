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
    <template #operate>
      <button
        class="p-panel-header-icon p-link mr-2"
        @click.stop="addCondition"
        v-tooltip="{ value: '添加分支', showDelay: 100, hideDelay: 300 }"
      >
        <span class="pi pi-plus"></span>
      </button>
    </template>
    <div class="flex flex-col">
      <SequenceActivity
        v-for="element in props.element.children"
        :element="element"
        @update="update"
        @delete="deleteNode($event)"
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
      key: "if",
      method: "if",
      label: "IF",
      toggleable: true,
      deleteable: false,
      icon_path: "riHome5Line",
      isPseudo: true,
      parameter_define: {
        inputs: [
          {
            category: null,
            key: "expression",
            label: "表达式",
            type: "any",
            default_value: "1:True",
            required: true,
            options: null,
            editor_type: "text",
            show_if: [],
          },
        ],
        outputs: [],
      },
      children: [],
    });
    props.element.children.push({
      id: nanoid(16),
      key: "else",
      method: "else",
      label: "Else",
      toggleable: true,
      deleteable: false,
      icon_path: "riHome5Line",
      isPseudo: true,
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
function addCondition() {
  props.element.children.splice(1, 0, {
    id: nanoid(16),
    key: "else if",
    method: "else if",
    label: "ElseIf",
    toggleable: true,
    deleteable: true,
    icon_path: "riHome5Line",
    isPseudo: true,
    parameter_define: {
      inputs: [
        {
          category: null,
          key: "expression",
          label: "表达式",
          type: "any",
          default_value: "1:True",
          required: true,
          options: null,
          editor_type: "text",
          show_if: [],
        },
      ],
      outputs: [],
    },
    children: [],
  });
}
function updateData(data) {
  for (const key in data) {
    if (Object.hasOwnProperty.call(data, key)) {
      props.element[key] = data[key];
    }
  }
}

function deleteNode({ id }) {
  const blockIndex = props.element.children.findIndex(
    (activity) => activity.id === id
  );

  if (blockIndex !== -1) {
    props.element.children.splice(blockIndex, 1);
  }
}
</script>
