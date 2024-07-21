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
    <div class="flex flex-row flex-wrap gap-2 p-2">
      <div
        class="flex h-32 w-32 items-center place-content-center bg-slate-100 select-none"
        v-for="image in images"
        :key="image.id"
        @dblclick.stop="showImage(image)"
        @contextmenu.stop="showContextMenu(image, $event)"
      >
        <img
          :src="'data:image/png;base64,' + image.image"
          alt="Image"
          class="hover:scale-105"
        />
      </div>
    </div>
    <ContextMenu ref="menu" :model="items" />
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
  </Panel>
</template>
<script setup>
import Image from "primevue/image";
import Panel from "primevue/panel";
import Dialog from "primevue/dialog";
import {
  HighlightImage,
  RemoveImage,
  SaveImage,
  StartCapture,
} from "@back/go/backend/App";
import { WindowMaximise, WindowMinimise } from "@back/runtime/runtime.js";
import { useToast } from "primevue/usetoast";
import { useConfirm } from "primevue/useconfirm";
import Button from "primevue/button";
import { ref, toRaw } from "vue";
import { nanoid } from "nanoid";
import InputText from "primevue/inputtext";
import InputNumber from "primevue/inputnumber";
import Dropdown from "primevue/dropdown";
import { matchTypes } from "@/utils/shared.js";
import ContextMenu from "primevue/contextmenu";

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
    const resp = await StartCapture();
    if (resp) {
      const data = JSON.parse(resp);
      data.id = nanoid(16);
      data.name = "图片" + props.images.length;
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

function showImage(image) {
  imageData.value = image;
  visible.value = true;
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
    await WindowMinimise();
    await HighlightImage(props.id, imageData.value);
    toast.add({
      severity: "success",
      summary: "找到元素",
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

const menu = ref();

function showContextMenu(image, event) {
  menu.value.show(event);
  imageData.value = image;
}

const items = ref([
  {
    label: "删除",
    icon: "pi pi-times-circle",
    command: () => {
      confirm.require({
        header: "提示",
        message: "确定要删除这条记录吗?",
        icon: "pi pi-info-circle",
        rejectClass: "p-button-secondary p-button-outlined p-button-sm",
        acceptClass: "p-button-danger p-button-sm",
        rejectLabel: "取消",
        acceptLabel: "删除",
        accept: async () => {
          await RemoveImage(props.id, imageData.value.id);
        },
        reject: () => {},
      });
    },
  },
  {
    label: "重新截图",
    icon: "pi pi-file-edit",
    command: async () => {
      try {
        await WindowMinimise();
        const resp = await StartCapture();
        if (resp) {
          const data = JSON.parse(resp);
          data.id = imageData.value.id;
          data.name = imageData.value.name;
          data.threshold = imageData.value.threshold;
          data.match_type = imageData.value.match_type;
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
    },
  },
  {
    label: "查看",
    icon: "pi pi-file-edit",
    command: () => {
      visible.value = true;
    },
  },
]);
</script>
<style scoped>
:deep(.p-splitter-panel) {
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

:deep(.p-dialog .p-dialog-header) {
  padding: 0.5rem;
}
</style>
