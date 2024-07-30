<template>
  <ActivityBase
    :id="props.element.id"
    :collapsed="props.element.collapsed"
    :toggleable="props.element.toggleable"
    :deleteable="props.element.deleteable"
    :commentable="props.element.key !== 'tools.comment'"
    :showComment="props.element.showComment"
    :label="props.element.label"
    :icon="props.element.icon_path"
    :color="props.element.color"
    :parameter_define="props.element.parameter_define"
    :parameter="props.element.parameter"
    @delete="$emit('delete', { id: props.element.id })"
    @update="updateData($event)"
    @comment="$emit('comment', cloneDeep(props.element))"
    @uncomment="$emit('uncomment', cloneDeep(props.element))"
    @drop="handleDrop"
  >
    <div
      v-for="(element, index) in props.element.children"
      :key="element.id"
      class="flex flex-col items-center px-2 overflow-auto text-sm m-0.5"
      @dragover.prevent
      @dragenter.prevent="dragEnter($event, element)"
    >
      <div
        v-if="dropBlockId === element.id"
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
        v-if="element.component === 'SequenceActivity'"
        v-bind="{ element: element }"
        @update="update"
        @delete="deleteNode($event)"
        @comment="onComment"
        @uncomment="onUncomment"
        draggable="true"
        @dragstart="dragStart($event, element)"
        @dragend="dragEnd($event, element)"
      />
      <component
        v-else
        :is="findComponent(element.component)"
        v-bind="{ element: element }"
        @delete="deleteNode($event)"
        @comment="onComment"
        @uncomment="onUncomment"
        draggable="true"
        @dragstart="dragStart($event, element)"
        @dragend="dragEnd($event, element)"
      />
    </div>
    <div
      v-if="dropBlockId === id"
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
import {
  computed,
  inject,
  nextTick,
  onMounted,
  provide,
  ref,
  watch,
} from "vue";
import ActivityBase from "./ActivityBase.vue";
import { nanoid } from "nanoid";
import { typeBuilder } from "@/utils/shared";
import { cloneDeep } from "lodash";
import { getVariableDefine } from "@/utils/helper.js";

const props = defineProps({
  element: {
    type: Object,
    default: {},
  },
});

const id = ref("");
onMounted(() => {
  if (!props.element.children) {
    props.element.children = [];
  }
  id.value = nanoid(16);
});

const emit = defineEmits(["update", "delete", "comment", "uncomment"]);

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
    (activity) => activity.id === id,
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

function onComment(data) {
  const blockIndex = props.element.children.findIndex(
    (activity) => activity.id === data.id,
  );
  if (blockIndex !== -1) {
    const commentBlock = {
      key: "tools.comment",
      label: "注释",
      icon_path: "riFileForbidFill",
      description: "注释掉其中的组件",
      isPseudo: true,
      parameter_define: {},
      component: "SequenceActivity",
      id: nanoid(16),
      parameter: {},
      children: [data],
    };
    const copyActivities = [...props.element.children];
    copyActivities.splice(blockIndex, 1, commentBlock);
    emit("update", { id: props.element.id, children: copyActivities });
  }
}

function onUncomment(data) {
  if (data.key === "tools.comment") {
    const blockIndex = props.element.children.findIndex(
      (activity) => activity.id === data.id,
    );
    if (blockIndex !== -1) {
      const copyActivities = [...props.element.children];
      copyActivities.splice(blockIndex, 1, ...data.children);
      emit("update", { id: props.element.id, children: copyActivities });
    }
  }
}

const { dragBlockId } = inject("dragBlockId");
const { dropBlockId } = inject("dropBlockId");

watch(dragBlockId, (value, oldValue) => {
  if (oldValue && !value) {
    const copyActivities = [...props.element.children];
    const removeIndex = copyActivities.findIndex(
      (item) => item.id === oldValue,
    );
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
}

function dragEnter(event, element) {
  event.stopPropagation();
  dropBlockId.value = element.id;
}

const newAddBlockId = ref(null);
provide("newAddBlockId", { newAddBlockId: newAddBlockId });

function handleDrop(event) {
  event.preventDefault();
  const droppedBlock = JSON.parse(
    event.dataTransfer.getData("activity") || null,
  );
  if (!droppedBlock || dropBlockId.value === "") return;

  const copyActivities = [...props.element.children];
  const newBlockId = nanoid(16);
  if (droppedBlock.id) {
    if (id.value === dropBlockId.value) {
      copyActivities.push({ ...droppedBlock, id: newBlockId });
      dragBlockId.value = "";
    } else {
      const dropIndex = copyActivities.findIndex(
        (item) => item.id === dropBlockId.value,
      );
      if (dropIndex > -1) {
        copyActivities.splice(dropIndex, 0, {
          ...droppedBlock,
          id: newBlockId,
        });
        dragBlockId.value = "";
      } else {
        return;
      }
    }
  } else {
    if (id.value === dropBlockId.value) {
      copyActivities.push({ ...droppedBlock, id: newBlockId, parameter: {} });
    } else {
      const dropIndex = copyActivities.findIndex(
        (item) => item.id === dropBlockId.value,
      );
      if (dropIndex > -1) {
        copyActivities.splice(dropIndex, 0, {
          ...droppedBlock,
          id: newBlockId,
          parameter: {},
        });
      } else {
        return;
      }
    }
  }
  selectedActivity.value = [newBlockId];
  emit("update", { id: props.element.id, children: copyActivities });
  if (!droppedBlock.id) {
    nextTick(() => {
      newAddBlockId.value = newBlockId;
    });
  }
  event.stopPropagation();
}

function updateData(data) {
  for (const key in data) {
    if (Object.hasOwnProperty.call(data, key)) {
      props.element[key] = data[key];
    }
  }
}

const { contextVariable } = inject("contextVariable");
const currentContextVariable = computed(() => {
  let child = props.element.children;
  let variables = [];
  child.forEach((element) => {
    element.parameter_define.outputs?.forEach((pd) => {
      const variableLabel =
        pd.key in element.parameter
          ? element.parameter[pd.key]
          : pd.default_value;
      if (
        variableLabel &&
        variables.findIndex((v) => v.label === variableLabel) === -1
      ) {
        variables.push(getVariableDefine(variableLabel, pd.type));
      }
    });
  });
  if (contextVariable.value) {
    contextVariable.value.forEach((variable) => {
      if (
        !variables.some((v) => {
          return v.label === variable.label;
        })
      ) {
        variables.push(variable);
      }
    });
  }
  return variables;
});
provide("contextVariable", { contextVariable: currentContextVariable });
</script>
