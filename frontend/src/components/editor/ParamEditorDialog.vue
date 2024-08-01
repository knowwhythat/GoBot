<template>
  <Dialog
    :visible="visible"
    v-on:update:visible="$emit('hide')"
    maximizable
    modal
    @show="setData"
    header="参数设置"
    :style="{ width: '60rem' }"
    :breakpoints="{ '1199px': '75vw', '575px': '90vw' }"
    :pt="{
      header: (options) => ({
        style: {
          padding: '0.75rem',
        },
      }),
    }"
  >
    <template #header>
      <div class="flex align-items-center gap-2">
        <span
          :class="props.color"
          class="inline-block rounded-lg dark:text-black"
        >
          <v-remixicon v-bind="getIconPath(props.icon)" />
        </span>
        <p class="w-64 truncate">{{ parameterData.label }}</p>
      </div>
    </template>

    <TabView>
      <TabPanel header="常规">
        <div class="px-2">
          <div class="flex mt-4">
            <span class="mr-2 w-32 truncate my-auto"> 组件名称 </span>
            <InputText
              class="w-full mr-2"
              type="text"
              v-model="parameterData.label"
              @keydown="(e) => e.stopPropagation()"
            />
            <a title="提示" target="_blank" rel="noopener" class="my-auto p-4">
              <v-remixicon name="riInformationLine" size="18" />
            </a>
          </div>
          <Divider align="left" type="solid">
            <b>输入参数</b>
          </Divider>
          <template
            v-for="input in props.parameter_define?.inputs"
            :key="input.key"
          >
            <div class="flex mb-4" v-if="show(input)">
              <span class="mr-2 w-32 truncate mt-4" :title="input.label">
                {{ input.label }}
              </span>
              <component
                v-if="extraEditorMap.hasOwnProperty(input.editor_type)"
                :is="extraEditorMap[input.editor_type]"
                :parameter="parameterData.parameter"
                @update="updateInnerValue($event.key, $event.value)"
              />
              <component
                v-else
                :is="findComponent(input.editor_type)"
                :value="
                  parameterData.parameter?.hasOwnProperty(input.key)
                    ? parameterData.parameter[input.key]
                    : input.default_value
                "
                :input="input"
                @update="updateInnerValue(input.key, $event)"
              />
              <a
                title="提示"
                target="_blank"
                rel="noopener"
                class="my-auto p-4"
              >
                <v-remixicon name="riInformationLine" size="18" />
              </a>
            </div>
          </template>

          <Divider
            v-if="props.parameter_define?.outputs?.length > 0"
            align="left"
            type="solid"
          >
            <b>输出参数</b>
          </Divider>
          <template
            v-for="output in props.parameter_define?.outputs"
            :key="output.key"
          >
            <div v-if="show(output)" class="flex mb-4">
              <span class="mr-2 w-32 truncate my-auto" :title="output.label">
                {{ output.label }}
              </span>
              <component
                :is="findComponent(output.editor_type)"
                :value="
                  parameterData.parameter?.hasOwnProperty(output.key)
                    ? parameterData.parameter[output.key]
                    : output.default_value
                "
                :input="output"
                @update="updateInnerValue(output.key, $event)"
              />
              <a
                title="提示"
                target="_blank"
                rel="noopener"
                class="my-auto p-4"
              >
                <v-remixicon name="riInformationLine" size="18" />
              </a>
            </div>
          </template>
        </div>
      </TabPanel>
      <TabPanel
        v-if="
          props.parameter_define?.extra &&
          props.parameter_define?.extra.length > 0
        "
        header="高级设置"
      >
        <div class="px-2">
          <div
            class="flex my-4"
            v-for="input in props.parameter_define?.extra"
            :key="input.key"
          >
            <span class="mr-2 w-32 truncate my-auto" :title="input.label">
              {{ input.label }}
            </span>
            <component
              :is="findComponent(input.editor_type)"
              :value="
                parameterData.parameter?.hasOwnProperty(input.key)
                  ? parameterData.parameter[input.key]
                  : input.default_value
              "
              :input="input"
              @update="updateInnerValue(input.key, $event)"
            />
            <a title="提示" target="_blank" rel="noopener" class="my-auto p-4">
              <v-remixicon name="riInformationLine" size="18" />
            </a>
          </div>
        </div>
      </TabPanel>
      <TabPanel header="错误处理" v-if="props.parameter_define?.error">
        <div class="px-2">
          <div class="flex my-4">
            <span class="mr-2 w-32 truncate my-auto"> 异常处理策略 </span>
            <Dropdown
              v-model="parameterData.parameter['exception']"
              :options="toDoTypes"
              optionLabel="name"
              option-value="value"
              class="mt-1 w-full"
            />
            <a title="提示" target="_blank" rel="noopener" class="my-auto p-4">
              <v-remixicon name="riInformationLine" size="18" />
            </a>
          </div>
          <div class="inline-flex">
            <InputSwitch
              v-model="parameterData.parameter['retry']"
              trueValue="True"
              falseValue="False"
            />
            <span class="ml-2"> 重试操作 </span>
          </div>
          <div v-if="parameterData.parameter['retry'] === 'True'" class="mr-2">
            <div class="inline-flex items-center">
              <span> 次数 </span>
              <v-remixicon
                title="重试操作的次数"
                name="riInformationLine"
                size="20"
                class="mr-2"
              />
              <InputText v-model="parameterData.parameter['retry_count']" />
            </div>
            <div class="inline-flex items-center ml-2">
              <span> 间隔(s) </span>
              <v-remixicon
                title="每次尝试等待的时间"
                name="riInformationLine"
                size="20"
                class="mr-2"
              />
              <InputText
                v-model="parameterData.parameter['retry_interval']"
                @keydown="(e) => e.stopPropagation()"
              />
            </div>
          </div>
        </div>
      </TabPanel>
    </TabView>
    <template #footer>
      <Button label="取消" icon="pi pi-times" @click="$emit('hide')" text />
      <Button
        v-if="props.runnable"
        label="运行"
        icon="pi pi-play"
        @click="$emit('run')"
      />
      <Button label="确认" icon="pi pi-check" @click="updateData" />
    </template>
  </Dialog>
