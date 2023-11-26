<template>
  <InputGroup>
    <Button
      :severity="data.isExpression ? 'Primary' : 'secondary'"
      @click="changeValueType"
    >
      <template #icon>
        <v-remixicon name="riReactjsFill" size="24" />
      </template>
    </Button>
    <InputText
      placeholder="Vote"
      :modelValue="data.value"
      @update:modelValue="updateValue($event)"
    />
  </InputGroup>
</template>
<script setup>
import InputGroup from "primevue/inputgroup";
import InputText from "primevue/inputtext";
import Button from "primevue/button";
import { reactive, watch } from "vue";
const props = defineProps({
  value: {
    type: String,
    default: "",
  },
});

const data = reactive({
  isExpression: false,
  value: "",
});

watch(
  () => props.value,
  (oldData, newData) => {
    if (newData) {
      if (newData.length > 2) {
        data.isExpression = newData.split(":")[0] === "1";
        data.value = newData.substring(2);
      } else {
        data.isExpression = false;
        data.value = newData;
      }
    }
  },
  { deep: true, immediate: true }
);

const emit = defineEmits(["update"]);

function changeValueType() {
  data.isExpression = !data.isExpression;
  updateValue();
}

function updateValue() {}
</script>
