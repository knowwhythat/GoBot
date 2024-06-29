<template>
  <div style="height: 100%">
    <DataView
      :value="showProjects"
      :layout="layout"
      paginator
      :rows="12"
      :alwaysShowPaginator="false"
    >
      <template #header>
        <Toolbar class="bg-gray-100 p-2">
          <template #start>
            <InputGroup>
              <InputGroupAddon>
                <i class="pi pi-search" />
              </InputGroupAddon>
              <InputText
                v-model="searchText"
                placeholder="输入关键字进行搜索"
              />
            </InputGroup>
          </template>
          <template #end>
            <div class="flex gap-4">
              <Button label="导入">
                <template #icon>
                  <v-remixicon name="riDownloadLine" />
                </template>
              </Button>
              <Button label="新建" @click="createProject">
                <template #icon>
                  <v-remixicon name="riAddLine" />
                </template>
              </Button>
            </div>
          </template>
        </Toolbar>
      </template>

      <template #grid="slotProps">
        <div
          class="grid sm:grid-cols-2 md:grid-cols-2 lg:grid-cols-4 gap-4 mt-2"
        >
          <div
            v-for="(item, index) in slotProps.items"
            :key="index"
            class="border-gray-200 border-current border rounded-lg hover:bg-gray-100"
            @dblclick="edit(item)"
          >
            <div class="p-4 border-1 surface-border border-round">
              <div
                class="flex flex-wrap align-items-center justify-content-between gap-2"
              >
                <div class="flex align-items-center gap-2">
                  <v-remixicon size="18" name="riHistoryLine" />
                  <span class="-mt-1">
                    {{ item?.updateTs.substring(0, 19).replace("T", " ") }}
                  </span>
                </div>
              </div>
              <div class="flex flex-row align-items-center gap-3 py-3">
                <div class="text-xl font-bold truncate">
                  {{ item?.name }}
                </div>
              </div>
              <div class="pb-1">{{}}</div>
              <div class="flex flex-row-reverse gap-1">
                <Button
                  rounded
                  outlined
                  icon="pi pi-ellipsis-v"
                  @click="toggle(index, $event)"
                  v-tooltip="'其他'"
                />
                <TieredMenu
                  ref="operateMenu"
                  id="overlay_tmenu"
                  :model="getItems(item?.id)"
                  popup
                />
                <Button
                  rounded
                  outlined
                  icon="pi pi-file-edit"
                  @click="edit(item)"
                  v-tooltip="'编辑'"
                />
                <Button rounded outlined @click="run(item)" v-tooltip="'运行'">
                  <template #icon>
                    <v-remixicon size="24" name="riPlayLine" />
                  </template>
                </Button>
                <Button
                  rounded
                  outlined
                  @click="openParamDialog(item)"
                  v-tooltip="'设置参数'"
                >
                  <template #icon>
                    <v-remixicon size="24" name="riGridLine" />
                  </template>
                </Button>
              </div>
            </div>
          </div>
        </div>
      </template>
    </DataView>
    <Dialog
      v-model:visible="dialogVisible"
      modal
      maximizable
      :header="(newProject.id ? '修改' : '新建') + '流程'"
    >
      <div class="flex justify-content-center">
        <label for="value" class="mt-3 w-32">流程名称</label>
        <InputText
          id="value"
          v-model="newProject.name"
          type="text"
          autocomplete="off"
          class="w-full ml-6"
        />
      </div>
      <div class="flex justify-content-center mt-2">
        <label for="value" class="w-32">是否使用流程图</label>
        <InputSwitch v-model="newProject.isFlow" />
      </div>
      <div class="flex justify-content-center mt-2">
        <label for="value" class="w-32">描述</label>
        <Editor v-model="newProject.desc" editorStyle="height: 240px" />
      </div>
      <template #footer>
        <Button
          label="取消"
          icon="pi pi-times"
          outlined
          severity="secondary"
          @click="dialogVisible = false"
        />
        <Button label="确定" icon="pi pi-check" @click="doCreateProject" />
      </template>
    </Dialog>
    <Dialog
      v-model:visible="paramDiaglogVisible"
      modal
      maximizable
      header="流程参数"
      :style="{ width: '60rem' }"
    >
      <ParamSettingView :params="projectParamDefine" v-model="projectParam" />
      <template #footer>
        <Button
          label="取消"
          icon="pi pi-times"
          outlined
          severity="secondary"
          @click="paramDiaglogVisible = false"
        />
        <Button label="确定" icon="pi pi-check" @click="saveParams" />
      </template>
    </Dialog>
  </div>
</template>

