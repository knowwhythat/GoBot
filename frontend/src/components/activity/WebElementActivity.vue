<template>
  <ActivityBase
    :runnable="true"
    :toggleable="false"
    :id="props.element.id"
    :breakpoint="props.element.breakpoint"
    :collapsed="true"
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
    @run="runActivity"
    @update="updateData($event)"
  >
    <template #top>
      <div class="flex gap-4 justify-center">
        <Button
          class="h-8"
          icon="pi pi-search"
          severity="help"
          label="拾取"
          @click="pickElement"
        />
        <Button
          class="h-8"
          icon="pi pi-arrow-right-arrow-left"
          severity="help"
          label="校验"
          @click="checkElement"
        />
      </div>
    </template>
    <Dialog
      v-model:visible="dialogVisible"
      modal
      header="编辑元素选择器"
      :style="{ width: '45rem' }"
    >
      <div class="mt-3 flex justify-content-center">
        <label for="value" class="mt-3 w-32">Frame选择器</label>
        <Dropdown
          v-model="selectedFrame"
          editable
          :options="frames"
          class="w-full md:w-14rem"
        />
      </div>
      <div class="mt-3 flex justify-content-center">
        <label for="value" class="mt-3 w-32">元素选择器</label>
        <Dropdown
          v-model="selectedSelector"
          editable
          :options="selectors"
          class="w-full md:w-14rem"
        />
      </div>
      <template #footer>
        <Button
          label="取消"
          icon="pi pi-times"
          @click="dialogVisible = false"
          text
        />
        <Button label="校验" icon="pi pi-sun" @click="check" />
        <Button label="确认" icon="pi pi-check" @click="updateInnerData" />
      </template>
    </Dialog>
  </ActivityBase>
</template>
<script setup>
import Dialog from "primevue/dialog";
import Dropdown from "primevue/dropdown";
import { StartCheck, StartPick, RunActivity } from "@back/go/backend/App";
import { ref, inject } from "vue";
import Button from "primevue/button";
import ActivityBase from "./ActivityBase.vue";
import { EventsOn, EventsOff } from "@back/runtime/runtime.js";
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
const dialogVisible = ref(false);
const frames = ref([]);
const selectedFrame = ref("");
const selectors = ref([]);
const selectedSelector = ref("");

async function pickElement() {
  const result = await StartPick();
  const resp = JSON.parse(result);
  frames.value = resp["message"]["framePath"];
  selectors.value = resp["message"]["elementPath"];
  if (frames.value.length > 0) {
    selectedFrame.value = frames.value[0];
  }
  if (selectors.value.length > 0) {
    selectedSelector.value = selectors.value[0];
  }
  dialogVisible.value = true;
}
async function check() {
  try {
    const resp = await StartCheck(selectedFrame.value, selectedSelector.value);
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
function updateInnerData() {
  props.element.parameter["frame_selector"] = "0:" + selectedFrame.value;
  props.element.parameter["element_selector"] = "0:" + selectedSelector.value;
  dialogVisible.value = false;
}

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
