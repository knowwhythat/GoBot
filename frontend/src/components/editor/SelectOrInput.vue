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
    <InputText
      :editable="true"
      :model-value="props.value"
      @update:modelValue="$emit('update', $event)"
      @keydown="(e) => e.stopPropagation()"
    ></InputText>
  </InputGroup>
</template>
<script setup>
import Button from "primevue/button";
import Tree from "primevue/tree";
import OverlayPanel from "primevue/overlaypanel";
import InputGroup from "primevue/inputgroup";
import InputText from "primevue/inputtext";
import { inject, ref } from "vue";

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

const { contextVariable } = inject("contextVariable");
const op = ref();
const toggle = (event) => {
  op.value.toggle(event);
};

function selectedNode(node) {
  op.value.toggle(false);
  emit("update", node.label);
}
</script>
