<template>
  <InputGroup>
    <Button
      icon="pi pi-chevron-right"
      severity="Primary"
      @click="toggle"
      aria-haspopup="true"
      aria-controls="overlay_panel"
    />
    <OverlayPanel ref="op" appendTo="body" :unstyled="true" :dismissable="true">
      <div class="w-[740px] flex flex-row flex-wrap gap-2 p-2 bg-white">
        <div
          class="flex h-32 w-32 items-center place-content-center bg-slate-100 select-none"
          v-for="image in imageElement"
          :key="image.id"
          @dblclick.stop="showImage(image)"
          @click.stop="selectedNode(image)"
        >
          <img
            :src="'data:image/png;base64,' + image.image"
            alt="Image"
            class="hover:scale-105"
          />
        </div>
      </div>
    </OverlayPanel>
    <InputText :placeholder="input.placeholder" v-model="content" disabled />
    <Button
      icon="pi pi-image"
      severity="secondary"
      @click="startCapture"
      aria-haspopup="true"
      aria-controls="overlay_panel"
      v-tooltip.top="{
        value: '截图',
        showDelay: 100,
        hideDelay: 300,
      }"
    />
    <Dialog
      v-model:visible="visible"
      modal
      header="图片预览"
      :style="{ width: '60rem' }"
    >
      <div class="flex w-full mb-2 gap-2 justify-center">
        <div class="flex flex-col gap-1">
          <span class="mt-3 w-24">图片名称:</span>
          <InputText class="w-full" v-model="imageData.name" type="text" />
        </div>
        <div class="flex flex-col gap-1">
          <span class="mt-3 w-24">匹配精确的:</span>
          <InputNumber class="w-full" v-model="imageData.threshold" />
        </div>
        <div class="flex flex-col gap-1">
          <span class="mt-3 w-24">匹配算法:</span>
          <Dropdown
            v-model="imageData.match_type"
            :options="matchTypes"
            optionLabel="name"
            optionValue="value"
            placeholder="选择匹配方式"
            class="w-full md:w-14rem"
          />
        </div>
      </div>
      <div class="w-full flex justify-center">
        <Image
          :src="'data:image/png;base64,' + imageData.image"
          alt="Image"
          preview
          image-class="max-h-96"
        />
      </div>
      <template #footer>
        <Button label="取消" icon="pi pi-times" text @click="visible = false" />
        <Button label="高亮" icon="pi pi-expand" @click="highlightImage" />
        <Button label="确认" icon="pi pi-check" @click="saveImage" />
      </template>
    </Dialog>
  </InputGroup>
</template>
<script setup>
import InputGroup from "primevue/inputgroup";
import InputText from "primevue/inputtext";
import OverlayPanel from "primevue/overlaypanel";
import Button from "primevue/button";
import { computed, inject, nextTick, ref } from "vue";
import { HighlightImage, SaveImage, StartCapture } from "@back/go/backend/App";
import { WindowMaximise, WindowMinimise } from "@back/runtime/runtime";
import { nanoid } from "nanoid";
import { matchTypes } from "@/utils/shared.js";
import Image from "primevue/image";
import InputNumber from "primevue/inputnumber";
import Dialog from "primevue/dialog";
import Dropdown from "primevue/dropdown";
import { useToast } from "primevue/usetoast";

const toast = useToast();
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

const { id, imageElement } = inject("imageElement");

const content = computed(() => {
  if (props.value.split(":").length > 1) {
    const key = props.value.substring(2);
    if (key) {
      const image = imageElement.value.find((item) => item.id === key);
      if (image) {
        return image["name"];
      }

      return "";
    } else {
      return props.value;
    }
  }
});
const op = ref();
const toggle = (event) => {
  op.value.toggle(event);
};

function selectedNode(node) {
  op.value.toggle(false);
  emit("update", "0:" + node.id);
}

const visible = ref(false);
const imageData = ref(null);
const { projectId } = inject("projectId");
async function startCapture() {
  try {
    await WindowMinimise();
    const resp = await StartCapture();
    if (resp) {
      const data = JSON.parse(resp);
      data.id = nanoid(16);
      data.name = "图片" + imageElement.value.length;
      data.threshold = 0.9;
      data.match_type = "template";
      imageData.value = data;
      visible.value = true;
    }
  } catch (err) {
    toast.add({
      severity: "error",
      summary: "截图失败",
      detail: err,
      life: 3000,
    });
  }
  await WindowMaximise();
}

async function saveImage() {
  try {
    await SaveImage(projectId, imageData.value);
    visible.value = false;
    await nextTick(() => {
      emit("update", "0:" + imageData.value.id);
    });
  } catch (err) {
    toast.add({
      severity: "error",
      summary: "保存失败",
      detail: err,
      life: 3000,
    });
  }
}

async function highlightImage() {
  try {
    await WindowMinimise();
    await HighlightImage(projectId, imageData.value);
    toast.add({
      severity: "success",
      summary: "找到图片",
      life: 3000,
    });
  } catch (err) {
    toast.add({
      severity: "error",
      summary: "高亮失败",
      detail: err,
      life: 3000,
    });
  }
  await WindowMaximise();
}

function showImage(image) {
  imageData.value = image;
  visible.value = true;
}
</script>
