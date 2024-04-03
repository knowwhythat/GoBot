<template>
  <div class="flex flex-col gap-2 w-full">
    <div class="flex flex-row gap-2 w-full">
      <InputGroup>
        <Button
          icon="pi pi-chevron-right"
          severity="Primary"
          @click="toggle"
          aria-haspopup="true"
          aria-controls="overlay_panel"
        />
        <OverlayPanel ref="op" appendTo="body" :unstyled="true">
          <Tree
            :value="contextVariable"
            selectionMode="single"
            :filter="true"
            filterPlaceholder="搜索"
            @node-select="selectedNode($event)"
          ></Tree>
        </OverlayPanel>
        <InputText
          placeholder="frame选择器"
          v-model="frameSelector"
          @keydown="(e) => e.stopPropagation()"
        />
        <Button
          :severity="frameIsExpression ? 'Primary' : 'secondary'"
          @click="changeFrameValueType"
        >
          <template #icon>
            <v-remixicon name="riReactjsFill" size="24" />
          </template>
        </Button>
      </InputGroup>
      <InputGroup>
        <Button
          icon="pi pi-chevron-right"
          severity="Primary"
          @click="toggle"
          aria-haspopup="true"
          aria-controls="overlay_panel"
        />
        <OverlayPanel ref="op" appendTo="body" :unstyled="true">
          <Tree
            :value="contextVariable"
            selectionMode="single"
            :filter="true"
            filterPlaceholder="搜索"
            @node-select="selectedNode($event)"
          ></Tree>
        </OverlayPanel>
        <InputText
          placeholder="元素选择器"
          v-model="elementSelector"
          @keydown="(e) => e.stopPropagation()"
        />
        <Button
          :severity="elementIsExpression ? 'Primary' : 'secondary'"
          @click="changeElementValueType"
        >
          <template #icon>
            <v-remixicon name="riReactjsFill" size="24" />
          </template>
        </Button>
      </InputGroup>
    </div>
    <div class="flex flex-row gap-2 justify-center w-full">
      <Button icon="pi pi-sun" @click="pickElement" label="拾取" />
      <Button
        icon="pi pi-arrow-right-arrow-left"
        @click="checkElement"
        label="校验"
      />
    </div>
  </div>
</template>
<script setup>
import Tree from "primevue/tree";
import InputGroup from "primevue/inputgroup";
import InputText from "primevue/inputtext";
import OverlayPanel from "primevue/overlaypanel";
import Button from "primevue/button";
import { ref, computed, inject, onMounted } from "vue";
import { StartCheck, StartPick } from "@back/go/backend/App";
import { useToast } from "primevue/usetoast";
const toast = useToast();
const props = defineProps({
  parameter: {
    type: Object,
    default: () => ({}),
  },
});

const emit = defineEmits(["update"]);

const frameIsExpression = ref(false);

const frameSelector = computed({
  get() {
    if (
      props.parameter["frame_selector"] &&
      props.parameter["frame_selector"].split(":").length > 1
    ) {
      frameIsExpression.value =
        props.parameter["frame_selector"].split(":")[0] === "1";
      return props.parameter["frame_selector"].substring(2);
    } else {
      frameIsExpression.value = false;
      return "";
    }
  },
  set(value) {
    emit("update", {
      key: "frame_selector",
      value: (frameIsExpression.value ? "1:" : "0:") + value,
    });
  },
});
function changeFrameValueType() {
  frameIsExpression.value = !frameIsExpression.value;
  emit("update", {
    key: "frame_selector",
    value: (frameIsExpression.value ? "1:" : "0:") + frameSelector.value,
  });
}

const elementIsExpression = ref(false);

const elementSelector = computed({
  get() {
    if (
      props.parameter["element_selector"] &&
      props.parameter["element_selector"].split(":").length > 1
    ) {
      elementIsExpression.value =
        props.parameter["element_selector"].split(":")[0] === "1";
      return props.parameter["element_selector"].substring(2);
    } else {
      elementIsExpression.value = false;
      return "";
    }
  },
  set(value) {
    emit("update", {
      key: "element_selector",
      value: (elementIsExpression.value ? "1:" : "0:") + value,
    });
  },
});
function changeElementValueType() {
  elementIsExpression.value = !elementIsExpression.value;
  emit("update", {
    key: "element_selector",
    value: (elementIsExpression.value ? "1:" : "0:") + elementSelector.value,
  });
}

async function checkElement() {
  try {
    if (!frameSelector.value) {
      toast.add({
        severity: "error",
        summary: "校验失败",
        detail: "元素选择器不能为空",
        life: 3000,
      });
      return;
    }
    const resp = await StartCheck(frameSelector.value, elementSelector.value);
    const result = JSON.parse(resp);
    if (result["message"] === "ok") {
      toast.add({
        severity: "success",
        summary: "校验成功",
        detail: "找到元素",
        life: 3000,
      });
    } else {
      toast.add({
        severity: "error",
        summary: "校验失败",
        detail: "元素未找到",
        life: 3000,
      });
    }
  } catch (err) {
    toast.add({
      severity: "error",
      summary: "校验失败",
      detail: err,
      life: 3000,
    });
  }
}

async function pickElement() {
  const result = await StartPick();
  const resp = JSON.parse(result);
  frames.value = resp["message"]["framePath"];
  selectors.value = resp["message"]["elementPath"];
  if (frames.value.length > 0) {
    selectedFrame.value = frames.value[0];
  }
  if (selectors.value.length > 0) {
    selectedSelector.value = selectors.value[0];
  }
  dialogVisible.value = true;
}

const { contextVariable } = inject("contextVariable");

const op = ref();
const toggle = (event) => {
  op.value.toggle(event);
};

function selectedNode(node) {
  op.value.toggle(false);
  isExpression.value = true;
  emit("update", "1:" + node.label);
}
</script>
