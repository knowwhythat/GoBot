<template>
  <div>
    <Panel
      class="main"
      header="登录"
      :pt="{
        header: (options) => ({
          style: '--wails-draggable:drag',
        }),
      }"
    >
      <template #icons>
        <SystemOperate
          @quit="Quit"
          @toggleMax="appStore.changeMainWindowState($event)"
        />
      </template>
      <div class="body flex justify-center items-center">
        <Image :src="qrCodeUrl" alt="Image" width="450" />
        <div
          class="flex flex-col ml-40 shadow-lg p-20 gap-4 bg-slate-50 rounded-3xl hover:bg-slate-100"
        >
          <div class="flex justify-center font-serif text-3xl">
            GoBot-自动化系统
          </div>
          <InputGroup class="py-2 w-full">
            <InputGroupAddon>
              <i class="pi pi-user pr-2"></i>
            </InputGroupAddon>
            <InputText
              placeholder="用户名"
              class="w-96"
              v-model="registerForm.username"
            />
          </InputGroup>
          <InputGroup class="py-2">
            <InputGroupAddon>
              <i class="pi pi-lock pr-2"></i>
            </InputGroupAddon>
            <Password
              v-model="registerForm.pwd"
              toggleMask
              placeholder="密码"
            />
          </InputGroup>
          <Button
            class="mt-6"
            label="注册"
            @click="register"
            v-tooltip="'用微信扫描左侧二维码后再进行注册'"
          ></Button>
          <div class="flex flex-row-reverse align-items-center py-2">
            <Button label="登录" link @click="router.replace('/')" />
          </div>
        </div>
      </div>
    </Panel>
  </div>
</template>

<script setup>
import { useRoute, useRouter } from "vue-router";
import { onMounted, reactive, ref } from "vue";
import Panel from "primevue/panel";
import InputGroupAddon from "primevue/inputgroupaddon";
import InputGroup from "primevue/inputgroup";
import SystemOperate from "@/components/SystemOperate.vue";
import { useToast } from "primevue/usetoast";
import { useAppStore } from "@/stores/app";
import { Quit } from "@back/runtime/runtime.js";
import { GetQrCode, Register } from "@back/go/backend/App";

const appStore = useAppStore();
const toast = useToast();
const registerForm = reactive({
  username: "",
  pwd: "",
});
const router = useRouter();
const route = useRoute();
const qrCodeUrl = ref();
onMounted(async () => {
  qrCodeUrl.value = await GetQrCode();
});

async function register() {
  if (!registerForm.username) {
    toast.add({
      severity: "warn",
      detail: "用户名不能为空",
      life: 3000,
    });
    return;
  }
  if (!registerForm.pwd) {
    toast.add({
      severity: "warn",
      detail: "密码不能为空",
      life: 3000,
    });
    return;
  }
  try {
    appStore.changeLoadingState(true);
    await Register(registerForm.username, registerForm.pwd);
    toast.add({
      severity: "success",
      detail: "注册成功",
      life: 3000,
    });
    await router.push("/");
    appStore.changeLoadingState(false);
  } catch (err) {
    appStore.changeLoadingState(false);
    toast.add({
      severity: "error",
      summary: "注册异常",
      detail: err,
      life: 3000,
    });
  }
}
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

:deep(.p-password-input) {
  @apply w-96;
}
</style>
