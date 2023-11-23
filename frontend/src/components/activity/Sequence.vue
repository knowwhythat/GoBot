<template>
  <div
    class="flex flex-col min-w-[300px] bg-gray-50 border-2 border-gray-200 rounded-xl m-4"
    style="cursor: grab"
  >
    <div class="flex bg-gray-300 rounded-t-xl p-2 truncate">
      <v-remixicon :name="props.element.icon_path" class="pr-2" />
      {{ props.element.label }}
    </div>
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
        <Sequence
          v-if="element.component == 'Sequence'"
          v-bind="{ element: element }"
          @dragstart="onDragStart(element, $event)"
          @dragend="onDragEnd(element, $event)"
          @update="update"
        />
        <component
          v-else
          :is="findComponent(element.component)"
          v-bind="{ element: element }"
          @dragstart="onDragStart(element, $event)"
          @dragend="onDragEnd(element, $event)"
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
  </div>
</template>
<script setup>
import { ref, onMounted, computed, shallowReactive, watch } from "vue";
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

function update({ id, children }) {
  if (id) {
    console.log(id);
    let target = props.element.children.filter((element) => element.id === id);
    target[0].children = children;
  } else {
    props.element.children = children;
  }
}

function onDragStart(item, event) {
  event.dataTransfer.setData(
    "activity",
    JSON.stringify({ ...item, from: "sequence", containerId: props.id })
  );
}

function onDragEnd(item, event) {
  const droppedBlock = JSON.parse(
    event.dataTransfer.getData("activity") || null
  );
  setTimeout(() => {
    const blockEl = document.querySelector(`[group-item-id="${item.itemId}"]`);

    if (blockEl) {
      const blockIndex = props.children.findIndex(
        (activity) => activity.itemId === item.itemId
      );

      if (blockIndex !== -1) {
        const copyActivities = [...props.children];
        copyActivities.splice(blockIndex, 1);
        emit("update", { id: props.element.id, children: copyActivities });
      }
    }
  }, 200);
}

function handleDrop(event) {
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
      shallowReactive({ ...droppedBlock, id: nanoid(10) })
    );
  } else {
    copyActivities.push(shallowReactive({ ...droppedBlock, id: nanoid(10) }));
  }

  emit("update", { id: props.element.id, children: copyActivities });
}
</script>
