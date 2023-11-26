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
    <draggable
      v-if="props.element.children"
      :model-value="props.element.children"
      item-key="id"
      group="sequence"
      class="flex flex-col items-center p-x-4 mb-8 overflow-auto nowheel scroll text-sm p-2"
      @mousedown.stop
      @dragover.prevent
      @drop="handleDrop"
      @update:modelValue="
        $emit('update', { id: props.element.id, children: $event })
      "
    >
      <template #item="{ element, index }">
        <SequenceActivity
          :index="index"
          v-if="element.component == 'SequenceActivity'"
          v-bind="{ element: element }"
          @update="update"
          @delete="deleteNode($event)"
        />
        <component
          v-else
          :index="index"
          :is="findComponent(element.component)"
          v-bind="{ element: element }"
          @delete="deleteNode($event)"
        />
      </template>
      <template #footer>
        <div
          class="w-full p-8 rounded-lg text-gray-600 dark:text-gray-200 border-2 text-center border-dashed"
        >
          拖放组件到此处
        </div>
      </template>
    </draggable>
  </ActivityBase>
</template>
<script setup>
import { ref, onMounted, computed, shallowReactive, watch } from "vue";
import ActivityBase from "./ActivityBase.vue";
import draggable from "vuedraggable";
import { nanoid } from "nanoid";

const props = defineProps({
  element: {
    type: Object,
    default: {},
  },
});

onMounted(() => {
  if (!props.element.children) {
    props.element.children = [];
  }
});

const emit = defineEmits(["update", "delete"]);

const activityComponents = import.meta.glob("@/components/activity/*.vue", {
  eager: true,
});
const activityMap = {};
for (let each in activityComponents) {
  const name = activityComponents[each].default.__name;
  activityMap[`${name}`] = activityComponents[each].default;
}

function findComponent(name) {
  return activityMap[name];
}

function deleteNode({ id }) {
  console.log(id);
  const blockIndex = props.element.children.findIndex(
    (activity) => activity.id === id
  );

  if (blockIndex !== -1) {
    const copyActivities = [...props.element.children];
    copyActivities.splice(blockIndex, 1);
    emit("update", { id: props.element.id, children: copyActivities });
  }
}

function update({ id, children }) {
  if (id) {
    console.log(id);
    let target = props.element.children.filter((element) => element.id === id);
    target[0].children = children;
  } else {
    props.element.children = children;
  }
}

function handleDrop(event) {
  console.log(event);
  event.preventDefault();
  event.stopPropagation();
  const droppedBlock = JSON.parse(
    event.dataTransfer.getData("activity") || null
  );
  if (!droppedBlock || droppedBlock.id) return;

  const copyActivities = [...props.element.children];

  const ancestorElement = event.target.closest(".activity-node");
  if (ancestorElement) {
    let index = ancestorElement.getAttribute("index");
    copyActivities.splice(
      index,
      0,
      shallowReactive({ ...droppedBlock, id: nanoid(16) })
    );
  } else {
    copyActivities.push(shallowReactive({ ...droppedBlock, id: nanoid(16) }));
  }

  emit("update", { id: props.element.id, children: copyActivities });
}
function updateData(data) {
  for (const key in data) {
    if (Object.hasOwnProperty.call(data, key)) {
      props.element[key] = data[key];
    }
  }
}
</script>