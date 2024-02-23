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
    @drop="handleDrop"
  >
    <div
      v-for="(element, index) in props.element.children"
      :key="element.id"
      :index="index"
      class="flex flex-col items-center px-2 overflow-auto text-sm m-0.5"
      @dragover.prevent
      @dragenter.prevent="dragEnter($event, element)"
    >
      <div
        v-if="dropBlockId == element.id"
        class="relative w-full h-4 flex items-center"
      >
        <div
          class="absolute top-0 -left-1"
          style="
            border: 8px solid transparent;
            border-right: none;
            border-left-color: rgb(147 51 234);
          "
        ></div>
        <div class="relative w-full bg-purple-600 h-0.5"></div>
        <div
          class="absolute top-0 -right-1"
          style="
            border: 8px solid transparent;
            border-left: none;
            border-right-color: rgb(147 51 234);
          "
        ></div>
      </div>
      <SequenceActivity
        :index="index"
        v-if="element.component == 'SequenceActivity'"
        v-bind="{ element: element }"
        @update="update"
        @delete="deleteNode($event)"
        draggable="true"
        @dragstart="dragStart($event, element)"
        @dragend="dragEnd($event, element)"
      />
      <component
        v-else
        :index="index"
        :is="findComponent(element.component)"
        v-bind="{ element: element }"
        @delete="deleteNode($event)"
        draggable="true"
        @dragstart="dragStart($event, element)"
        @dragend="dragEnd($event, element)"
      />
    </div>
    <div v-if="dropBlockId == id" class="relative w-full h-4 flex items-center">
      <div
        class="absolute top-0 -left-1"
        style="
          border: 8px solid transparent;
          border-right: none;
          border-left-color: rgb(147 51 234);
        "
      ></div>
      <div class="relative w-full bg-purple-600 h-0.5"></div>
      <div
        class="absolute top-0 -right-1"
        style="
          border: 8px solid transparent;
          border-left: none;
          border-right-color: rgb(147 51 234);
        "
      ></div>
    </div>
    <div
      class="w-full p-2 rounded-lg text-gray-600 dark:text-gray-200 border-2 text-center border-dashed"
      @dragover.prevent
      @dragenter.prevent="dragEnter($event, { id })"
    >
      拖放组件到此处
    </div>
  </ActivityBase>
</template>
<script setup>
import { ref, onMounted, computed, inject, provide, watch } from "vue";
import ActivityBase from "./ActivityBase.vue";
import { nanoid } from "nanoid";
import { typeBuilder } from "@/utils/shared";

const props = defineProps({
  element: {
    type: Object,
    default: {},
  },
});

const id = ref(null);
onMounted(() => {
  if (!props.element.children) {
    props.element.children = [];
  }
  id.value = nanoid(16);
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

const { selectedActivity } = inject("selectedActivity");
function deleteNode({ id }) {
  const blockIndex = props.element.children.findIndex(
    (activity) => activity.id === id
  );

  if (blockIndex !== -1) {
    const copyActivities = [...props.element.children];
    copyActivities.splice(blockIndex, 1);
    emit("update", { id: props.element.id, children: copyActivities });
  }
  const index = selectedActivity.value.findIndex((item) => item === id);
  selectedActivity.value.splice(index, 1);
}

function update({ id, children }) {
  if (id) {
    let target = props.element.children.filter((element) => element.id === id);
    target[0].children = children;
  } else {
    props.element.children = children;
  }
}

const { dragBlockId } = inject("dragBlockId");
const { dropBlockId } = inject("dropBlockId");

watch(dragBlockId, (value, oldValue) => {
  if (oldValue && !value) {
    const copyActivities = [...props.element.children];
    const removeIndex = copyActivities.findIndex((item) => item.id == oldValue);
    if (removeIndex > -1) {
      copyActivities.splice(removeIndex, 1);
      emit("update", { id: props.element.id, children: copyActivities });
    }
  }
});

function dragStart(event, element) {
  event.stopPropagation();
  dragBlockId.value = element.id;
  event.dataTransfer.setData("activity", JSON.stringify(element));
}

function dragEnd(event, element) {
  dropBlockId.value = "";
  // dragBlockId.value = "";
}

function dragEnter(event, element) {
  event.stopPropagation();
  dropBlockId.value = element.id;
}

function handleDrop(event) {
  event.preventDefault();
  event.stopPropagation();
  const droppedBlock = JSON.parse(
    event.dataTransfer.getData("activity") || null
  );
  if (!droppedBlock || dropBlockId.value == "") return;

  const copyActivities = [...props.element.children];
  if (droppedBlock.id) {
    console.log(dragBlockId.value);
    dragBlockId.value = "";
    if (id.value == dropBlockId.value) {
      copyActivities.push({ ...droppedBlock, id: nanoid(16), parameter: {} });
    } else {
      const dropIndex = copyActivities.findIndex(
        (item) => item.id == dropBlockId.value
      );
      if (dropIndex > -1) {
        copyActivities.splice(dropIndex, 0, {
          ...droppedBlock,
          id: nanoid(16),
          parameter: {},
        });
      }
    }
  } else {
    if (id.value == dropBlockId.value) {
      copyActivities.push({ ...droppedBlock, id: nanoid(16), parameter: {} });
    } else {
      const dropIndex = copyActivities.findIndex(
        (item) => item.id == dropBlockId.value
      );
      if (dropIndex > -1) {
        copyActivities.splice(dropIndex, 0, {
          ...droppedBlock,
          id: nanoid(16),
          parameter: {},
        });
      }
    }
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
const { contextVariable } = inject("contextVariable");
const currentContextVariable = computed({
  get() {
    let child = props.element.children;
    let variables = [];
    child.forEach((element) => {
      element.parameter_define.outputs?.forEach((pd) => {
        variables.push({
          label:
            pd.key in element.parameter
              ? element.parameter[pd.key]
              : pd.default_value,
          type: pd.type,
          icon: typeBuilder[pd.type] ? typeBuilder[pd.type] : "pi pi-building",
          key: nanoid(8),
        });
      });
    });
    if (contextVariable.value) {
      contextVariable.value.forEach((variable) => {
        if (
          !variables.some((v) => {
            if (v.label === variable.label) {
              return true;
            } else {
              return false;
            }
          })
        ) {
          variables.push(variable);
        }
      });
    }
    return variables;
  },
});
provide("contextVariable", { contextVariable: currentContextVariable });
</script>
