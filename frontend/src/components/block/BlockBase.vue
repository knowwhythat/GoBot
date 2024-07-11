<template>
  <div class="block-base relative w-48" @dblclick.stop="$emit('edit')">
    <div
      class="top-0 w-full absolute block-menu-container hidden"
      style="transform: translateY(-100%)"
    >
      <div class="inline-flex items-center dark:text-gray-300 block-menu">
        <button
          v-if="!block.disableDelete"
          title="删除"
          @click.stop="$emit('delete')"
        >
          <v-remixicon size="20" name="riDeleteBin7Line" />
        </button>
        <button
          v-if="!block.disableSetting"
          title="设置"
          @click.stop="$emit('settings', { data, id })"
        >
          <v-remixicon size="20" name="riSettings3Line" />
        </button>
        <button title="从这里运行" @click.stop="$emit('run')">
          <v-remixicon size="20" name="riPlayLine" />
        </button>
        <button
          v-if="!block.disableEdit"
          title="编辑"
          @click.stop="$emit('edit')"
        >
          <v-remixicon size="20" name="riPencilLine" />
        </button>
        <slot name="action" />
      </div>
    </div>
    <slot name="prepend" />
    <div
      class="z-10 relative bg-gray-200 hover:bg-gray-300 focus:ring p-2 rounded-xl"
    >
      <slot></slot>
    </div>
    <slot name="append" />
  </div>
</template>
<script setup>
import { BlockType } from "@/utils/shared.js";

const props = defineProps({
  data: {
    type: Object,
    default: () => ({}),
  },
  id: {
    type: String,
    default: "",
  },
});
defineEmits(["delete", "edit", "run", "settings"]);
const block = BlockType[props.data.nodeType];
</script>
<style>
.block-menu {
  @apply mb-1 rounded-md bg-slate-400;

  button {
    padding-left: 6px;
    padding-right: 6px;
    @apply focus:ring-0 py-1 hover:text-primary;
  }
}
</style>
