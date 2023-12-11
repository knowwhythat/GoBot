<template>
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
        filterPlaceholder="搜索变量"
        @node-select="selectedNode($event)"
      ></Tree>
    </OverlayPanel>
    <InputText :placeholder="input.placeholder" v-model="content" />
    <Button
      :severity="isExpression ? 'Primary' : 'secondary'"
      @click="changeValueType"
    >
      <template #icon>
        <v-remixicon name="riReactjsFill" size="24" />
      </template>
    </Button>
  </InputGroup>
</template>
<script setup>
import Tree from "primevue/tree";
import InputGroup from "primevue/inputgroup";
import InputText from "primevue/inputtext";
import OverlayPanel from "primevue/overlaypanel";
import Button from "primevue/button";
import { ref, computed, inject } from "vue";
const props = defineProps({
  value: {
    type: String,
    default: "",
  },
  input: {
    type: Object,
    default: () => ({}),
  },
});

const emit = defineEmits(["update"]);
const isExpression = ref(false);

const content = computed({
  get() {
    if (props.value.split(":").length > 1) {
      isExpression.value = props.value.split(":")[0] === "1";
      return props.value.substring(2);
    } else {
      isExpression.value = false;
      return props.value;
    }
  },
  set(value) {
    emit("update", (isExpression.value ? "1:" : "0:") + value);
  },
});

function changeValueType() {
  isExpression.value = !isExpression.value;
  if (isExpression.value) {
    emit("update", "1:" + content.value);
  } else {
    emit("update", "0:" + content.value);
  }
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
