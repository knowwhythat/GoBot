<template>
  <div class="main">
    <div class="card">
      <Toast position="top-center" group="tc" />
      <div class="dock-window dock-advanced">
        <Dock :model="items">
          <template #item="{ item }">
            <a v-tooltip.top="item.label" href="#" class="p-dock-link" @click="onDockItemClick($event, item)">
              <img :alt="item.label" :src="item.icon" style="width: 100%" />
            </a>
          </template>
        </Dock>

        <Dialog v-model:visible="displayTerminal" header="Terminal" :breakpoints="{ '960px': '80vw' }"
          :style="{ width: '80vw', height: '80vh' }" :maximizable="true">
          <Terminal prompt="$" />
        </Dialog>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onBeforeUnmount } from 'vue'
import Dock from 'primevue/dock'
import { useToast } from 'primevue/usetoast'
import Dialog from 'primevue/dialog'
import Toast from 'primevue/toast'
import Terminal from 'primevue/terminal'
import TerminalService from 'primevue/terminalservice'

onMounted(() => {
  TerminalService.on('command', commandHandler);
})

onBeforeUnmount(() => {
  TerminalService.off('command', commandHandler);
})

const displayTerminal = ref(false);
const toast = useToast();
const items = ref([
  {
    label: 'Terminal',
    icon: "https://primefaces.org/cdn/primevue//images/dock/terminal.svg",
    command: () => {
      displayTerminal.value = true;
    }
  },
  {
    label: 'App Store',
    icon: "https://primefaces.org/cdn/primevue//images/dock/appstore.svg",
    command: () => {
      toast.add({ severity: 'error', summary: 'An unexpected error occurred while signing in.', detail: 'UNTRUSTED_CERT_TITLE', group: 'tc', life: 3000 });
    }
  },
  {
    label: 'Safari',
    icon: "https://primefaces.org/cdn/primevue//images/dock/safari.svg",
    command: () => {
      toast.add({ severity: 'warn', summary: 'Safari has stopped working', group: 'tc', life: 3000 });
    }
  },
  {
    label: 'GitHub',
    icon: "https://primefaces.org/cdn/primevue//images/dock/github.svg",
  },
  {
    label: 'Trash',
    icon: "https://primefaces.org/cdn/primevue//images/dock/trash.png",
    command: () => {
      toast.add({ severity: 'info', summary: 'Empty Trash', life: 3000 });
    }
  }
]);

const onDockItemClick = (event, item) => {
  if (item.command) {
    item.command();
  }

  event.preventDefault();
};

const commandHandler = (text) => {
  let response;
  let argsIndex = text.indexOf(' ');
  let command = argsIndex !== -1 ? text.substring(0, argsIndex) : text;

  switch (command) {
    case "date":
      response = 'Today is ' + new Date().toDateString();
      break;

    case "greet":
      response = 'Hola ' + text.substring(argsIndex + 1);
      break;

    case "random":
      response = Math.floor(Math.random() * 100);
      break;

    default:
      response = "Unknown command: " + command;
  }

  TerminalService.emit('response', response);
};
</script>

<style scoped>
.main {
  width: 100vw;
  height: 100vh;
  position: relative;
  background-image: url('@/assets/images/window.jpg');
  background-repeat: no-repeat;
  background-size: cover;
  padding: 0;
}
</style>