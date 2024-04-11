<template>
  <div>
    <DataTable
      :value="executions"
      :globalFilterFields="['name']"
      v-model:selection="selectedExecution"
      v-model:filters="filters"
      selectionMode="multiple"
      :alwaysShowPaginator="false"
      paginator
      :rows="5"
      :rowsPerPageOptions="[5, 10, 20, 50]"
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
          <template #end>
            <div class="flex gap-4">
              <Button
                label="删除"
                @click="deleteExecution"
                :disabled="!selectedExecution"
              >
                <template #icon>
                  <v-remixicon name="riDeleteBin7Line" />
                </template>
              </Button>
            </div>
          </template>
        </Toolbar>
      </template>
      <Column field="projectName" header="项目名称" style="width: 25%"></Column>
      <Column field="triggerType" header="触发方式" style="width: 10%"></Column>
      <Column header="执行结果" style="width: 10%">
        <template #body="slotProps">
          <Tag
            v-if="slotProps.data.executeResult == 1"
            icon="pi pi-check"
            severity="success"
            value="成功"
          ></Tag>
          <Tag
            v-else-if="slotProps.data.executeResult == 0"
            icon="pi pi-times"
            severity="danger"
            value="失败"
          ></Tag>
          <Tag
            v-else
            icon="pi pi-exclamation-triangle"
            severity="warning"
            value="手动停止"
          ></Tag>
        </template>
      </Column>
      <Column header="开始时间" style="width: 15%">
        <template #body="slotProps">
          {{ slotProps.data.startTs.substring(0, 19).replace("T", " ") }}
        </template>
      </Column>
      <Column header="结束时间" style="width: 15%">
        <template #body="slotProps">
          {{ slotProps.data.endTs.substring(0, 19).replace("T", " ") }}
        </template>
      </Column>
      <Column header="操作" style="width: 15%">
        <template #body="slotProps">
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
        </template>
      </Column>
    </DataTable>
    <Dialog v-model:visible="dialogVisible" modal maximizable>
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
import Tag from "primevue/tag";
import InputGroup from "primevue/inputgroup";
import InputGroupAddon from "primevue/inputgroupaddon";
import Button from "primevue/button";
import DataTable from "primevue/datatable";
import { FilterMatchMode } from "primevue/api";
import Column from "primevue/column";
import Dialog from "primevue/dialog";
import Toolbar from "primevue/toolbar";
import { useToast } from "primevue/usetoast";
import { ref, onMounted, onUnmounted, watch } from "vue";
import {
  QueryAllExecution,
  RemoveExecution,
  SelectExecutionLog,
} from "@back/go/backend/App";
import { EventsOn, EventsOff } from "@back/runtime/runtime";
import { useAppStore } from "@/stores/app";
import { useConfirm } from "primevue/useconfirm";
const confirm = useConfirm();
const appStore = useAppStore();
const toast = useToast();
const executions = ref();
const filters = ref({
  global: { value: null, matchMode: FilterMatchMode.CONTAINS },
});
const selectedExecution = ref(null);
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
  executions.value = await QueryAllExecution();
}

function deleteExecution() {
  if (!selectedExecution.value) {
    return;
  }
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
      RemoveExecution(selectedExecution.value)
        .then((result) => {
          appStore.changeLoadingState(false);
          loadExecutions();
        })
        .catch((error) => {
          appStore.changeLoadingState(false);
          toast.add({
            severity: "error",
            summary: "删除失败",
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
async function showLogs(id) {
  const result = await SelectExecutionLog(id);
  console.log(result);
  logs.value = result.split(/\r\n|\n/);
  dialogVisible.value = true;
}
</script>
