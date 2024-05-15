<template>
  <div class="w-textarea" ref="wTextarea">
    <div
      class="w-textarea_input"
      ref="wTextareaContent"
      :id="contentId"
      @click="inputClick($event)"
      @focus="isLocked = true"
      @blur="isLocked = false"
      @keydown.delete="handleDelete($event)"
      @input="handleInput($event.target)"
    ></div>
    <div class="resize" @mousedown="startResize"></div>
  </div>
</template>
<script setup>
import { ref, onMounted, onUnmounted, nextTick } from "vue";
import { nanoid } from "nanoid";

const props = defineProps({
  minHeight: {
    type: Number,
    default: 30,
  },
  maxHeight: {
    type: Number,
    default: 30,
  },
});

const modelValue = defineModel({ type: String });

const wTextarea = ref(null);
const wTextareaContent = ref(null);
const contentId = `content${nanoid(16)}`;

const isLocked = ref(false);
const currentTagId = ref(null);

const savedRange = ref(null);

function initData(content) {
  // 初始化数据
  if (content) {
    console.log(content);
    const data = JSON.parse(content);
    console.log(data);
    data.forEach((text) => {
      if (text && text.length >= 3) {
        if (text.substring(0, 1) === "0") {
          wTextareaContent.value.appendChild(
            document.createTextNode(text.substring(2)),
          );
        } else {
          wTextareaContent.value.appendChild(createTag(text.substring(2)));
        }
      }
    });
  }
  // 每次光标变化的时候，保存 range
  document.addEventListener("selectionchange", selectHandler);
}

onUnmounted(() => {
  document.removeEventListener("selectionchange", selectHandler);
});

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
    savedRange.value.insertNode(node);
  } else {
    wTextareaContent.value.appendChild(node);
  }

  let target = wTextareaContent.value;
  handleInput(target);
}

function handleInput(target) {
  nextTick(() => {
    if (target.childNodes) {
      let data = [];
      target.childNodes.forEach((node) => {
        if (node.nodeType === 3) {
          data.push("0:" + node.textContent);
        } else {
          data.push("2:" + node.textContent);
        }
      });
      modelValue.value = JSON.stringify(data);
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
  isLocked.value = true;
  const TAG_NAME = e.target.nodeName;
  if (currentTagId.value) {
    let target = document.getElementById(currentTagId.value);
    target.className = "";
  }
  if (TAG_NAME === "wise".toUpperCase()) {
    if (currentTagId.value !== e.target.id) {
      currentTagId.value = e.target.id;
      e.target.className = "active";
    } else {
      currentTagId.value = null;
    }
  } else {
    currentTagId.value = null;
  }
}

function selectHandler(e) {
  // 监听选定文本的变动
  let sel = window.getSelection();
  let range = sel.rangeCount > 0 ? sel.getRangeAt(0) : null;

  if (
    range &&
    range.commonAncestorContainer.ownerDocument.activeElement.id === contentId
  ) {
    savedRange.value = range;
  }
}

const resizing = ref(false);
let original_height = 0;
let original_mouse_y = 0;
const startResize = (e) => {
  resizing.value = true;
  e.preventDefault();
  e.stopPropagation();
  original_height = wTextareaContent.value.getBoundingClientRect().height;
  original_mouse_y = e.pageY;
  document.addEventListener("mousemove", handleResize, true);
  document.addEventListener("mouseup", stopResize, true);
};

const handleResize = (evt) => {
  if (resizing.value) {
    let height = original_height + (evt.pageY - original_mouse_y);
    if (height < props.minHight) {
      height = props.minHight;
    }
    if (props.maxHeight > 0 && height > props.maxHeight) {
      height = props.maxHeight;
    }
    wTextareaContent.value.style.height = height + "px";
  }
};

const stopResize = () => {
  resizing.value = false;
  document.removeEventListener("mousemove", handleResize);
  document.removeEventListener("mouseup", stopResize);
};

defineExpose({
  addTag,
  initData,
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
    background: #dcdfe6;
  }
}
</style>

<style lang="scss" scoped>
$borderColor: #dcdfe6;
$bgColor: #f5f7fa;
$textColor: #666;

.w-textarea {
  width: 100%;
  box-sizing: border-box;
  border-radius: 4px;
  border: 1px solid $borderColor;
  margin: 1px;
  overflow: auto;
  position: relative;
  text-align: left;

  &_input {
    width: 100%;
    min-height: 56px;
    box-sizing: border-box;
    padding: 2px;
    word-break: break-word;
    // 允许编辑，禁止富文本
    -webkit-user-modify: read-write-plaintext-only !important;

    &:focus {
      outline: none;
    }
  }

  .resize {
    position: absolute;
    right: 0px;
    bottom: 0px;
    width: 0px;
    height: 0px;
    border-bottom: 8px dashed rgba(0, 0, 0, 0.5);
    border-left: 8px dashed transparent;
    cursor: nw-resize;
  }
}
</style>