</template>
<script setup>
import { getIconPath } from "@/utils/helper";
import InputSwitch from "primevue/inputswitch";
import Dropdown from "primevue/dropdown";
import Divider from "primevue/divider";
import Dialog from "primevue/dialog";
import InputText from "primevue/inputtext";
import Button from "primevue/button";
import TabView from "primevue/tabview";
import TabPanel from "primevue/tabpanel";
import { ref, watch } from "vue";

const props = defineProps({
  icon: {
    type: String,
    default: "",
  },
  color: {
    type: String,
    default:
      "bg-green-200 dark:bg-green-300 fill-green-200 dark:fill-green-300",
  },
  dialogShow: {
    type: Boolean,
    default: false,
  },
  parameter_define: {
    type: Object,
    default: () => ({}),
  },
  nodeData: {
    type: Object,
    default: () => ({}),
  },
  runnable: {
    type: Boolean,
    default: false,
  },
});
const emit = defineEmits(["hide", "update", "run"]);

const visible = ref(false);
watch(
  () => props.dialogShow,
  (value, oldValue) => {
    visible.value = value;
  },
  { immediate: true },
);
const extraEditorComponents = import.meta.glob(
  "@/components/editor/extra/*.vue",
  {
    eager: true,
  },
);
const extraEditorMap = {};
for (let each in extraEditorComponents) {
  const name = extraEditorComponents[each].default.__name;
  extraEditorMap[`${name}`] = extraEditorComponents[each].default;
}

const editorComponents = import.meta.glob("@/components/editor/*.vue", {
  eager: true,
});
const editorMap = {};
for (let each in editorComponents) {
  const name = editorComponents[each].default.__name;
  editorMap[`${name}`] = editorComponents[each].default;
}

function findComponent(name) {
  if (editorMap[name]) {
    return editorMap[name];
  }
  return editorMap["ExpressionTextInput"];
}

function show(input) {
  if (!input.show_if || input.show_if.length === 0) {
    return true;
  }
  return input.show_if.every((si) => {
    const [key, value] = si.split("=");
    if (key in parameterData.value.parameter) {
      const v = parameterData.value.parameter[key];
      return v.substring(2) === value;
    } else {
      const pd = props.parameter_define.inputs.filter((pd) => pd.key === key);
      if (pd.length > 0 && pd[0].default_value.length >= 2) {
        return pd[0].default_value.substring(2) === value;
      }
      return false;
    }
  });
}

const toDoTypes = [
  { value: "error", name: "抛出错误" },
  { value: "continue", name: "继续执行" },
];
const parameterData = ref({
  label: "",
  parameter: {},
});

function setData() {
  parameterData.value.label = props.nodeData.label;
  parameterData.value.parameter = props.nodeData.parameter;
  if (!parameterData.value.parameter.exception) {
    parameterData.value.parameter.exception = "error";
  }
}

function updateInnerValue(key, value) {
  parameterData.value.parameter[key] = value;
}

function updateData() {
  emit("update", { ...parameterData.value });
}
</script>
