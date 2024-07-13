<template>
  <InputGroup>
    <InputText v-model="content" @keydown="(e) => e.stopPropagation()" />
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
import InputGroup from "primevue/inputgroup";
import InputText from "primevue/inputtext";
import Button from "primevue/button";
import { computed, ref } from "vue";

const props = defineProps({
  modelValue: {
    type: String,
    default: "",
  },
});

const emit = defineEmits(["update:modelValue"]);
const isExpression = ref(false);

const content = computed({
  get() {
    if (props.modelValue.split(":").length > 1) {
      isExpression.value = props.modelValue.split(":")[0] === "1";
      return props.modelValue.substring(2);
    } else {
      isExpression.value = false;
      return props.modelValue;
    }
  },
  set(value) {
    emit("update:modelValue", (isExpression.value ? "1:" : "0:") + value);
  },
});

function changeValueType() {
  isExpression.value = !isExpression.value;
  if (isExpression.value) {
    emit("update:modelValue", "1:" + content.value);
  } else {
    emit("update:modelValue", "0:" + content.value);
  }
}
</script>
