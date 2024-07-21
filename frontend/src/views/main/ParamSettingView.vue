<template>
  <template v-for="input in props.params" :key="input.name">
    <div v-if="input.direction === 'In'" class="flex mb-4">
      <span class="mr-2 w-48 truncate my-auto text-left" :title="input.name">
        {{ input.name }}
      </span>
      <InputText
        class="flex-auto"
        v-if="input.type === 'str' || input.type === 'any'"
        v-model="modelValue[input.name]"
      />
      <Password
        class="flex-auto"
        v-else-if="input.type === 'password'"
        toggleMask
        v-model="modelValue[input.name]"
      />
      <FilePathInput
        class="flex-auto"
        v-else-if="input.type === 'filePath'"
        type="file"
        v-model="modelValue[input.name]"
      />
      <FilePathInput
        class="flex-auto"
        v-else-if="input.type === 'dirPath'"
        type="directory"
        v-model="modelValue[input.name]"
      />
      <InputNumber
        class="flex-auto"
        showButtons
        v-else-if="input.type === 'number'"
        v-model="modelValue[input.name]"
      />
      <InputSwitch
        v-else-if="input.type === 'bool'"
        v-model="modelValue[input.name]"
      />
      <Dropdown
        class="flex-auto"
        v-else-if="input.type === 'single'"
        v-model="modelValue[input.name]"
        :options="input.options"
        optionLabel="key"
        optionValue="value"
      />
      <MultiSelect
        class="flex-auto"
        v-else-if="input.type === 'multiple'"
        :model-value="modelValue[input.name].split(',')"
        @update:model-value="modelValue[input.name] = $event.join(',')"
        :options="input.options"
        optionLabel="key"
        optionValue="value"
      />
    </div>
  </template>
</template>

<script setup>
import InputText from "primevue/inputtext";
import Password from "primevue/password";
import InputNumber from "primevue/inputnumber";
import InputSwitch from "primevue/inputswitch";
import Dropdown from "primevue/dropdown";
import MultiSelect from "primevue/multiselect";
import FilePathInput from "@/components/common/FilePathInput.vue";

const props = defineProps({
  modelValue: {
    type: Object,
    default: () => ({}),
  },
  params: {
    type: Array,
    default: () => [],
  },
});

const emit = defineEmits(["update:modelValue"]);
</script>
<style scoped>
:deep(.p-inputgroup) {
  width: auto;
}
:deep(.p-password-input) {
  width: 100%;
}
</style>
