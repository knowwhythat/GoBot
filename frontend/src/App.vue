<template>
  <div>
    <Toast />
    <ConfirmDialog />
    <router-view></router-view>
  </div>
</template>
<script setup>
import Toast from "primevue/toast";
import ConfirmDialog from "primevue/confirmdialog";
</script>
