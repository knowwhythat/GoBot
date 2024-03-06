<template>
  <ActivityBase
    :runnable="true"
    :toggleable="false"
    :id="props.element.id"
    :breakpoint="props.element.breakpoint"
    :collapsed="true"
    :deleteable="props.element.deleteable"
    :label="props.element.label"
    :icon="props.element.icon_path"
    :color="props.element.color"
    :parameter_define="props.element.parameter_define"
    :parameter="props.element.parameter"
    @delete="$emit('delete', { id: props.element.id })"
    @update="updateData($event)"
    @run="runActivity"
  >
    <template #top>
      <div class="flex gap-4 justify-center">
        <Button
          class="h-8"
          icon="pi pi-arrow-right-arrow-left"
          severity="help"
          label="校验"
          @click="checkElement"
        />
      </div>
    </template>
  </ActivityBase>
</template>
<script setup>
import { StartCheck, RunActivity } from "@back/go/main/App";
import { inject } from "vue";
import Button from "primevue/button";
import ActivityBase from "./ActivityBase.vue";
import { EventsOn, EventsOff } from "@back/runtime/runtime.js";
import { useToast } from "primevue/usetoast";
const toast = useToast();
const props = defineProps({
  element: {
    type: Object,
    default: {},
  },
});
const emit = defineEmits(["delete"]);

async function checkElement() {
  try {
    let framePath = "";
    if (
      props.element.parameter["frame_selector"] &&
      props.element.parameter["frame_selector"].length > 2
    ) {
      framePath = props.element.parameter["frame_selector"].substring(2);
    }
    let elementPath = "";
    if (
      props.element.parameter["element_selector"] &&
      props.element.parameter["element_selector"].length > 2
    ) {
      elementPath = props.element.parameter["element_selector"].substring(2);
    } else {
      toast.add({
        severity: "error",
        summary: "校验失败",
        detail: "元素选择器不能为空",
        life: 3000,
      });
      return;
    }
    const resp = await StartCheck(framePath, elementPath);
    const result = JSON.parse(resp);
    if (result["message"] === "ok") {
      toast.add({
        severity: "success",
        summary: "校验成功",
        detail: "找到元素",
        life: 3000,
      });
    } else {
      toast.add({
        severity: "error",
        summary: "校验失败",
        detail: "元素未找到",
        life: 3000,
      });
    }
  } catch (err) {
    toast.add({
      severity: "error",
      summary: "校验失败",
      detail: err,
      life: 3000,
    });
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
