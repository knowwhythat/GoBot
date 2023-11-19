<template>
  <div
    class="activity-node p-4 border-2 m-2 rounded-lg bg-input flex items-center group"
    :group-item-id="element.itemId"
    style="cursor: grab"
    :activity-data="element"
  >
    <div class="leading-tight flex-1 overflow-hidden">
      <p class="w-80 text-overflow">
        {{ element.label }}
      </p>
    </div>
    <div class="invisible group-hover:visible ml-8">
      <v-remixicon
        name="riPencilLine"
        size="18"
        class="cursor-pointer inline-block mr-2"
        @click="editBlock(element)"
      />
      <v-remixicon
        name="riSettings3Line"
        size="18"
        class="cursor-pointer inline-block mr-2"
        @click="editItemSettings(element)"
      />
      <v-remixicon
        name="riDeleteBin7Line"
        size="18"
        class="cursor-pointer inline-block"
        @click="deleteItem(index, element.itemId)"
      />
    </div>
  </div>
</template>
<script setup>
import { onMounted } from "vue";
const props = defineProps({
  element: {
    type: Object,
    default: {},
  },
});
onMounted(() => {
  console.log(props);
});
</script>
