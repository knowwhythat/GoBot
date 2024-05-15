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
    <Textarea
      v-if="isExpression"
      :placeholder="input.placeholder"
      v-model="expressionContent"
      variant="filled"
      rows="1"
      cols="30"
      style="min-height: 56px; max-height: 200px"
      @keydown="(e) => e.stopPropagation()"
    />
    <TagTextarea
      v-else
      ref="tagTextarea"
      v-model="rawContent"
      :minHeight="56"
      :maxHeight="200"
    />
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
import Textarea from "primevue/textarea";
import InputGroup from "primevue/inputgroup";
import OverlayPanel from "primevue/overlaypanel";
import Button from "primevue/button";
import TagTextarea from "@/components/common/TagTextarea.vue";
import { ref, inject, onMounted, watch } from "vue";
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

const expressionContent = ref();
watch(expressionContent, (value, oldValue) => {
  emit("update", "1:" + value);
});
const rawContent = ref();

watch(rawContent, (value, oldValue) => {
  emit("update", "3:" + value);
});
const tagTextarea = ref(null);

onMounted(() => {
  if (props.value.split(":").length > 1) {
    if (props.value.split(":")[0] === "1") {
      isExpression.value = true;
      expressionContent.value = props.value.substring(2);
    } else if (props.value.split(":")[0] === "0") {
      isExpression.value = false;
      rawContent.value = props.value;
    } else if (props.value.split(":")[0] === "3") {
      isExpression.value = false;
      if (props.value.substring(2)) {
        rawContent.value = props.value.substring(2);
      }
    }
    tagTextarea.value.initData(rawContent.value);
  }
});

function changeValueType() {
  isExpression.value = !isExpression.value;
  if (isExpression.value) {
    emit("update", "1:" + expressionContent.value);
  } else {
    emit("update", "3:" + rawContent.value);
  }
}
const { contextVariable } = inject("contextVariable");
const op = ref();
const toggle = (event) => {
  op.value.toggle(event);
};

function selectedNode(node) {
  op.value.toggle(false);
  if (!isExpression.value) {
    tagTextarea.value.addTag(node.label);
  }
}
</script>
