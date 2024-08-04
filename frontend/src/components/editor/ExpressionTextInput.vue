<template>
  <div class="flex flex-row w-full">
    <Button
      class="flex-none"
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
    <div
      v-show="!isExpression"
      class="w-textarea max-h-14 overflow-hidden flex-1 border border-solid rounded border-slate-300 focus:border-violet-500"
      ref="wTextarea"
    >
      <div
        class="w-textarea_input p-4 focus:outline-none"
        ref="wTextareaContent"
        @click="inputClick($event)"
        @keydown="(e) => e.stopPropagation()"
        @keydown.delete="handleDelete($event)"
        @input="handleInput($event.target)"
        @mouseup="selectHandler($event)"
        @keyup="selectHandler($event)"
      ></div>
    </div>
    <InputText
      class="flex-1"
      ref="expressionInput"
      v-show="isExpression"
      :placeholder="input.placeholder"
      v-model="content"
      @keydown="(e) => e.stopPropagation()"
    />
    <Button
      class="flex-none"
      :severity="isExpression ? 'Primary' : 'secondary'"
      @click="changeValueType"
    >
      <template #icon>
        <v-remixicon name="riReactjsFill" size="24" />
      </template>
    </Button>
  </div>
</template>
<script setup>
import Tree from "primevue/tree";
import InputText from "primevue/inputtext";
import OverlayPanel from "primevue/overlaypanel";
import Button from "primevue/button";
import { inject, nextTick, onMounted, ref, watch } from "vue";
import { nanoid } from "nanoid";

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

const wTextarea = ref(null);
const wTextareaContent = ref(null);
const expressionInput = ref(null);

const currentTagId = ref(null);

const savedRange = ref(null);

const isExpression = ref(false);
const content = ref("");

watch(content, () => {
  emit("update", "1:" + content.value);
});
watch(
  () => props.value,
  () => {},
);
onMounted(() => {
  // 初始化数据
  if (props.value && props.value.length >= 2) {
    const prefix = props.value.substring(0, 2);
    if (prefix === "1:") {
      isExpression.value = true;
      content.value = props.value.substring(2);
    } else if (prefix === "0:") {
      isExpression.value = false;
      wTextareaContent.value.appendChild(
        document.createTextNode(props.value.substring(2)),
      );
    } else if (prefix === "2:") {
      isExpression.value = false;
      wTextareaContent.value.appendChild(createTag(props.value.substring(2)));
    } else {
      isExpression.value = false;
      const data = props.value.substring(2);
      if (data) {
        try {
          const dataArray = JSON.parse(data);
          dataArray.forEach((text) => {
            if (text && text.length >= 2) {
              if (text.substring(0, 2) === "0:") {
                wTextareaContent.value.appendChild(
                  document.createTextNode(text.substring(2)),
                );
              } else {
                wTextareaContent.value.appendChild(
                  createTag(text.substring(2)),
                );
              }
            }
          });
        } catch (err) {
          log.error("变量设置错误", err);
        }
      }
    }
  }
});

function changeValueType() {
  isExpression.value = !isExpression.value;
  if (isExpression.value) {
    emit("update", "1:" + content.value);
  } else {
    handleInput(wTextareaContent.value);
  }
}

function selectedNode(node) {
  if (isExpression.value) {
    content.value += node.name;
  } else {
    addTag(node.name);
  }
  op.value.toggle(false);
}

const { contextVariable } = inject("contextVariable");

const op = ref();
const toggle = (event) => {
  op.value.toggle(event);
};

function addTag(text) {
  insertNode(createTag(text));
}

function createTag(text) {
  // 创建模版标签
  let node = document.createElement("wise");
  node.innerText = text;
  // 添加id便于删除
  node.id = nanoid(16);
  return node;
}

function insertNode(node) {
  // 在内容中插入标签
  if (savedRange.value) {
    savedRange.value.deleteContents();
    savedRange.value.insertNode(node);
    savedRange.value.collapse();
  } else {
    wTextareaContent.value.appendChild(node);
  }

  let target = wTextareaContent.value;
  handleInput(target);
}

function handleInput(target) {
  nextTick(() => {
    if (target.childNodes) {
      if (target.childNodes.length === 0) {
        emit("update", "0:");
      } else if (target.childNodes.length === 1) {
        const node = target.childNodes[0];
        if (node.nodeType === 3) {
          emit("update", "0:" + node.textContent);
        } else {
          emit("update", "2:" + node.textContent);
        }
      } else {
        const dataArray = [];
        target.childNodes.forEach((node) => {
          if (node.nodeType === 3) {
            if (node.textContent) {
              dataArray.push("0:" + node.textContent);
            }
          } else {
            dataArray.push("2:" + node.textContent);
          }
        });
        emit("update", "3:" + JSON.stringify(dataArray));
      }
    }
  });
}

function handleDelete(e) {
  // 监听“删除”事件
  if (currentTagId.value) {
    // 若已选中模版标签，直接删除dom节点
    let t = document.getElementById(currentTagId.value);
    wTextareaContent.value.removeChild(t);
    currentTagId.value = null;
    // 阻止浏览器默认的删除事件，并手动更新数据
    e.preventDefault();
    handleInput(e.target);
  }
}

function inputClick(e) {
  const TAG_NAME = e.target.nodeName;
  if (currentTagId.value) {
    let target = document.getElementById(currentTagId.value);
    target.className = "";
    currentTagId.value = null;
  } else if (TAG_NAME === "wise".toLocaleUpperCase()) {
    currentTagId.value = e.target.id;
    e.target.className = "active";
  } else {
    currentTagId.value = null;
  }
}

function selectHandler(e) {
  // 监听选定文本的变动
  let sel = window.getSelection();
  savedRange.value = sel.rangeCount > 0 ? sel.getRangeAt(0) : null;
}

defineExpose({
  addTag,
});
</script>
<style lang="scss">
// 给标签默认样式，不可scoped
.w-textarea {
  wise {
    display: inline;
    color: #26a2ff;
    padding: 0 4px;
    white-space: nowrap;
    cursor: default;
    -webkit-user-modify: read-only !important;
    user-select: none;
    i {
      color: crimson;
      padding-left: 1px;
      cursor: pointer;
    }
  }

  .active {
    background: #e5e8f0;
  }
}
</style>

<style lang="scss" scoped>
.w-textarea {
  &_input {
    width: 100%;
    // 允许编辑，禁止富文本
    -webkit-user-modify: read-write-plaintext-only !important;
    :focus-visible {
      border: none;
    }
  }
}
</style>
