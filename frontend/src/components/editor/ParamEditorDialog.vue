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
        <p class="w-64 truncate">{{ nodeData.label }}</p>
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
              v-model="nodeData.label"
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
              <span class="mr-2 w-32 truncate my-auto" :title="input.label">
                {{ input.label }}
              </span>
              <component
                :is="findComponent(input.editor_type)"
                :value="
                  nodeData.parameter?.hasOwnProperty(input.key)
                    ? nodeData.parameter[input.key]
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

          <Divider align="left" type="solid">
            <b>输出参数</b>
          </Divider>
          <div
            class="flex mb-4"
            v-for="output in props.parameter_define?.outputs"
            :key="output.key"
          >
            <span class="mr-2 w-32 truncate my-auto" :title="output.label">
              {{ output.label }}
            </span>
            <component
              :is="findComponent(output.editor_type)"
              :value="
                nodeData.parameter?.hasOwnProperty(output.key)
                  ? nodeData.parameter[output.key]
                  : output.default_value
              "
              :input="output"
              @update="updateInnerValue(output.key, $event)"
            />
            <a title="提示" target="_blank" rel="noopener" class="my-auto p-4">
              <v-remixicon name="riInformationLine" size="18" />
            </a>
          </div>
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
                nodeData.parameter?.hasOwnProperty(input.key)
                  ? nodeData.parameter[input.key]
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
              v-model="nodeData.parameter['exception']"
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
              v-model="nodeData.parameter['retry']"
              trueValue="True"
              falseValue="False"
            />
            <span class="ml-2"> 重试操作 </span>
          </div>
          <div v-if="nodeData.parameter['retry'] == 'True'" class="mr-2">
            <div class="inline-flex items-center">
              <span> 次数 </span>
              <v-remixicon
                title="重试操作的次数"
                name="riInformationLine"
                size="20"
                class="mr-2"
              />
              <InputText v-model="nodeData.parameter['retry_count']" />
            </div>
            <div class="inline-flex items-center ml-2">
              <span> 间隔(s) </span>
              <v-remixicon
                title="每次尝试等待的时间"
                name="riInformationLine"
                size="20"
                class="mr-2"
              />
              <InputText v-model="nodeData.parameter['retry_interval']" />
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
import { reactive, watch, ref, Comment, computed } from "vue";
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
  { immediate: true }
);

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
  if (!input.show_if || input.show_if.length == 0) {
    return true;
  }
  const result = input.show_if.every((si) => {
    const [key, value] = si.split("=");
    if (key in nodeData.parameter) {
      const v = nodeData.parameter[key];
      return v.substring(2) === value;
    } else {
      return (
        props.parameter_define.inputs
          .filter((pd) => pd.key === key)[0]
          .default_value.substring(2) === value
      );
    }
  });
  return result;
}
const toDoTypes = [
  { value: "error", name: "抛出错误" },
  { value: "continue", name: "继续执行" },
];
const nodeData = reactive({
  label: "",
  parameter: {},
});

function setData() {
  nodeData.label = props.nodeData.label;
  nodeData.parameter = props.nodeData.parameter;
}

function updateInnerValue(key, value) {
  console.log(key, value);
  nodeData.parameter[key] = value;
}

function updateData() {
  emit("update", { ...nodeData });
}
</script>
