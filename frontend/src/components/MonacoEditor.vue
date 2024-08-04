<template>
  <div ref="editor" class="h-full"></div>
</template>
<script setup>
import * as monaco from "monaco-editor";
import "monaco-editor/esm/vs/basic-languages/python/python.contribution";
import { onBeforeUnmount, onMounted, ref, watch } from "vue";

const props = defineProps({
  code: {
    type: String,
    default: "",
  },
});

const emit = defineEmits(["update"]);
const editor = ref(null);

self.MonacoEnvironment = {
  getWorker: function (workerId, label) {
    const getWorkerModule = (moduleUrl, label) => {
      return new Worker(self.MonacoEnvironment.getWorkerUrl(moduleUrl), {
        name: label,
        type: "module",
      });
    };

    return getWorkerModule(
      "/monaco-editor/esm/vs/editor/editor.worker?worker",
      label,
    );
  },
};

const loadMonacoEditor = async () => {
  return new Promise((resolve) => {
    resolve(monaco);
  });
};

let editorInstance;
onMounted(async () => {
  await loadMonacoEditor();
  editorInstance = monaco.editor.create(editor.value, {
    value: props.code,
    language: "python",
    theme: "vs",
    lineNumbers: "on",
    scrollBeyondLastLine: false,
    overviewRulerBorder: false,
    automaticLayout: true,
    folding: true,
    renderLineHighlight: "gutter",
    contextmenu: true,
    useTabStops: false,
    quickSuggestions: true,
    accessibilitySupport: "on",
    acceptSuggestionOnEnter: "on",
    autoClosingBrackets: "always",
  });

  editorInstance.onDidChangeModelContent(() => {
    emit("update", editorInstance.getValue());
  });

  editorInstance.onKeyDown((listener, thisArgs, disposables) => {
    console.log(listener.keyCode);
    if (listener.ctrlKey && listener.keyCode === monaco.KeyCode.KeyS) {
      //不拦截保存
    } else {
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
  editorInstance.pushUndoStop();

  editorInstance.executeEdits("name-of-edit", [
    {
      range: editorInstance.getModel().getFullModelRange(), // full range
      text: val, // target value here
    },
  ]);

  editorInstance.pushUndoStop();
};

onBeforeUnmount(() => {
  editorInstance?.dispose();
});
</script>

<style scoped></style>
