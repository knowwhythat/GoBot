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
          <Button
            label="停止"
            severity="success"
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
        </template>
      </Column>
    </DataTable>
  </div>
</template>
<script setup>
import Button from "primevue/button";
import InputGroup from "primevue/inputgroup";
import InputGroupAddon from "primevue/inputgroupaddon";
import DataTable from "primevue/datatable";
import { FilterMatchMode } from "primevue/api";
import Column from "primevue/column";
import Toolbar from "primevue/toolbar";
import { useToast } from "primevue/usetoast";
import { ref, onMounted, onUnmounted, watch } from "vue";
import { GetRunningFlows, TerminateMainFlow } from "@back/go/backend/App";
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
</script>
