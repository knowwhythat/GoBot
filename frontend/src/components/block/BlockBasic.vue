<template>
  <block-base
    :id="id"
    :data="data"
    class="block-basic group"
    @edit="$emit('edit', id)"
    @delete="$emit('delete', id)"
    @run="$emit('run', id)"
    @settings="$emit('settings', $event)"
  >
    <Handle
      v-if="data?.nodeType !== 'Start'"
      :id="`${id}-input`"
      type="target"
      :position="appStore.inputPosition"
    />
    <div class="flex items-center">
      <span
        :class="block.color"
        class="inline-block p-2 mr-2 rounded-lg dark:text-black"
      >
        <v-remixicon :name="block?.icon || 'riGlobalLine'" />
      </span>
      <div class="overflow-hidden flex-1">
        <p class="font-semibold leading-tight text-overflow whitespace-nowrap">
          {{ data.label }}
        </p>
      </div>
    </div>
    <div
      v-if="data.errorEnable && data.toDo === 'fallback'"
      class="fallback flex items-center justify-end"
    >
      <v-remixicon title="异常处理策略" name="riInformationLine" size="18" />
      <span class="ml-1"> 异常执行 </span>
    </div>
    <Handle
      :id="`${id}-output`"
      type="source"
      :position="appStore.outputPosition"
    />
    <Handle
      v-if="data.errorEnable && data.toDo === 'fallback'"
      :id="`${id}-fallback`"
      type="source"
      :position="appStore.outputPosition"
      :style="fallbackStyle"
    />
  </block-base>
</template>
<script setup>
import { useAppStore } from "@/stores/app";
import { Handle } from "@vue-flow/core";
import BlockBase from "./BlockBase.vue";
import { BlockType } from "@/utils/shared.js";
import { computed } from "vue";

const props = defineProps({
  id: {
    type: String,
    default: "",
  },
  data: {
    type: Object,
    default: () => ({}),
  },
});
defineEmits(["delete", "edit", "run", "settings"]);

const appStore = useAppStore();
const block = BlockType[props.data.nodeType];

const fallbackStyle = computed(() => {
  if (appStore.settings.editor.layout === "vertical") {
    return { left: "auto", right: "10px" };
  } else {
    return { top: "auto", bottom: "10px" };
  }
});
</script>
