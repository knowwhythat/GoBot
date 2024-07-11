<template>
  <div
    class="bg-white dark:bg-gray-800 ml-2 inline-flex items-center rounded-lg"
  >
    <Button class="hoverable p-3" @click="toggleActiveSearch">
      <v-remixicon name="riSearch2Line" />
    </Button>
    <AutoComplete
      v-if="state.active"
      inputId="search-blocks"
      v-model="state.query"
      :suggestions="state.autocompleteItems"
      :completeOnFocus="true"
      @item-select="onItemSelected"
      @complete="search"
      optionLabel="name"
      placeholder="搜索"
      :forceSelection="true"
    />
  </div>
</template>
<script setup>
import { reactive } from "vue";
import AutoComplete from "primevue/autocomplete";
import Button from "primevue/button";

const props = defineProps({
  editor: {
    type: Object,
    default: () => ({}),
  },
});

const initialState = {
  rectX: 0,
  rectY: 0,
  position: {
    x: 0,
    y: 0,
    zoom: 1,
  },
};

const state = reactive({
  query: "",
  active: false,
  selected: false,
  allItems: [],
  autocompleteItems: [],
});

function toggleActiveSearch() {
  state.active = !state.active;

  if (state.active) {
    extractBlocks();
    document.querySelector("#search-blocks")?.focus();
  } else {
    clearState();
  }
}

function extractBlocks() {
  const editorContainer = document.querySelector(".vue-flow");
  editorContainer.classList.add("add-transition");
  const { width, height } = editorContainer.getBoundingClientRect();

  initialState.rectX = width / 2;
  initialState.rectY = height / 2;
  initialState.position = props.editor.getTransform();

  state.allItems = props.editor.getNodes.value.map(
    ({ computedPosition, id, data }) => ({
      id,
      position: computedPosition,
      description: data.description || "",
      name: data.label,
    }),
  );
}

function search(event) {
  if (!event.query.trim().length) {
    state.autocompleteItems = [...state.allItems];
  } else {
    state.autocompleteItems = state.allItems.filter((item) => {
      return item.name.toLowerCase().startsWith(event.query.toLowerCase());
    });
  }
}

function clearHighlightedNodes() {
  document.querySelectorAll(".search-select-node").forEach((el) => {
    el.classList.remove("search-select-node");
  });
}

function clearState() {
  if (!state.selected) {
    props.editor.setTransform(initialState.position);
  }

  state.query = "";
  state.active = false;
  state.selected = false;

  Object.assign(initialState, {
    rectX: 0,
    rectY: 0,
    position: {
      x: 0,
      y: 0,
      zoom: 1,
    },
  });

  clearHighlightedNodes();

  setTimeout(() => {
    const editorContainer = document.querySelector(".vue-flow");
    editorContainer.classList.remove("add-transition");
  }, 500);
}

function onSelectItem(event) {
  const item = event.value;
  const { x, y } = item.position;
  const { rectX, rectY } = initialState;

  clearHighlightedNodes();
  document
    .querySelector(`[data-id="${item.id}"]`)
    ?.classList.add("search-select-node");

  props.editor.setTransform({
    zoom: 1,
    x: -(x - rectX),
    y: -(y - rectY),
  });
}

function onItemSelected(event) {
  if (!event.value.id) {
    return;
  }
  state.selected = true;

  const node = props.editor.getNode.value(event.value.id);
  props.editor.addSelectedNodes([node]);

  onSelectItem(event);
}
</script>
<style scoped>
input {
  transition: width 300ms ease;
}
</style>
<style>
.search-select-node > div {
  @apply ring-4;
}

.vue-flow.add-transition .vue-flow__transformationpane {
  transition: transform 250ms ease;
}
</style>
