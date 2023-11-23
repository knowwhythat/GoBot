<template>
  <Panel toggleable class="w-96">
    <template #header>
      <div class="flex align-items-center gap-2">
        <Avatar
          image="https://primefaces.org/cdn/primevue/images/avatar/amyelsner.png"
          size="large"
          shape="circle"
        />
        <span class="font-bold">Amy Elsner</span>
      </div>
    </template>
    <template #icons>
      <button class="p-panel-header-icon p-link mr-2" @click="toggle">
        <span class="pi pi-cog"></span>
      </button>
      <Menu ref="menu" id="config_menu" :model="items" popup />
    </template>
    <slot></slot>
  </Panel>
</template>
<script setup>
import Button from "primevue/button";
import Panel from "primevue/panel";
</script>
