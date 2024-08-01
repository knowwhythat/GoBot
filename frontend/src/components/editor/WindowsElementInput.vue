<template>
  <InputGroup>
    <Button
      icon="pi pi-chevron-right"
      severity="Primary"
      @click="toggle"
      aria-haspopup="true"
      aria-controls="overlay_panel"
    />
    <OverlayPanel ref="op" appendTo="body" :unstyled="true">
      <Tree
        scrollHeight="300px"
        :value="windowsElement"
        selectionMode="single"
        :filter="true"
        filterPlaceholder="搜索元素"
        @node-select="selectedNode($event)"
      ></Tree>
    </OverlayPanel>
    <InputText :placeholder="input.placeholder" v-model="content" disabled />
    <Button
      icon="pi pi-bullseye"
      severity="secondary"
      @click="pickElement"
      aria-haspopup="true"
      aria-controls="overlay_panel"
      v-tooltip.top="{
        value: '拾取元素',
        showDelay: 100,
        hideDelay: 300,
      }"
    />
    <Button
      icon="pi pi-align-left"
      severity="info"
      @click="deepPickElement"
      aria-haspopup="true"
      aria-controls="overlay_panel"
      v-tooltip.top="{
        value: '深度拾取',
        showDelay: 100,
        hideDelay: 300,
      }"
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
  </InputGroup>
</template>
<script setup>
import Tree from "primevue/tree";
import InputGroup from "primevue/inputgroup";
import InputText from "primevue/inputtext";
import OverlayPanel from "primevue/overlaypanel";
import Button from "primevue/button";
import { computed, inject, ref } from "vue";
import {
  GetSelectedWindowsElement,
  SaveWindowsElement,
  StartPickWindowsElement,
} from "@back/go/backend/App";
import {
  EventsOnce,
  WindowMaximise,
  WindowMinimise,
} from "@back/runtime/runtime";
import ElementOptionDialog from "@/components/element/ElementOptionDialog.vue";
import ElementTreeDialog from "@/components/element/ElementTreeDialog.vue";

const props = defineProps({
  value: {
    type: String,
    default: "",
  },
  input: {
    type: Object,
    default: () => ({}),
  },
});

const emit = defineEmits(["update"]);

const { id, windowsElement } = inject("windowsElement");

const content = computed({
  get: () => {
    if (props.value.split(":").length > 1) {
      const key = props.value.substring(2);
      if (key) {
        let ele;
        windowsElement.value.forEach((element) => {
          element.children.forEach((item) => {
            if (item.key === key) {
              ele = item;
            }
          });
        });
        if (ele) {
          return ele.label;
        } else {
          return "元素不存在";
        }
      }
      return "";
    } else {
      return props.value;
    }
  },
});

const op = ref();
const toggle = (event) => {
  op.value.toggle(event);
};

function selectedNode(node) {
  if (!node["children"]) {
    op.value.toggle(false);
    emit("update", "0:" + node.key);
  }
}
const { projectId } = inject("projectId");
const treeDialogShow = ref(false);
const dialogShow = ref(false);
const pathOption = ref({});
const needInit = ref(false);

async function pickElement() {
  try {
    await StartPickWindowsElement();
    await WindowMinimise();
    EventsOnce("windowsEvent", async (resp) => {
      await WindowMaximise();
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
    });
  } catch (err) {
    toast.add({
      severity: "error",
      summary: "拾取失败",
      detail: err,
      life: 3000,
    });
  }
}

function deepPickElement() {
  treeDialogShow.value = true;
}

async function saveElement(element) {
  const image = element.image;
  delete element.image;

  await SaveWindowsElement(
    projectId,
    element.id,
    image,
    JSON.stringify(element),
  );
  dialogShow.value = false;
  emit("update", "0:" + element.id);
}

async function selectElement(id) {
  try {
    await GetSelectedWindowsElement(id);
    await WindowMinimise();
    EventsOnce("windowsEvent", async (resp) => {
      await WindowMaximise();
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
    });
  } catch (err) {
    toast.add({
      severity: "error",
      summary: "拾取失败",
      detail: err,
      life: 3000,
    });
  }
}
</script>
