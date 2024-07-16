<template>
  <Panel
    header="图片库"
    :pt="{
      header: (options) => ({
        class: ['bg-slate-100', 'p-1'],
      }),
      root: (options) => ({
        style: {
          height: '100%',
          overflow: 'hidden',
        },
      }),
    }"
    :toggleable="false"
  >
    <template #icons>
      <div class="flex align-items-center gap-2">
        <button
          class="mr-4 hover:text-purple-500"
          @click="doStartCapture"
          v-tooltip.left="{
            value: '截图',
            showDelay: 100,
            hideDelay: 300,
          }"
        >
          <v-remixicon name="riScreenshotLine" size="16" />
        </button>
        <button class="hover:text-purple-500" @click="$emit('hide')">
          <span class="pi pi-chevron-down"></span>
        </button>
      </div>
    </template>
    <div class="h-full flex flex-row flex-wrap">
      <Image
        v-for="image in images"
        :key="image.id"
        :src="image.path"
        alt="Image"
        preview
        class="flex"
        image-class="flex justify-center items-center object-scale-down"
      />
    </div>
    <Dialog
      v-model:visible="visible"
      modal
      header="图片预览"
      :style="{ width: '60rem' }"
    >
      <div class="flex mb-1">
        <span class="mt-3 w-24">图片名称:</span>
        <InputText class="w-full" v-model="imageData.name" type="text" />
      </div>
      <div class="w-full flex justify-center">
        <Image :src="'data:image/png;base64,' + imageData.image" alt="Image" />
      </div>
      <template #footer>
        <Button label="取消" icon="pi pi-times" text @click="visible = false" />
        <Button label="高亮" icon="pi pi-expand" @click="highlightImage" />
        <Button label="确认" icon="pi pi-check" @click="saveImage" />
      </template>
    </Dialog>
  </Panel>
</template>
<script setup>
import Image from "primevue/image";
import Panel from "primevue/panel";
import Dialog from "primevue/dialog";

import { HighlightImage, SaveImage, StartCapture } from "@back/go/backend/App";
import { WindowMaximise, WindowMinimise } from "@back/runtime/runtime.js";
import { useToast } from "primevue/usetoast";
import { useConfirm } from "primevue/useconfirm";
import { sleep } from "@/utils/helper.js";
import Button from "primevue/button";
import { ref } from "vue";
import { nanoid } from "nanoid";
import InputText from "primevue/inputtext";

const toast = useToast();
const confirm = useConfirm();
const emit = defineEmits(["hide"]);
const props = defineProps({
  id: {
    type: String,
    default: "",
  },
  images: {
    type: Array,
    default: () => [],
  },
});
const visible = ref(false);
const imageData = ref(null);

async function doStartCapture() {
  try {
    await WindowMinimise();
    await sleep(500);
    const resp = await StartCapture();
    if (resp) {
      const data = JSON.parse(resp);
      data.id = nanoid(16);
      data.name = "图片" + props.images.length;
      imageData.value = data;
      visible.value = true;
    }
    await WindowMaximise();
  } catch (err) {
    toast.add({
      severity: "error",
      summary: "截图失败",
      detail: err,
      life: 3000,
    });
  }
}

async function saveImage() {
  try {
    await SaveImage(props.id, imageData.value);
    visible.value = false;
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
    await HighlightImage(props.id, imageData.value.id);
  } catch (err) {
    toast.add({
      severity: "error",
      summary: "高亮失败",
      detail: err,
      life: 3000,
    });
  }
}
</script>
<style scoped>
:deep(.p-splitter-panel) {
  overflow: hidden;
}

:deep(.p-component) {
  height: 100%;
  overflow: hidden;
}

:deep(.p-toggleable-content) {
  height: calc(100% - 40px);
}

:deep(.p-panel-content) {
  height: 100%;
  overflow: auto;
}

:deep(.p-panel-content) {
  padding: 0.75rem;
}
</style>
