<template>
  <block-base :id="componentId" :data="data" :block-id="id" :block-data="block" class="w-64" @edit="$emit('edit')"
    @delete="$emit('delete', id)" @update="$emit('update', $event)" @settings="$emit('settings', $event)">
    <Handle :id="`${id}-input-1`" type="target" :position="appStore.inputPosition" />
    <div class="flex items-center">
      <div :class="block.color" class="inline-block text-sm mr-4 p-2 rounded-lg dark:text-black">
        <v-remixicon name="riAB" size="20" class="inline-block mr-1" />
        <span>{{ block.name }}</span>
      </div>
      <div class="flex-grow pointer-events-none" />
      <p class="p-1 rounded-lg dark:text-black" :class="block.color">
        <v-remixicon name="riAddLine" size="20" class="m-1" />
      </p>
    </div>
    <ul v-if="data.conditions && data.conditions.length !== 0" class="mt-4 space-y-2">
      <li v-for="item in data.conditions" :key="item.id"
        class="flex items-center flex-1 p-2 bg-box-transparent rounded-lg w-full relative"
        @dblclick.stop="$emit('edit', { editCondition: item.id })">
        <p v-if="item.name" class="text-overflow w-full text-right" :title="item.name">
          {{ item.name }}
        </p>
        <template v-else>
          <p class="w-5/12 text-overflow text-right">
            {{ item.compareValue || '_____' }}
          </p>
          <p class="w-2/12 text-center mx-1 font-mono">
            {{ item.type }}
          </p>
          <p class="w-5/12 text-overflow">
            {{ item.value || '_____' }}
          </p>
        </template>
        <Handle :id="`${id}-output-${item.id}`" :position="appStore.outputPosition" style="margin-right: -33px"
          type="source" />
      </li>
      <p v-if="data.conditions && data.conditions.length !== 0" class="text-right text-gray-600 dark:text-gray-200">
        <span title="Fallback"> &#9432; </span>
        Fallback
      </p>
    </ul>
    <Handle v-if="data.conditions.length > 0" :id="`${id}-output-fallback`" :position="appStore.outputPosition"
      type="source" style="top: auto; bottom: 10px" />
  </block-base>
</template>
<script setup>
import { useAppStore } from '@/stores/app';
import { Handle } from '@vue-flow/core';
import { useComponentId } from '@/composable/componentId';
import { getBlocks } from '@/utils/getSharedData';
import BlockBase from './BlockBase.vue';
import Button from 'primevue/button';

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
defineEmits(['delete', 'settings', 'edit', 'update']);

const appStore = useAppStore();

const componentId = useComponentId('block-conditions');
const blocks = getBlocks();
const block = blocks[props.label];
</script>
<style>
.condition-handle {
  position: relative !important;
  top: 82px !important;
  margin-bottom: 32px !important;
}
</style>
