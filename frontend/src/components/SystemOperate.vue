<template>
  <div class="flex gap-4 -mt-5">
    <slot></slot>
    <div class="hover:bg-slate-200 p-2" @click="WindowMinimise">
      <v-remixicon size="20" viewBox="0 0 1024 1024" name="minIcon" />
    </div>
    <div class="hover:bg-slate-200 p-2" @click="toggleMaximise">
      <v-remixicon
        v-if="isMax"
        size="20"
        viewBox="0 0 1024 1024"
        name="restoreIcon"
      />
      <v-remixicon v-else size="20" viewBox="0 0 1024 1024" name="maxIcon" />
    </div>
    <div class="hover:bg-slate-200 p-2" @click="$emit('quit')">
      <v-remixicon size="20" viewBox="0 0 1024 1024" name="closeIcon" />
    </div>
  </div>
</template>
<script setup>
import {
  WindowIsMaximised,
  WindowMinimise,
  WindowToggleMaximise,
} from "@back/runtime/runtime";
import { nextTick, onMounted, onUnmounted, ref } from "vue";

const emit = defineEmits(["toggleMax", "min", "quit"]);

onMounted(async () => {
  isMax.value = await WindowIsMaximised();
  window.addEventListener("resize", handleResize);
});

onUnmounted(() => {
  window.removeEventListener("resize", handleResize);
});

const isMax = ref(false);

async function handleResize() {
  isMax.value = await WindowIsMaximised();
}

async function toggleMaximise() {
  await WindowToggleMaximise();
  nextTick(async () => {
    emit("toggleMax", await WindowIsMaximised());
  });
}
</script>
