<template>
  <div :class="[data.color || 'white', colors[data.color || 'white']]" class="p-4 rounded-lg block-note"
    style="min-width: 192px">
    <div class="pb-2 border-b flex items-center">
      <v-remixicon name="riFileEditLine" size="20" />
      <p class="flex-1 ml-2 mr-2 font-semibold">注释</p>
      <ConfirmPopup class="note-color" group="setting">
        <template #message="slotProps">
          <div class="p-2">
            <p class="mb-1 ml-1 text-sm text-gray-600 dark:text-gray-200">颜色</p>
            <div class="flex items-center space-x-2">
              <span v-for="(color, colorId) in colors" :key="colorId" :class="color" style="border-width: 3px"
                class="h-8 w-8 rounded-full inline-block cursor-pointer" @click="updateData({ color: colorId })" />
            </div>
            <p class="py-2 mb-1 ml-1 text-sm text-gray-600 dark:text-gray-200">字体大小</p>
            <div class="flex flex-row gap-2 flex-wrap">
              <div v-for="(font, key) in fontSize" :key="key" class="flex align-items-center">
                <RadioButton input-class="border-2 border-lime-200" v-model="data.fontSize" :inputId="font.name"
                  name="dynamic" :value="key" />
                <label :for="font.name" class="ml-2">{{ font.name }}</label>
              </div>
            </div>
          </div>
        </template>
      </ConfirmPopup>
      <Button @click="setting($event)">
        <template #icon>
          <v-remixicon name="riSettings3Line" size="20" class="cursor-pointer" />
        </template>
      </Button>
      <hr class="mx-2 h-7 border-r" />
      <v-remixicon name="riDeleteBin7Line" size="20" class="cursor-pointer" @click="$emit('delete', id)" />
    </div>
    <textarea :value="data.note" :style="initialSize" :class="[fontSize[data.fontSize || 'regular'].class]"
      placeholder="在这里填写注释..." cols="30" rows="7" style="resize: both; min-width: 160px; min-height: 88px"
      @input="updateData({ note: $event.target.value })" class="bg-transparent" @mousedown.stop @mouseup="onMouseup" />
  </div>
</template>
<script setup>
import { debounce } from '@/utils/helper';
import { useConfirm } from "primevue/useconfirm";
import ConfirmPopup from 'primevue/confirmpopup';
import Button from 'primevue/button';
import RadioButton from 'primevue/radiobutton';
import Textarea from 'primevue/textarea';

const props = defineProps({
  id: {
    type: String,
    default: '',
  },
  label: {
    type: String,
    default: '',
  },
  data: {
    type: Object,
    default: () => ({}),
  },
});
const emit = defineEmits(['update', 'delete']);

const confirm = useConfirm();

const setting = (event) => {
  confirm.require({
    target: event.currentTarget,
    group: 'setting',
    acceptClass: 'invisible',
    rejectClass: 'invisible'
  });
};

const initialSize = {
  width: `${props.data.width}px`,
  height: `${props.data.height}px`,
};

const colors = {
  white: 'bg-white dark:bg-gray-800',
  red: 'bg-red-200 dark:bg-red-300',
  indigo: 'bg-indigo-200 dark:bg-indigo-300',
  green: 'bg-green-200 dark:bg-green-300',
  amber: 'bg-amber-200 dark:bg-amber-300',
  sky: 'bg-sky-200 dark:bg-sky-300',
  violet: 'bg-violet-200 dark:bg-violet-300',
};
const fontSize = {
  regular: {
    name: '常规',
    class: 'text-base',
  },
  medium: {
    name: '中',
    class: 'text-xl',
  },
  large: {
    name: '大',
    class: 'text-2xl',
  },
  'extra-large': {
    name: '超大',
    class: 'text-3xl',
  },
};

const updateData = debounce((data) => {
  emit('update', data);
}, 250);

function onMouseup({ target }) {
  let { height, width } = target.style;
  width = parseInt(width, 10);
  height = parseInt(height, 10);

  if (width === props.data.width && height === props.data.height) return;

  updateData({ height, width });
}
</script>
<style>
.note-color .ui-popover__trigger {
  @apply flex items-center;
}

.block-note * {
  border-color: rgb(0 0 0 / 12%);
}

.dark .block-note {
  &:not(.white) {
    @apply text-gray-900;
  }

  &.white * {
    border-color: rgb(255 255 255 / 12%);
  }

  * {
    border-color: rgb(0 0 0 / 12%);
  }
}
</style>
