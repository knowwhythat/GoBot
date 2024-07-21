<template>
  <div>
    <DataTable
      :value="schedules"
      :globalFilterFields="['name']"
      v-model:selection="selectedTrigger"
      v-model:filters="filters"
      selectionMode="single"
      :alwaysShowPaginator="false"
      paginator
      :rows="10"
      :rowsPerPageOptions="[5, 10, 20, 50]"
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
                @click="deleteSchedule"
                :disabled="!selectedTrigger"
              >
                <template #icon>
                  <v-remixicon name="riDeleteBin7Line" />
                </template>
              </Button>
              <Button
                label="修改"
                @click="updateSchedule"
                :disabled="!selectedTrigger"
              >
                <template #icon>
                  <v-remixicon name="riEditBoxLine" />
                </template>
              </Button>
              <Button label="新建" @click="CreateSchedule">
                <template #icon>
                  <v-remixicon name="riAddLine" />
                </template>
              </Button>
            </div>
          </template>
        </Toolbar>
      </template>
      <Column field="name" header="名称" style="width: 20%"></Column>
      <Column field="desc" header="描述" style="width: 20%"></Column>
      <Column field="cron" header="定时设置" style="width: 20%"></Column>
      <Column header="是否启用" style="width: 10%">
        <template #body="slotProps">
          <InputSwitch
            :trueValue="true"
            :falseValue="false"
            v-model="slotProps.data.enabled"
            @update:modelValue="changeState($event, slotProps.data.id)"
          />
        </template>
      </Column>
      <Column header="创建日期" style="width: 15%">
        <template #body="slotProps">
          {{ slotProps.data.createTs.substring(0, 10) }}
        </template>
      </Column>
      <Column header="更新日期" style="width: 15%">
        <template #body="slotProps">
          {{ slotProps.data.updateTs.substring(0, 10) }}
        </template>
      </Column>
    </DataTable>
    <Dialog v-model:visible="dialogVisible" modal maximizable>
      <template #header>
        <div
          class="inline-flex align-items-center justify-content-center gap-2"
        >
          <v-remixicon size="20" name="riAddLine" />
          <span class="font-bold white-space-nowrap">新建触发器</span>
        </div>
      </template>
      <div class="flex gap-4 flex-col w-[56rem]">
        <div class="flex justify-content-center">
          <label for="name" class="mt-3 w-32">触发器名称</label>
          <InputText
            placeholder="输入触发器名称"
            inputId="name"
            v-model="newSchedule.name"
            type="text"
            autocomplete="off"
            class="w-full"
          />
        </div>
        <div class="flex justify-content-center">
          <label for="projectId" class="mt-3 w-32">流程名称</label>
          <Dropdown
            v-model="newSchedule.projectId"
            filter
            :options="projects"
            optionLabel="name"
            optionValue="id"
            placeholder="选择流程"
            class="w-full"
          />
        </div>
        <div class="flex flex-wrap gap-8">
          <label for="cronType" class="w-32">频率</label>
          <div class="flex align-items-center">
            <RadioButton
              v-model="cronConfig.cronType"
              inputId="minute"
              name="minute"
              value="minute"
            />
            <label for="minute" class="ml-2">每分</label>
          </div>
          <div class="flex align-items-center">
            <RadioButton
              v-model="cronConfig.cronType"
              inputId="hour"
              name="hour"
              value="hour"
            />
            <label for="hour" class="ml-2">每时</label>
          </div>
          <div class="flex align-items-center">
            <RadioButton
              v-model="cronConfig.cronType"
              inputId="day"
              name="day"
              value="day"
            />
            <label for="day" class="ml-2">每天</label>
          </div>
          <div class="flex align-items-center">
            <RadioButton
              v-model="cronConfig.cronType"
              inputId="week"
              name="week"
              value="week"
            />
            <label for="week" class="ml-2">每周</label>
          </div>
          <div class="flex align-items-center">
            <RadioButton
              v-model="cronConfig.cronType"
              inputId="month"
              name="month"
              value="month"
            />
            <label for="month" class="ml-2">月份</label>
          </div>
          <div class="flex align-items-center">
            <RadioButton
              v-model="cronConfig.cronType"
              inputId="senior"
              name="senior"
              value="senior"
            />
            <label for="senior" class="ml-2">高级</label>
          </div>
        </div>
        <div
          v-if="cronConfig.cronType == 'senior'"
          class="flex align-items-center gap-2 ml-32"
        >
          <label for="cron" class="w-32 mt-3">时间表达式:</label>
          <InputText
            v-model="cronConfig.cron"
            inputId="cron"
            type="text"
            autocomplete="off"
            class="w-full"
          />
        </div>
        <div v-if="cronConfig.cronType == 'month'" class="flex gap-2 ml-32">
          <div class="flex align-items-center">
            <label for="month" class="w-12 mt-3">月份</label>
            <InputNumber
              v-model="cronConfig.month"
              inputId="month"
              mode="decimal"
              showButtons
              :min="1"
              :max="11"
            />
          </div>
        </div>
        <div
          v-if="cronConfig.cronType == 'month' || cronConfig.cronType == 'week'"
          class="flex flex-wrap gap-8 ml-32"
        >
          <div class="flex align-items-center">
            <Checkbox
              inputId="monday"
              v-model="cronConfig.week[0]"
              :binary="true"
            />
            <label for="monday" class="ml-2">周一</label>
          </div>
          <div class="flex align-items-center">
            <Checkbox
              inputId="tuesday"
              v-model="cronConfig.week[1]"
              :binary="true"
            />
            <label for="tuesday" class="ml-2">周二</label>
          </div>
          <div class="flex align-items-center">
            <Checkbox
              inputId="wednesday"
              v-model="cronConfig.week[2]"
              :binary="true"
            />
            <label for="wednesday" class="ml-2">周三</label>
          </div>
          <div class="flex align-items-center">
            <Checkbox
              inputId="thursday"
              v-model="cronConfig.week[3]"
              :binary="true"
            />
            <label for="thursday" class="ml-2">周四</label>
          </div>
          <div class="flex align-items-center">
            <Checkbox
              inputId="friday"
              v-model="cronConfig.week[4]"
              :binary="true"
            />
            <label for="friday" class="ml-2">周五</label>
          </div>
          <div class="flex align-items-center">
            <Checkbox
              inputId="saturday"
              v-model="cronConfig.week[5]"
              :binary="true"
            />
            <label for="saturday" class="ml-2">周六</label>
          </div>
          <div class="flex align-items-center">
            <Checkbox
              inputId="sunday"
              v-model="cronConfig.week[6]"
              :binary="true"
            />
            <label for="sunday" class="ml-2">周日</label>
          </div>
        </div>
        <div class="flex gap-2 ml-32">
          <div
            v-if="
              cronConfig.cronType == 'month' ||
              cronConfig.cronType == 'week' ||
              cronConfig.cronType == 'day' ||
              cronConfig.cronType == 'hour'
            "
            class="flex align-items-center w-54"
          >
            <label for="hour" class="w-12 mt-3">小时</label>
            <InputNumber
              v-model="cronConfig.hour"
              inputId="hour"
              mode="decimal"
              showButtons
              :min="1"
              :max="11"
            />
          </div>
          <div
            v-if="
              cronConfig.cronType == 'month' ||
              cronConfig.cronType == 'week' ||
              cronConfig.cronType == 'day' ||
              cronConfig.cronType == 'hour' ||
              cronConfig.cronType == 'minute'
            "
            class="flex align-items-center w-54"
          >
            <label for="minitus" class="w-12 mt-3">分钟</label>
            <InputNumber
              v-model="cronConfig.minute"
              inputId="minitus"
              mode="decimal"
              showButtons
              :min="1"
              :max="59"
            />
          </div>
        </div>
        <div class="flex gap-2 ml-32">
          <div class="flex align-items-center">
            {{ newSchedule.desc }}
            <i
              class="pi pi-stopwatch"
              style="font-size: 1.5rem; margin-left: 4px"
              v-tooltip.bottom="{
                value: nextTriggerTime,
                pt: {
                  text: {
                    style: {
                      width: '210px',
                    },
                  },
                },
              }"
            ></i>
          </div>
        </div>
        <div class="flex justify-content-center">
          <label for="enableTrigger" class="w-32">是否启用</label>
          <InputSwitch inputId="enableTrigger" v-model="newSchedule.enabled" />
        </div>
      </div>
      <template #footer>
        <Button
          label="取消"
          icon="pi pi-times"
          outlined
          severity="secondary"
          @click="dialogVisible = false"
        />
        <Button label="确定" icon="pi pi-check" @click="doCreateSchedule" />
      </template>
    </Dialog>
  </div>
