<template>
  <div class="main">
    <Panel
      class="body"
      :header="title"
      :pt="{
        header: (options) => ({
          style: '--wails-draggable:drag',
        }),
      }"
    >
      <template #icons>
        <SystemOperate
          @quit="confirmQuit"
          @toggleMax="appStore.changeMainWindowState($event)"
        >
          <div
            class="hover:bg-slate-200 p-2"
            @mouseenter="showUserInfo"
            @mouseleave="hideUserInfo"
          >
            <v-remixicon size="20" name="riUser3Line" />
          </div>
          <OverlayPanel ref="userInfo">
            <div
              class="flex flex-column gap-3 w-25rem"
              @mouseenter="showUserInfo"
              @mouseleave="hideUserInfo"
            >
              <p class="items-center content-center">Admin</p>
              <Button
                @click="logout"
                icon="pi pi-sign-out"
                severity="danger"
                rounded
                aria-label="Cancel"
              />
            </div>
          </OverlayPanel>
        </SystemOperate>
      </template>
      <router-view v-slot="{ Component }">
        <transition name="slide-left">
          <component :is="Component" />
        </transition>
      </router-view>
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
import { onMounted, ref } from "vue";
import { useRouter } from "vue-router";
import { Quit, WindowMaximise, WindowUnmaximise } from "@back/runtime/runtime";
import OverlayPanel from "primevue/overlaypanel";
import Panel from "primevue/panel";
import Dock from "primevue/dock";
import TerminalSvg from "@/assets/images/terminal.svg";
import DesignerSvg from "@/assets/images/designer.svg";
import TriggerSvg from "@/assets/images/trigger.svg";
import RecordSvg from "@/assets/images/record.svg";
import SettingSvg from "@/assets/images/setting.svg";
import SystemOperate from "@/components/SystemOperate.vue";
import { useConfirm } from "primevue/useconfirm";
import { useAppStore } from "@/stores/app";
import { Logout } from "@back/go/backend/App";

const appStore = useAppStore();
const confirm = useConfirm();

onMounted(async () => {
  if (appStore.mainWindow.isMax) {
    await WindowMaximise();
  } else {
    await WindowUnmaximise();
  }
});

const userInfo = ref();
const userInfoShow = ref(false);

const showUserInfo = async (e) => {
  userInfoShow.value = true;
  userInfo.value.show(e);
};
const hideUserInfo = async () => {
  userInfoShow.value = false;
  setTimeout(() => {
    if (!userInfoShow.value) {
      userInfo.value.hide();
    }
  }, 300);
};
const router = useRouter();
const title = ref("工作流");
const items = ref([
  {
    label: "工作流",
    icon: DesignerSvg,
    command: () => {
      title.value = "工作流";
      router.push("/main");
    },
  },
  {
    label: "触发器",
    icon: TriggerSvg,
    command: () => {
      title.value = "触发器";
      router.push("/main/trigger");
    },
  },
  {
    label: "执行记录",
    icon: RecordSvg,
    command: () => {
      title.value = "执行记录";
      router.push("/main/record");
    },
  },
  {
    label: "运行监控",
    icon: TerminalSvg,
    command: () => {
      title.value = "运行监控";
      router.push("/main/monitor");
    },
  },
  {
    label: "设置",
    icon: SettingSvg,
    command: () => {
      title.value = "设置";
      router.push("/main/setting");
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
const logout = async () => {
  confirm.require({
    message: "退出后所有的任务将停止,是否退出登录?",
    header: "退出登录",
    icon: "pi pi-exclamation-triangle",
    rejectClass: "p-button-secondary p-button-outlined",
    rejectLabel: "取消",
    acceptLabel: "确认",
    accept: async () => {
      await Logout();
      await router.replace("/");
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

:deep(.p-panel .p-panel-content) {
  border: none;
}
</style>
