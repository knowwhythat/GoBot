<template>
  <ActivityBase
    :runnable="true"
    :toggleable="false"
    :collapsed="true"
    :id="props.element.id"
    :breakpoint="props.element.breakpoint"
    :deleteable="props.element.deleteable"
    :commentable="true"
    :label="props.element.label"
    :icon="props.element.icon_path"
    :color="props.element.color"
    :parameter_define="props.element.parameter_define"
    :parameter="props.element.parameter"
    @delete="$emit('delete', { id: props.element.id })"
    @comment="$emit('comment', cloneDeep(props.element))"
    @uncomment="$emit('uncomment', cloneDeep(props.element))"
    @update="updateData($event)"
    @run="runActivity"
  />
</template>
<script setup>
import ActivityBase from "./ActivityBase.vue";
import { inject } from "vue";
import { RunActivity } from "@back/go/backend/App";
import { EventsOff, EventsOn } from "@back/runtime/runtime.js";
import { useToast } from "primevue/usetoast";
import { cloneDeep } from "lodash";

const toast = useToast();
const props = defineProps({
  element: {
    type: Object,
    default: {},
  },
});
const emit = defineEmits(["delete", "comment", "uncomment"]);

function updateData(data) {
  for (const key in data) {
    if (Object.hasOwnProperty.call(data, key)) {
      props.element[key] = data[key];
    }
  }
}
const { projectId } = inject("projectId");

async function runActivity() {
  EventsOn("repl", (data, out) => {
    if (out) {
      toast.add({
        severity: "error",
        summary: "执行失败",
        detail: out,
        life: 3000,
      });
    } else {
      toast.add({
        severity: "success",
        summary: "执行成功",
        life: 3000,
      });
    }
    EventsOff("repl");
  });
  await RunActivity(projectId, JSON.stringify(props.element));
}
</script>
