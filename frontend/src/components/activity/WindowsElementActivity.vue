<template>
  <ActivityBase
    :runnable="true"
    :toggleable="true"
    :id="props.element.id"
    :breakpoint="props.element.breakpoint"
    :collapsed="props.element.collapsed"
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
    <div class="max-h-32 flex justify-center mb-2">
      <Image
        class="my-auto"
        v-if="imagePath"
        :src="imagePath"
        alt="Image"
        imageStyle="max-width:350px;max-height:100px"
        preview
      />
      <div class="text-center my-auto" v-else>元素图片</div>
    </div>
    <div class="flex gap-4 justify-center">
      <SplitButton
        severity="Primary"
        label="拾取"
        @click="pickElement"
        :model="pickItems"
      >
        <template #icon>
          <v-remixicon name="riCursorLine" size="24" />
        </template>
      </SplitButton>
      <Button
        severity="info"
        label="编辑"
        @click="editElement"
        :disabled="!imagePath"
      >
        <template #icon>
          <v-remixicon name="riEditBoxLine" size="24" />
        </template>
      </Button>
      <Button
        icon="pi pi-arrow-right-arrow-left"
        severity="help"
        :disabled="!imagePath"
        label="校验"
        @click="checkElement()"
      />
      <ElementOptionDialog
        :dialogShow="dialogShow"
        :pathOption="pathOption"
        :needInit="needInit"
        @hide="dialogShow = false"
        @update="saveElement($event)"
      />
      <ElementTreeDialog
        :dialogShow="treeDialogShow"
        @hide="treeDialogShow = false"
        @confirm="selectElement($event)"
      />
    </div>
  </ActivityBase>
</template>
<script setup>
import Image from "primevue/image";
import SplitButton from "primevue/splitbutton";
import {
  StartPickWindowsElement,
  StartCheckWindowsElement,
  RunActivity,
  SaveWindowsElement,
  GetWindowsElement,
  GetElementImage,
  GetSelectedWindowsElement,
} from "@back/go/main/App";
import { inject, onMounted, ref } from "vue";
import Button from "primevue/button";
import ActivityBase from "./ActivityBase.vue";
import { EventsOn, EventsOff } from "@back/runtime/runtime.js";
import ElementOptionDialog from "../element/ElementOptionDialog.vue";
import ElementTreeDialog from "../element/ElementTreeDialog.vue";
import { useToast } from "primevue/usetoast";
const toast = useToast();
const props = defineProps({
  element: {
    type: Object,
    default: {},
  },
});
const emit = defineEmits(["delete"]);

const { projectId } = inject("projectId");
const imagePath = ref("");
onMounted(async () => {
  if ("windows_element" in props.element.parameter) {
    const image = await GetElementImage(
      projectId,
      props.element.parameter["windows_element"].substring(2)
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
            data[key]["windows_element"].substring(2)
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

const treeDialogShow = ref(false);

const pickItems = [
  {
    label: "深度拾取",
    icon: "pi pi-filter-slash",
    command: () => {
      treeDialogShow.value = true;
    },
  },
];

async function selectElement(id) {
  try {
    const resp = await GetSelectedWindowsElement(projectId, id);
    const result = JSON.parse(resp);
    if (result.result === "ok") {
      needInit.value = true;
      pathOption.value = JSON.parse(result.data);
      treeDialogShow.value = false;
      dialogShow.value = true;
    } else {
      toast.add({
        severity: "error",
        summary: "拾取失败",
        detail: result.reason,
        life: 3000,
      });
    }
  } catch (err) {
    toast.add({
      severity: "error",
      summary: "拾取失败",
      detail: err,
      life: 3000,
    });
  }
}

const dialogShow = ref(false);
const pathOption = ref({});
const needInit = ref(false);
async function pickElement() {
  try {
    const resp = await StartPickWindowsElement();
    const result = JSON.parse(resp);
    if (result.result === "ok") {
      needInit.value = true;
      pathOption.value = JSON.parse(result.data);
      dialogShow.value = true;
    } else {
      toast.add({
        severity: "error",
        summary: "拾取失败",
        detail: result.reason,
        life: 3000,
      });
    }
  } catch (err) {
    toast.add({
      severity: "error",
      summary: "拾取失败",
      detail: err,
      life: 3000,
    });
  }
}

async function editElement() {
  const paths = await GetWindowsElement(
    projectId,
    props.element.parameter["windows_element"].substring(2)
  );
  pathOption.value = {
    id: props.element.parameter["windows_element"].substring(2),
    paths: JSON.parse(paths),
  };
  needInit.value = false;
  dialogShow.value = true;
}
const { dataChanged, updateDataChanged } = inject("dataChanged");
async function saveElement(paths) {
  console.log(paths);
  await SaveWindowsElement(
    projectId,
    pathOption.value.id,
    JSON.stringify(paths)
  );
  props.element.parameter["windows_element"] = "0:" + pathOption.value.id;
  const image = await GetElementImage(
    projectId,
    props.element.parameter["windows_element"].substring(2)
  );
  imagePath.value = "data:image/png;base64," + image;
  dialogShow.value = false;
  updateDataChanged();
}
async function checkElement(paths) {
  try {
    const paths = await GetWindowsElement(
      projectId,
      props.element.parameter["windows_element"].substring(2)
    );
    const resp = await StartCheckWindowsElement(paths);
    const result = JSON.parse(resp);
    if (result["result"] === "ok") {
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
