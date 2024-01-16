<template>
  <InputGroup>
    <InputText v-model="content" />
    <Button
      :severity="isExpression ? 'Primary' : 'secondary'"
      @click="changeType"
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
import { ref, computed, inject } from "vue";
const props = defineProps({
  modelValue: {
    type: [String, Number],
    default: "",
  },
});

const emit = defineEmits(["update:modelValue"]);
const isExpression = ref(false);

function changeType() {
  isExpression.value = !isExpression.value;
  emit(
    "update:modelValue",
    isExpression.value ? `{{ ${content.value} }}` : content.value
  );
}

var regex = /{{ (.*?) }}/;
const content = computed({
  get() {
    if (regex.test(props.modelValue)) {
      isExpression.value = true;
      return regex.exec(props.modelValue)[1];
    } else {
      isExpression.value = false;
      return props.modelValue;
    }
  },
  set(value) {
    emit("update:modelValue", isExpression.value ? `{{ ${value} }}` : value);
  },
});
</script>
