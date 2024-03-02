<template>
  <div class="main">
    <Panel
      class="body"
      :header="title"
      :pt="{
        header: (options) => ({
          style: {
            'user-select': 'none',
          },
        }),
      }"
    >
      <template #icons>
        <SystemOperate @quit="confirmQuit" />
      </template>
      <router-view></router-view>
    </Panel>

    <div>
      <Dock :model="items">
        <template #item="{ item }">
          <a
            v-tooltip.top="item.label"
            href="#"
            class="p-dock-link"
            @click="onDockItemClick($event, item)"
          >
            <img :alt="item.label" :src="item.icon" style="width: 100%" />
          </a>
        </template>
      </Dock>
    </div>
  </div>
</template>

<script setup>
import { ref } from "vue";
import { useRouter } from "vue-router";
import { Quit } from "@back/runtime/runtime";
import Panel from "primevue/panel";
import Dock from "primevue/dock";
import TerminalSvg from "@/assets/images/terminal.svg";
import DesignerSvg from "@/assets/images/designer.svg";
import TriggerSvg from "@/assets/images/trigger.svg";
import RecordSvg from "@/assets/images/record.svg";
import SystemOperate from "@/components/SystemOperate.vue";
import { useConfirm } from "primevue/useconfirm";
const confirm = useConfirm();

const router = useRouter();
const title = ref("工作流");
const items = ref([
  {
    label: "工作流",
    icon: DesignerSvg,
    command: () => {
      title.value = "工作流";
      router.push("/");
    },
  },
  {
    label: "触发器",
    icon: TriggerSvg,
    command: () => {
      title.value = "触发器";
      router.push("/trigger");
    },
  },
  {
    label: "执行记录",
    icon: RecordSvg,
    command: () => {
      title.value = "执行记录";
      router.push("/record");
    },
  },
  {
    label: "日志",
    icon: TerminalSvg,
    command: () => {
      title.value = "日志";
      router.push("/log");
    },
  },
]);

const onDockItemClick = (event, item) => {
  if (item.command) {
    item.command();
  }

  event.preventDefault();
};

const confirmQuit = () => {
  confirm.require({
    message: "关闭程序后所有的任务将停止,是否确定关闭?",
    header: "关闭程序",
    icon: "pi pi-exclamation-triangle",
    rejectClass: "p-button-secondary p-button-outlined",
    rejectLabel: "取消",
    acceptLabel: "确认",
    accept: () => {
      Quit();
    },
    reject: () => {},
  });
};
</script>

<style scoped>
.main {
  width: 100vw;
  height: 100vh;
  position: relative;
  background-repeat: no-repeat;
  background-size: cover;
  padding: 0;
}

.body {
  height: calc(100vh - 98px);
}
</style>