</template>
<script setup>
import InputGroup from "primevue/inputgroup";
import InputGroupAddon from "primevue/inputgroupaddon";
import Button from "primevue/button";
import DataTable from "primevue/datatable";
import { FilterMatchMode } from "primevue/api";
import Column from "primevue/column";
import Dialog from "primevue/dialog";
import Dropdown from "primevue/dropdown";
import RadioButton from "primevue/radiobutton";
import Checkbox from "primevue/checkbox";
import Toolbar from "primevue/toolbar";
import { useToast } from "primevue/usetoast";
import { onMounted, ref, watch } from "vue";
import InputSwitch from "primevue/inputswitch";
import {
  AddOrUpdateSchedule,
  GetNextTriggerTime,
  ListProject,
  QuerySchedulePage,
  RemoveSchedule,
  ToggleScheduleById,
} from "@back/go/backend/App";
import { useAppStore } from "@/stores/app";
import { useConfirm } from "primevue/useconfirm";

const confirm = useConfirm();
const appStore = useAppStore();
const toast = useToast();
const schedules = ref();
const filters = ref({
  global: { value: null, matchMode: FilterMatchMode.CONTAINS },
});
const selectedTrigger = ref(null);
const projects = ref([]);
onMounted(async () => {
  await loadSchedules();
  projects.value = await ListProject();
});

