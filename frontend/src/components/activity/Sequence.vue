<template>
  <div
    class="flex flex-col min-w-[300px] bg-gray-50 border-2 border-gray-200 rounded-xl hover:ring-2 m-4"
    style="cursor: grab"
  >
    <div v-if="props.showLabel" class="bg-primary rounded-t-xl p-2 truncate">
      {{ props.label }}
    </div>
    <draggable
      :model-value="props.children"
      item-key="id"
      class="flex flex-col items-center p-x-4 mb-8 overflow-auto nowheel scroll text-sm p-2"
      @mousedown.stop
      @dragover.prevent
      @drop="handleDrop"
      @update:modelValue="$emit('update', { children: $event })"
    >
      <template #item="{ element, index }">
        <BaseActivity
          v-bind="{ element: element, component: element.component }"
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
import BaseActivity from "@/components/BaseActivity.vue";

const props = defineProps({
  id: {
    type: String,
    default: nanoid(8),
  },
  label: {
    type: String,
    default: "子模块",
  },
  showLabel: {
    type: Boolean,
    default: false,
  },
  children: {
    type: Array,
    default: () => [],
  },
});

const emit = defineEmits(["update", "delete"]);

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
  console.log(event);
  setTimeout(() => {
    const blockEl = document.querySelector(`[group-item-id="${item.itemId}"]`);

    if (blockEl) {
      const blockIndex = props.children.findIndex(
        (activity) => activity.itemId === item.itemId
      );

      if (blockIndex !== -1) {
        const copyActivities = [...props.children];
        copyActivities.splice(blockIndex, 1);
        emit("update", { children: copyActivities });
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

  const copyActivities = [...props.children];

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

  emit("update", { children: copyActivities });
}
</script>
