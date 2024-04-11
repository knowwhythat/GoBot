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
        :value="windowsElement"
        selectionMode="single"
        :filter="true"
        filterPlaceholder="搜索元素"
        @node-select="selectedNode($event)"
      ></Tree>
    </OverlayPanel>
    <InputText :placeholder="input.placeholder" v-model="content" disabled />
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

const { id, windowsElement } = inject("windowsElement");

const content = computed({
  get() {
    if (props.value.split(":").length > 1) {
      const key = props.value.substring(2);

      return props.value.substring(2);
    } else {
      return props.value;
    }
  },
});

const op = ref();
const toggle = (event) => {
  op.value.toggle(event);
};

function selectedNode(node) {
  if (!node["children"]) {
    op.value.toggle(false);
    emit("update", "0:" + node.key);
  }
}
</script>
