<template>
  <div style="height: 100%">
    <DataView :value="projects" :layout="layout">
      <template #header>
        <Toolbar class="bg-gray-100 p-2">
          <template #start>
            <div class="flex gap-1">
              <span class="p-input-icon-left">
                <i class="pi pi-search" />
                <InputText
                  v-model="searchText"
                  style="margin-left: 40px"
                  placeholder="Search"
                />
              </span>
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
          <template #end>
            <SelectButton
              v-model="layout"
              :options="options"
              optionValue="value"
              optionLabel="value"
              dataKey="value"
              aria-labelledby="custom"
            >
              <template #option="slotProps">
                <i :class="slotProps.option.icon" size="32"></i>
              </template>
            </SelectButton>
          </template>
        </Toolbar>
      </template>

      <template #list="slotProps">
        <div
          class="border-gray-200 border-current border rounded-lg p-2 m-2"
          v-for="(item, index) in slotProps.items"
          :key="index"
        >
          <div class="activity-node p-4 bg-input flex items-center group">
            <div class="leading-tight flex-1 overflow-hidden">
              <div class="text-xl font-bold truncate">
                {{ item?.name }}
              </div>
              <div
                class="flex flex-wrap align-items-center justify-content-between gap-2 pt-2"
              >
                <div class="flex align-items-center gap-2">
                  <v-remixicon size="18" name="riPriceTag3Line" />
                  <span class="font-semibold">测试</span>
                </div>
              </div>
              <div
                class="flex flex-wrap align-items-center justify-content-between gap-2 pt-2"
              >
                <div class="flex align-items-center gap-2">
                  <span class="font-semibold">2023-11-24 10:00</span>
                </div>
              </div>
            </div>
            <div class="invisible group-hover:visible ml-8">
              <v-remixicon
                name="riPencilLine"
                size="24"
                class="cursor-pointer inline-block mr-2"
                @click="edit(item?.id)"
              />
              <v-remixicon
                name="riSettings3Line"
                size="24"
                class="cursor-pointer inline-block mr-2"
                @click="editItemSettings(element)"
              />
              <v-remixicon
                name="riDeleteBin7Line"
                size="24"
                class="cursor-pointer inline-block"
                @click="deleteItem(index, element.itemId)"
              />
            </div>
          </div>
        </div>
      </template>

      <template #grid="slotProps">
        <div
          class="grid sm:grid-cols-2 md:grid-cols-2 lg:grid-cols-4 gap-4 mt-2"
        >
          <div
            v-for="(item, index) in slotProps.items"
            :key="index"
            class="border-gray-200 border-current border rounded-lg"
          >
            <div class="p-4 border-1 surface-border border-round">
              <div
                class="flex flex-wrap align-items-center justify-content-between gap-2"
              >
                <div class="flex align-items-center gap-2">
                  <v-remixicon size="18" name="riPriceTag3Line" />
                  <span class="font-semibold">{{ item?.category }}</span>
                </div>
              </div>
              <div class="flex flex-column align-items-center gap-3 py-5">
                <div class="text-xl font-bold truncate">
                  {{ item?.name }}
                </div>
              </div>
              <div class="flex flex-row-reverse gap-1">
                <SplitButton
                  label="设置"
                  icon="pi pi-play"
                  @click="run(item?.id)"
                  :model="getItems(item?.id)"
                >
                  <template #icon>
                    <v-remixicon size="24" name="riSettings3Line" />
                  </template>
                </SplitButton>
                <Button
                  icon="pi pi-file-edit"
                  @click="edit(item?.id)"
                  v-tooltip="'编辑'"
                />
                <Button icon="pi pi-caret-right" v-tooltip="'运行'">
                  <template #icon>
                    <v-remixicon size="24" name="riPlayLine" />
                  </template>
                </Button>
              </div>
            </div>
          </div>
        </div>
      </template>
    </DataView>
    <Dialog
      v-model:visible="newProject.dialogVisible"
      modal
      maximizable
      header="Header"
    >
      <template #header>
        <div
          class="inline-flex align-items-center justify-content-center gap-2"
        >
          <v-remixicon size="20" name="riAddLine" />
          <span class="font-bold white-space-nowrap">新建流程</span>
        </div>
      </template>
      <div class="flex justify-content-center">
        <label for="value" class="mt-3 w-20">流程名称</label>
        <InputText
          id="value"
          v-model="newProject.name"
          type="text"
          autocomplete="off"
        />
      </div>
      <div class="flex justify-content-center mt-2">
        <label for="value" class="w-20">描述</label>
        <Editor v-model="newProject.desc" editorStyle="height: 240px" />
      </div>
      <template #footer>
        <Button
          label="取消"
          icon="pi pi-times"
          outlined
          severity="secondary"
          @click="newProject.dialogVisible = false"
        />
        <Button
          label="确定"
          icon="pi pi-check"
          @click="doCreateProject"
          :loading="newProject.loading"
        />
      </template>
    </Dialog>
  </div>
</template>

<script setup>
import { ref, onMounted, reactive, watch } from "vue";
import { useRouter } from "vue-router";
import DataView from "primevue/dataview";
import Toolbar from "primevue/toolbar";
import InputText from "primevue/inputtext";
import SelectButton from "primevue/selectbutton";
import SplitButton from "primevue/splitbutton";
import Button from "primevue/button";
import Dialog from "primevue/dialog";
import Editor from "primevue/editor";
import { useToast } from "primevue/usetoast";
import { throttle } from "lodash-es";
import { CreateProject, ListProject, DeleteProject } from "@back/go/main/App";
import { useConfirm } from "primevue/useconfirm";
const toast = useToast();
const router = useRouter();
const confirm = useConfirm();

const projects = ref();
const layout = ref("grid");
const options = ref([
  { icon: "pi pi-table", value: "grid" },
  { icon: "pi pi-list", value: "list" },
]);

onMounted(() => {
  listProject();
});

const searchText = ref("");
watch(
  searchText,
  throttle((newVal, oldVal) => {
    listProject(newVal);
  }, 1000)
);

function listProject(name) {
  ListProject(name)
    .then((result) => {
      projects.value = result.list;
      console.log(projects);
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

function getItems(id) {
  return [
    {
      label: "编辑",
      icon: "pi pi-file-edit",
      command: () => {},
    },
    {
      label: "删除",
      icon: "pi pi-times",
      command: () => {
        confirm.require({
          message: "确定要删除这条记录吗?",
          icon: "pi pi-info-circle",
          rejectClass: "p-button-secondary p-button-outlined p-button-sm",
          acceptClass: "p-button-danger p-button-sm",
          rejectLabel: "取消",
          acceptLabel: "删除",
          accept: () => {
            DeleteProject(id)
              .then((result) => {
                listProject();
              })
              .catch((error) => {
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
          summary: "Delete",
          detail: "Data Deleted",
          life: 3000,
        });
      },
    },
  ];
}

const run = () => {
  toast.add({
    severity: "success",
    summary: "Success",
    detail: "Data Saved",
    life: 3000,
  });
};

const edit = (id) => {
  router.push("/design?id=" + id);
};

const newProject = reactive({
  dialogVisible: false,
  loading: false,
  name: "",
  desc: "",
});

function createProject() {
  newProject.dialogVisible = true;
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
  newProject.loading = true;
  CreateProject(newProject.name)
    .then((result) => {
      newProject.loading = false;
      newProject.dialogVisible = false;
      listProject();
    })
    .catch((error) => {
      console.error(error);
      newProject.loading = false;
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
