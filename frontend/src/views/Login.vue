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
        <Image :src="LoginSvg" alt="Image" width="450" />
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
              v-model="loginForm.username"
            />
          </InputGroup>
          <InputGroup class="py-2">
            <InputGroupAddon>
              <i class="pi pi-lock pr-2"></i>
            </InputGroupAddon>
            <Password
              v-model="loginForm.pwd"
              :feedback="false"
              toggleMask
              placeholder="密码"
            />
          </InputGroup>
          <div class="flex">
            <div class="flex align-items-center py-2">
              <Checkbox
                v-model="loginForm.rememberMe"
                inputId="rememberMe"
                :binary="true"
              />
              <label for="rememberMe" class="ml-2"> 记住密码 </label>
            </div>
            <div class="flex align-items-center py-2 ml-6">
              <Checkbox
                v-model="loginForm.autoLogin"
                inputId="autoLogin"
                :binary="true"
                @change="autoLoginChange"
              />
              <label for="autoLogin" class="ml-2"> 自动登录 </label>
            </div>
          </div>
          <Button class="mt-6" label="登录" @click="login"></Button>
        </div>
      </div>
    </Panel>
  </div>
</template>

<script setup>
import { useRouter } from "vue-router";
import { onMounted, reactive, ref } from "vue";
import Panel from "primevue/panel";
import InputGroupAddon from "primevue/inputgroupaddon";
import InputGroup from "primevue/inputgroup";
import SystemOperate from "@/components/SystemOperate.vue";
import { useToast } from "primevue/usetoast";
import { useAppStore } from "@/stores/app";
import { Quit } from "@back/runtime/runtime.js";
import LoginSvg from "@/assets/images/login.svg";
import { Login, GetLoginData } from "@back/go/backend/App";
const appStore = useAppStore();
const toast = useToast();
const loginForm = reactive({
  username: "",
  pwd: "",
  rememberMe: false,
  autoLogin: false,
});
const router = useRouter();
onMounted(async () => {
  const data = await GetLoginData();
  loginForm.username = data.username;
  loginForm.pwd = data.pwd;
  loginForm.rememberMe = data.rememberMe;
  loginForm.autoLogin = data.autoLogin;
  if (data.autoLogin) {
    await login();
  }
});
function autoLoginChange(event) {
  if (loginForm.autoLogin) {
    loginForm.rememberMe = true;
  }
}
async function login() {
  if (!loginForm.username) {
    toast.add({
      severity: "warn",
      detail: "用户名不能为空",
      life: 3000,
    });
    return;
  }
  if (!loginForm.pwd) {
    toast.add({
      severity: "warn",
      detail: "密码不能为空",
      life: 3000,
    });
    return;
  }
  try {
    appStore.changeLoadingState(true);
    await Login(loginForm);
    await router.push("/main");
    appStore.changeLoadingState(false);
  } catch (err) {
    appStore.changeLoadingState(false);
    toast.add({
      severity: "error",
      summary: "登录异常",
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
