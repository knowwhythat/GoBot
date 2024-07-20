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
    @update="updateData($event)"
    @run="runActivity"
  >
    <template #top>
      <Image
        class="flex justify-center items-center"
        v-if="imagePath"
        :src="imagePath"
        alt="Image"
        image-class="max-h-8"
        :pt="{
          root: (options) => ({
            style: 'max-width:500px;max-height:32px',
          }),
        }"
        preview
      />
    </template>
  </ActivityBase>
</template>
<script setup>
import Image from "primevue/image";
import { RunActivity, GetElementImage } from "@back/go/backend/App";
import { inject, onMounted, ref } from "vue";
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

const { projectId } = inject("projectId");
const imagePath = ref("");
onMounted(async () => {
  if ("windows_element" in props.element.parameter) {
    const image = await GetElementImage(
      projectId,
      props.element.parameter["windows_element"].substring(2),
    );
    imagePath.value = "data:image/png;base64," + image;
  }
});

function updateData(data) {
  for (const key in data) {
    if (Object.hasOwnProperty.call(data, key)) {
      props.element[key] = data[key];
      if (key === "parameter") {
        console.log(data[key]);
        if (data[key]["windows_element"]) {
          GetElementImage(
            projectId,
            data[key]["windows_element"].substring(2),
          ).then((data) => {
            imagePath.value = "data:image/png;base64," + data;
          });
        } else {
          imagePath.value = null;
        }
      }
    }
  }
}

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
