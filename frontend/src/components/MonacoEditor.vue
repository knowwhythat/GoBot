<template>
  <div class="h-full" ref="dom"></div>
</template>

<script setup>
import { onBeforeUnmount, onMounted, ref, watch } from "vue";
import * as monaco from "monaco-editor";
import EditorWorker from "monaco-editor/esm/vs/editor/editor.worker?worker";
import "monaco-editor/esm/vs/basic-languages/python/python.contribution";

self.MonacoEnvironment = {
  getWorker(workerId, label) {
    return new EditorWorker();
  },
};

const props = defineProps({
  modelValue: String,
});

const emit = defineEmits(["update:modelValue"]);
const dom = ref();
let instance;

onMounted(() => {
  const pythonModel = monaco.editor.createModel(props.modelValue, "python");
  instance = monaco.editor.create(dom.value, {
    model: pythonModel, // theme: "vs-dark",
    tabSize: 2,
    automaticLayout: true,
    scrollBeyondLastLine: false,
  });

  instance.onDidChangeModelContent(() => {
    const value = instance.getValue();
    emit("update:modelValue", value);
  });

  instance.onKeyDown((listener, thisArgs, disposables) => {
    if (
      (listener.ctrlKey && listener.keyCode === monaco.KeyCode.KeyX) ||
      (listener.ctrlKey && listener.keyCode === monaco.KeyCode.KeyC) ||
      (listener.ctrlKey && listener.keyCode === monaco.KeyCode.KeyV)
    ) {
      listener.stopPropagation();
    }
  });
});

watch(
  () => props.code,
  (now, old) => {
    setEditorValue(now);
  },
);

const setEditorValue = (val) => {
  instance.pushUndoStop();

  instance.executeEdits("name-of-edit", [
    {
      range: instance.getModel().getFullModelRange(), // full range
      text: val, // target value here
    },
  ]);

  instance.pushUndoStop();
};

onBeforeUnmount(() => {
  instance?.dispose();
});
</script>