async function loadSchedules() {
  schedules.value = await QuerySchedulePage();
}

const dialogVisible = ref(false);
const newSchedule = ref({
  projectId: "",
  name: "",
  enabled: false,
  desc: "",
  cron: "",
  config: "",
});

const cronConfig = ref({
  cronType: "minute",
  minute: 1,
  hour: 1,
  day: 1,
  week: [false, false, false, false, false, false, false],
  month: 1,
  cron: "",
});
const nextTriggerTime = ref("");
watch(
  () => cronConfig,
  async (value, oldValue) => {
    if (value.value.cronType === "minute") {
      newSchedule.value.desc = `每隔${value.value.minute}分钟执行`;
      newSchedule.value.cron = `0 0/${value.value.minute} * * * ?`;
    } else if (value.value.cronType === "hour") {
      newSchedule.value.desc = `每间隔${value.value.hour}小时,在第${value.value.minute}分钟执行`;
      newSchedule.value.cron = `0 ${value.value.minute} 0/${value.value.hour} * * ?`;
    } else if (value.value.cronType === "day") {
      newSchedule.value.desc = `每天的${value.value.hour}时${value.value.minute}分执行`;
      newSchedule.value.cron = `0 ${value.value.minute} ${value.value.hour} * * ?`;
    } else if (value.value.cronType === "week") {
      let weekEnable = "";
      let weekEnableNumber = "";
      for (let index = 0; index < 7; index++) {
        if (cronConfig.value.week[index]) {
          weekEnable += weekname[index] + "、";
          weekEnableNumber += index + 1 + ",";
        }
      }

      if (weekEnable.length > 0) {
        if (weekEnable.charAt(weekEnable.length - 1) === "、") {
          weekEnable = weekEnable.substring(0, weekEnable.length - 1);
          weekEnableNumber = weekEnableNumber.substring(
            0,
            weekEnableNumber.length - 1,
          );
        }
        newSchedule.value.desc = `每${weekEnable}的${value.value.hour}时${value.value.minute}分执行`;
        newSchedule.value.cron = `0 ${value.value.minute} ${value.value.hour} * * ${weekEnableNumber}`;
      } else {
        newSchedule.value.desc = `每天的${value.value.hour}时${value.value.minute}分执行`;
        newSchedule.value.cron = `0 ${value.value.minute} ${value.value.hour} * * ?`;
      }
    } else if (value.value.cronType === "month") {
      let weekEnable = "";
      let weekEnableNumber = "";
      for (let index = 0; index < 7; index++) {
        if (cronConfig.value.week[index]) {
          weekEnable += weekname[index] + "、";
          weekEnableNumber += index + 1 + ",";
        }
      }

      if (weekEnable.length > 0) {
        if (weekEnable.charAt(weekEnable.length - 1) === "、") {
          weekEnable = weekEnable.substring(0, weekEnable.length - 1);
          weekEnableNumber = weekEnableNumber.substring(
            0,
            weekEnableNumber.length - 1,
          );
        }
        newSchedule.value.desc = `${value.value.month}月每${weekEnable}的${value.value.hour}时${value.value.minute}分执行`;
        newSchedule.value.cron = `0 ${value.value.minute} ${value.value.hour} * ${value.value.month} ${weekEnableNumber}`;
      } else {
        newSchedule.value.desc = `${value.value.month}月每天的${value.value.hour}时${value.value.minute}分执行`;
        newSchedule.value.cron = `0 ${value.value.minute} ${value.value.hour} * ${value.value.month} ${weekEnableNumber}`;
      }
    } else if (value.value.cronType === "senior") {
      newSchedule.value.cron = value.value.cron;
      newSchedule.value.desc = "用户自定义";
    }
    try {
      const result = await GetNextTriggerTime(newSchedule.value.cron);
      nextTriggerTime.value = `预期将在以下时间点执行: ${result}`;
    } catch (err) {
      nextTriggerTime.value = "定时任务表达式不符合规范:" + err;
    }
  },
  { immediate: true, deep: true },
);
const weekname = ["周一", "周二", "周三", "周四", "周五", "周六", "周日"];
function deleteSchedule() {
  if (!selectedTrigger.value) {
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
      RemoveSchedule(selectedTrigger.value.id)
        .then((result) => {
          appStore.changeLoadingState(false);
          loadSchedules();
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
function updateSchedule() {
  if (!selectedTrigger.value) {
    return;
  }
  cronConfig.value = JSON.parse(selectedTrigger.value.config);
  newSchedule.value = selectedTrigger.value;
  nextTriggerTime.value = "";
  dialogVisible.value = true;
}
function CreateSchedule() {
  cronConfig.value = {
    cronType: "minute",
    minute: 1,
    hour: 1,
    day: 1,
    week: [false, false, false, false, false, false, false],
    month: 1,
    cron: "",
  };
  newSchedule.value = {
    projectId: "",
    name: "",
    enable: false,
    desc: "",
    cron: "",
    config: "",
  };
  nextTriggerTime.value = "";
  dialogVisible.value = true;
}
async function doCreateSchedule() {
  if (!newSchedule.value.name) {
    toast.add({
      severity: "warn",
      summary: "提示",
      detail: "触发器名称不能为空",
      life: 3000,
    });
    return;
  }
  if (!newSchedule.value.projectId) {
    toast.add({
      severity: "warn",
      summary: "提示",
      detail: "流程名称不能为空",
      life: 3000,
    });
    return;
  }
  newSchedule.value.config = JSON.stringify(cronConfig.value);
  appStore.changeLoadingState(true);
  try {
    await AddOrUpdateSchedule(newSchedule.value);
    await loadSchedules();
    dialogVisible.value = false;
  } catch (err) {
    toast.add({
      severity: "error",
      summary: "添加失败",
      detail: err,
      life: 3000,
    });
  }
  appStore.changeLoadingState(false);
}

async function changeState(event, id) {
  try {
    await ToggleScheduleById(id, event);
  } catch (err) {
    toast.add({
      severity: "error",
      summary: "失败",
      detail: err,
      life: 3000,
    });
  }
}
</script>
<style scoped>
:deep(.p-inputnumber-input) {
  @apply w-24;
}
</style>
