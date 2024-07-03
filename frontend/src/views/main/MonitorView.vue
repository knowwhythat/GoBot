<template>
  <div>
    <DataTable
      :value="executions"
      :globalFilterFields="['name']"
      v-model:filters="filters"
      :rowHover="true"
      scrollable
      scrollHeight="flex"
      resizableColumns
    >
      <template #header>
        <Toolbar class="bg-gray-100 p-2">
          <template #start>
            <InputGroup>
              <InputGroupAddon>
                <i class="pi pi-search" />
              </InputGroupAddon>
              <InputText
                v-model="filters['global'].value"
                placeholder="输入关键字进行搜索"
              />
            </InputGroup>
          </template>
        </Toolbar>
      </template>
      <Column field="projectName" header="项目名称" style="width: 25%"></Column>
      <Column field="triggerType" header="触发方式" style="width: 10%"></Column>
      <Column header="开始时间" style="width: 15%">
        <template #body="slotProps">
          {{ slotProps.data.startTs.substring(0, 19).replace("T", " ") }}
        </template>
      </Column>
      <Column header="操作" style="width: 15%">
        <template #body="slotProps">
          <div class="flex flex-row gap-2">
            <Button
              label="停止"
              severity="danger"
              text
              :pt="{
                root: (options) => ({
                  style: {
                    padding: '0rem',
                  },
                }),
              }"
              @click="terminateFlow(slotProps.data.id)"
            />
            <Button
              label="查看日志"
              severity="success"
              text
              :pt="{
                root: (options) => ({
                  style: {
                    padding: '0rem',
                  },
                }),
              }"
              @click="showLogs(slotProps.data.id)"
            />
          </div>
        </template>
      </Column>
    </DataTable>
    <Dialog
      v-model:visible="dialogVisible"
      modal
      maximizable
      @after-hide="hideLogView"
    >
      <template #header>
        <div
          class="inline-flex align-items-center justify-content-center gap-2"
        >
          <v-remixicon size="20" name="riFileTextLine" />
          <span class="font-bold white-space-nowrap">运行日志</span>
        </div>
      </template>
      <div class="flex gap-4 flex-col w-[56rem]">
        <p v-for="log in logs">{{ log }}</p>
      </div>
      <template #footer>
        <Button
          label="取消"
          icon="pi pi-times"
          outlined
          severity="secondary"
          @click="dialogVisible = false"
        />
      </template>
    </Dialog>
  </div>
</template>
<script setup>
import Dialog from "primevue/dialog";
import Button from "primevue/button";
import InputGroup from "primevue/inputgroup";
import InputGroupAddon from "primevue/inputgroupaddon";
import DataTable from "primevue/datatable";
import { FilterMatchMode } from "primevue/api";
import Column from "primevue/column";
import Toolbar from "primevue/toolbar";
import { useToast } from "primevue/usetoast";
import { onMounted, onUnmounted, ref } from "vue";
import {
  GetRunningFlows,
  SelectExecutionLog,
  TerminateMainFlow,
  StartMonitorLog,
  StopMonitorLog,
} from "@back/go/backend/App";
import { EventsOff, EventsOn } from "@back/runtime/runtime";
import { useAppStore } from "@/stores/app";
import { useConfirm } from "primevue/useconfirm";

const confirm = useConfirm();
const appStore = useAppStore();
const toast = useToast();
const executions = ref();
const filters = ref({
  global: { value: null, matchMode: FilterMatchMode.CONTAINS },
});
onMounted(async () => {
  await loadExecutions();
  EventsOn("execute_event", (data) => {
    loadExecutions();
  });
});

onUnmounted(() => {
  EventsOff("execute_event");
});

async function loadExecutions() {
  executions.value = await GetRunningFlows();
}

function terminateFlow(id) {
  confirm.require({
    header: "提示",
    message: "确定要停止该流程吗?",
    icon: "pi pi-info-circle",
    rejectClass: "p-button-secondary p-button-outlined p-button-sm",
    acceptClass: "p-button-danger p-button-sm",
    rejectLabel: "取消",
    acceptLabel: "确定",
    accept: () => {
      appStore.changeLoadingState(true);
      TerminateMainFlow(id)
        .then((result) => {
          appStore.changeLoadingState(false);
          loadExecutions();
        })
        .catch((error) => {
          appStore.changeLoadingState(false);
          toast.add({
            severity: "error",
            summary: "失败",
            detail: error,
            life: 3000,
          });
        });
    },
    reject: () => {},
  });
}

const dialogVisible = ref(false);
const logs = ref([]);
let monitorId = "";
async function showLogs(id) {
  logs.value = [];
  monitorId = id;
  await StartMonitorLog(id);
  dialogVisible.value = true;
  EventsOn("log_event", (data) => {
    if (data.id === id) {
      logs.value.push(data.text);
    }
  });
}

async function hideLogView() {
  EventsOff("log_event");
  await StopMonitorLog(monitorId);
  monitorId = "";
}
</script>
