<template>
  <div>
    <Toast />
    <ConfirmDialog />
    <BlockUI :blocked="appStore.loading" fullScreen />
    <router-view></router-view>
  </div>
</template>
<script setup>
import Toast from "primevue/toast";
import ConfirmDialog from "primevue/confirmdialog";
import BlockUI from "primevue/blockui";
import { useAppStore } from "@/stores/app";
const appStore = useAppStore();
</script>
