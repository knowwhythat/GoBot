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
        :value="nodes"
        selectionMode="single"
        :filter="true"
        filterPlaceholder="搜索变量"
        @node-select="selectedNode($event)"
      ></Tree>
    </OverlayPanel>
    <InputText
      :editable="true"
      v-model="props.value"
      @update:modelValue="$emit('update', $event)"
    ></InputText>
  </InputGroup>
</template>
<script setup>
import Button from "primevue/button";
import Tree from "primevue/tree";
import OverlayPanel from "primevue/overlaypanel";
import InputGroup from "primevue/inputgroup";
import InputText from "primevue/inputtext";
import Dropdown from "primevue/dropdown";
import { reactive, onMounted, ref, watch } from "vue";
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

const nodes = ref([
  {
    key: "0-0",
    label: "Work",
    data: "Work Folder",
    icon: "pi pi-fw pi-cog",
    children: [
      {
        key: "0-0-0",
        label: "Expenses.doc",
        icon: "pi pi-fw pi-file",
        data: "Expenses Document",
      },
      {
        key: "0-0-1",
        label: "Resume.doc",
        icon: "pi pi-fw pi-file",
        data: "Resume Document",
      },
    ],
  },
  {
    key: "0-1",
    label: "Home",
    data: "Home Folder",
    icon: "pi pi-fw pi-home",
    children: [
      {
        key: "0-1-0",
        label: "Invoices.txt",
        icon: "pi pi-fw pi-file",
        data: "Invoices for this month",
      },
    ],
  },
]);
const op = ref();
const toggle = (event) => {
  op.value.toggle(event);
};

function selectedNode(node) {
  op.value.toggle(false);
  emit("update", node.label);
}
</script>