<script setup>
import { ref, onMounted, reactive, watch } from "vue";
import { useRouter } from "vue-router";
import TieredMenu from "primevue/tieredmenu";
import InputGroup from "primevue/inputgroup";
import InputGroupAddon from "primevue/inputgroupaddon";
import InputSwitch from "primevue/inputswitch";
import DataView from "primevue/dataview";
import Toolbar from "primevue/toolbar";
import InputText from "primevue/inputtext";
import Button from "primevue/button";
import Dialog from "primevue/dialog";
import Editor from "primevue/editor";
import ParamSettingView from "@/views/main/ParamSettingView.vue";
import { useToast } from "primevue/usetoast";
import { throttle } from "lodash-es";
import {
  AddOrUpdateProject,
  ListProject,
  DeleteProject,
  RunMainFlow,
  GetSubFlow,
} from "@back/go/backend/App";
import { useConfirm } from "primevue/useconfirm";
import { useAppStore } from "@/stores/app";
const appStore = useAppStore();
const toast = useToast();
const router = useRouter();
const confirm = useConfirm();

const projects = ref();
const showProjects = ref();
const layout = ref("grid");

onMounted(() => {
  listProject();
});

const searchText = ref("");
watch(
  searchText,
  throttle((newVal, oldVal) => {
    showProjects.value = projects.value.filter((p) => p.name.includes(newVal));
  }, 1000)
);

function listProject() {
  ListProject()
    .then((result) => {
      projects.value = result;
      showProjects.value = result;
    })
    .catch((error) => {
      toast.add({
        severity: "error",
        summary: "加载失败",
        detail: error,
        life: 3000,
      });
      console.error(error);
    });
}
const operateMenu = ref();
const toggle = (index, event) => {
  operateMenu.value[index].toggle(event);
};
function getItems(id) {
  return [
    {
      label: "编辑",
      icon: "pi pi-file-edit",
      command: () => {
        const pro = projects.value.find((p) => p.id === id);
        newProject.id = pro.id;
        newProject.name = pro.name;
        newProject.desc = pro.description;
        dialogVisible.value = true;
      },
    },
    {
      label: "删除",
      icon: "pi pi-times",
      command: () => {
        confirm.require({
          header: "提示",
          message: "确定要删除这条记录吗?",
          icon: "pi pi-info-circle",
          rejectClass: "p-button-secondary p-button-outlined p-button-sm",
          acceptClass: "p-button-danger p-button-sm",
          rejectLabel: "取消",
          acceptLabel: "删除",
          accept: () => {
            appStore.changeLoadingState(true);
            DeleteProject(id)
              .then((result) => {
                appStore.changeLoadingState(false);
                listProject();
              })
              .catch((error) => {
                appStore.changeLoadingState(false);
                console.error(error);
              });
          },
          reject: () => {},
        });
      },
    },
    {
      label: "发布",
      icon: "pi pi-send",
      command: () => {
        toast.add({
          severity: "warn",
          detail: "功能尚未发布",
          life: 3000,
        });
      },
    },
  ];
}

const paramDiaglogVisible = ref(false);

const projectParamDefine = ref([]);
const projectParam = ref({});
const selectedProject = ref({});

const openParamDialog = async (item) => {
  selectedProject.value = item;
  const flow = await GetSubFlow(item.id, "main");
  projectParamDefine.value = JSON.parse(flow).params;
  let inputParam = item.inputParam ?? {};
  projectParamDefine.value.forEach((element) => {
    if (!(element.name in inputParam)) {
      inputParam[element.name] = element.value;
    }
  });
  projectParam.value = inputParam;
  paramDiaglogVisible.value = true;
};

const saveParams = () => {
  selectedProject.value.inputParam = projectParam.value;
  AddOrUpdateProject(selectedProject.value)
    .then((result) => {
      paramDiaglogVisible.value = false;
      listProject();
    })
    .catch((error) => {
      toast.add({
        severity: "error",
        summary: "异常",
        detail: error,
        life: 3000,
      });
    });
};

const run = (item) => {
  toast.add({
    severity: "success",
    detail: "流程启动成功",
    life: 3000,
  });
  RunMainFlow(item.id).catch((err) => {
    toast.add({
      severity: "error",
      summary: "流程运行异常",
      detail: err,
      life: 3000,
    });
  });
};

const edit = (item) => {
  router.push("/design?id=" + item.id);
};

const dialogVisible = ref(false);
const newProject = reactive({
  id: "",
  name: "",
  isFlow: false,
  desc: "",
});

function createProject() {
  newProject.id = "";
  newProject.name = "";
  newProject.isFlow = false;
  newProject.desc = "";
  dialogVisible.value = true;
}

function doCreateProject() {
  if (!newProject.name) {
    toast.add({
      severity: "warn",
      summary: "提示",
      detail: "流程名称不能为空",
      life: 3000,
    });
    return;
  }
  appStore.changeLoadingState(true);
  AddOrUpdateProject(newProject)
    .then((result) => {
      appStore.changeLoadingState(false);
      dialogVisible.value = false;
      listProject();
    })
    .catch((error) => {
      appStore.changeLoadingState(false);
      toast.add({
        severity: "error",
        summary: "异常",
        detail: error,
        life: 3000,
      });
    });
}
</script>
<style scoped>
:deep(.p-dataview-content) {
  height: calc(100vh - 298px);
  overflow: auto;
}

:deep(.p-grid) {
  @apply sm:grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4;
}
</style>
