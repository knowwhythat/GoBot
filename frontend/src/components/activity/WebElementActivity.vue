<template>
  <ActivityBase
    :toggleable="true"
    :id="props.element.id"
    :breakpoint="props.element.breakpoint"
    :collapsed="props.element.collapsed"
    :deleteable="props.element.deleteable"
    :label="props.element.label"
    :icon="props.element.icon_path"
    :color="props.element.color"
    :parameter_define="props.element.parameter_define"
    :parameter="props.element.parameter"
    @delete="$emit('delete', { id: props.element.id })"
    @update="updateData($event)"
  >
    <InputGroup class="mb-2">
      <InputText
        :model-value="frameSelector"
        @update:model-value="chageFrameSelector($event)"
        placeholder="Frame选择器"
      />
      <Button
        :severity="frameSelectorExpression ? 'Primary' : 'secondary'"
        @click="chageFrameSelectorExpression"
      >
        <template #icon>
          <v-remixicon name="riReactjsFill" size="24" />
        </template>
      </Button>
    </InputGroup>
    <InputGroup class="mb-2">
      <InputText
        :model-value="elementSelector"
        @update:model-value="chageElementSelector($event)"
        placeholder="元素选择器"
      />
      <Button
        :severity="elementSelectorExpression ? 'Primary' : 'secondary'"
        @click="chageElementSelectorExpression"
      >
        <template #icon>
          <v-remixicon name="riReactjsFill" size="24" />
        </template>
      </Button>
    </InputGroup>
    <div class="flex gap-4 justify-center">
      <Button severity="Primary" label="拾取" @click="pickElement">
        <template #icon>
          <v-remixicon name="riCursorLine" size="24" />
        </template>
      </Button>
      <Button
        icon="pi pi-arrow-right-arrow-left"
        severity="help"
        label="校验"
        @click="checkElement"
      />
    </div>
  </ActivityBase>
</template>
<script setup>
import { inject, ref, watch } from "vue";
import InputGroup from "primevue/inputgroup";
import InputText from "primevue/inputtext";
import Button from "primevue/button";
import ActivityBase from "./ActivityBase.vue";
const props = defineProps({
  element: {
    type: Object,
    default: {},
  },
});
const emit = defineEmits(["delete"]);

const frameSelectorExpression = ref(false);
const frameSelector = ref("");
const elementSelectorExpression = ref(false);
const elementSelector = ref("");
watch(
  () => props.element.parameter,
  (value, oldValue) => {
    if (
      "frame_selector" in value &&
      value["frame_selector"].split(":").length > 1
    ) {
      frameSelectorExpression.value =
        value["frame_selector"].split(":")[0] === "1";
      frameSelector.value = value["frame_selector"].substring(2);
    } else {
      const defaultValue = props.element.parameter_define.inputs.filter(
        (pd) => pd.key === "frame_selector"
      )[0].default_value;
      frameSelectorExpression.value = defaultValue.split(":")[0] === "1";
      frameSelector.value = defaultValue.substring(2);
    }
    if (
      "element_selector" in value &&
      value["element_selector"].split(":").length > 1
    ) {
      elementSelectorExpression.value =
        value["element_selector"].split(":")[0] === "1";
      elementSelector.value = value["element_selector"].substring(2);
    } else {
      const defaultValue = props.element.parameter_define.inputs.filter(
        (pd) => pd.key === "element_selector"
      )[0].default_value;
      elementSelectorExpression.value = defaultValue.split(":")[0] === "1";
      elementSelector.value = defaultValue.substring(2);
    }
  },
  { immediate: true, deep: true }
);
const { dataChanged, updateDataChanged } = inject("dataChanged");
function chageFrameSelector(data) {
  if (frameSelectorExpression.value) {
    props.element.parameter["frame_selector"] = "1:" + data;
  } else {
    props.element.parameter["frame_selector"] = "0:" + data;
  }
  updateDataChanged();
}
function chageFrameSelectorExpression() {
  frameSelectorExpression.value = !frameSelectorExpression.value;
  if (frameSelectorExpression.value) {
    props.element.parameter["frame_selector"] = "1:" + frameSelector.value;
  } else {
    props.element.parameter["frame_selector"] = "0:" + frameSelector.value;
  }
  updateDataChanged();
}

function chageElementSelector(data) {
  if (elementSelectorExpression.value) {
    props.element.parameter["element_selector"] = "1:" + data;
  } else {
    props.element.parameter["element_selector"] = "0:" + data;
  }
  updateDataChanged();
}
function chageElementSelectorExpression() {
  elementSelectorExpression.value = !elementSelectorExpression.value;
  if (elementSelectorExpression.value) {
    props.element.parameter["element_selector"] = "1:" + elementSelector.value;
  } else {
    props.element.parameter["element_selector"] = "0:" + elementSelector.value;
  }
  updateDataChanged();
}

function updateData(data) {
  for (const key in data) {
    if (Object.hasOwnProperty.call(data, key)) {
      props.element[key] = data[key];
    }
  }
}

function pickElement() {}

function checkElement() {}
</script>
