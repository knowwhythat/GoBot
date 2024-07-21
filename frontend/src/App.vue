<template>
  <div>
    <Toast />
    <ConfirmDialog />
    <Dialog
      v-model:visible="appStore.loading"
      :modal="true"
      :closable="false"
      :showHeader="false"
      :blockScroll="true"
      :draggable="false"
      :pt="{
        root: (options) => ({
          style: {
            'box-shadow': 'none',
          },
        }),
        content: (options) => ({
          style: {
            background: '#0000',
          },
        }),
      }"
    >
      <div class="mb-2">请稍等……</div>
      <ProgressSpinner
        style="width: 50px; height: 50px"
        strokeWidth="8"
        fill="var(--surface-ground)"
        animationDuration=".5s"
        aria-label="Custom ProgressSpinner"
      />
    </Dialog>
    <router-view v-slot="{ Component }">
      <transition name="slide-left">
        <component :is="Component" />
      </transition>
    </router-view>
  </div>
</template>
<script setup>
import Dialog from "primevue/dialog";
import Toast from "primevue/toast";
import ProgressSpinner from "primevue/progressspinner";
import ConfirmDialog from "primevue/confirmdialog";
import { useAppStore } from "@/stores/app";

const appStore = useAppStore();
</script>
<style scoped>
.slide-right-enter-active,
.slide-right-leave-active,
.slide-left-enter-active,
.slide-left-leave-active {
  height: 100%;
  will-change: transform;
  transition: all 500ms;
  position: absolute;
  backface-visibility: hidden;
}
.slide-right-enter {
  opacity: 0;
  transform: translate3d(-100%, 0, 0);
}
.slide-right-leave-active {
  opacity: 0;
  transform: translate3d(100%, 0, 0);
}
.slide-left-enter {
  opacity: 0;
  transform: translate3d(100%, 0, 0);
}
.slide-left-leave-active {
  opacity: 0;
  transform: translate3d(-100%, 0, 0);
}

.van-badge--fixed {
  z-index: 1000;
}
</style>
