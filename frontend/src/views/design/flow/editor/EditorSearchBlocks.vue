<template>
  <div class="bg-white dark:bg-gray-800 ml-2 inline-flex items-center rounded-lg">
    <Button class="hoverable p-2 rounded-lg rounded-lg" icon @click="toggleActiveSearch">
      <v-remixicon name="riSearch2Line" />
    </Button>
    <AutoComplete ref="autocompleteEl" v-model="state.query" :suggestions="state.autocompleteItems"
      @complete="searchNodes" placeholder="搜索">
      <template #option="slotProps">
        <div class="flex-1 overflow-hidden">
          <p class="text-overflow">
            {{ slotProps.option.name }}
          </p>
          <p class="text-sm text-overflow leading-none text-gray-600 dark:text-gray-300">
            {{ slotProps.option.description }}
          </p>
        </div>
      </template>
    </AutoComplete>
  </div>
</template>
<script setup>
import { reactive, ref } from 'vue';
import AutoComplete from 'primevue/autocomplete';
import Button from 'primevue/button'

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

const autocompleteEl = ref(null);
const state = reactive({
  query: '',
  active: false,
  selected: false,
  autocompleteItems: [],
});


function searchNodes({ item, text }) {
  const query = text.toLocaleLowerCase();

  return (
    item.name.toLocaleLowerCase().includes(query) ||
    item.description.toLocaleLowerCase().includes(query)
  );
}
function toggleActiveSearch() {
  state.active = !state.active;

  if (state.active) {
    document.querySelector('#search-blocks')?.focus();
  }
}
function extractBlocks() {
  const editorContainer = document.querySelector('.vue-flow');
  editorContainer.classList.add('add-transition');
  const { width, height } = editorContainer.getBoundingClientRect();

  initialState.rectX = width / 2;
  initialState.rectY = height / 2;
  initialState.position = props.editor.getTransform();

  state.autocompleteItems = props.editor.getNodes.value.map(
    ({ computedPosition, id, data, label }) => ({
      id,
      position: computedPosition,
      description: data.description || '',
      name: t(`workflow.blocks.${label}.name`),
    })
  );
}
function clearHighlightedNodes() {
  document.querySelectorAll('.search-select-node').forEach((el) => {
    el.classList.remove('search-select-node');
  });
}
function clearState() {
  if (!state.selected) {
    props.editor.setTransform(initialState.position);
  }

  state.query = '';
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

  autocompleteEl.value.state.showPopover = false;
  clearHighlightedNodes();

  setTimeout(() => {
    const editorContainer = document.querySelector('.vue-flow');
    editorContainer.classList.remove('add-transition');
  }, 500);
}
function blurInput() {
  document.querySelector('#search-blocks')?.blur();
}
function onSelectItem({ item }) {
  const { x, y } = item.position;
  const { rectX, rectY } = initialState;

  clearHighlightedNodes();
  document
    .querySelector(`[data-id="${item.id}"]`)
    ?.classList.add('search-select-node');

  props.editor.setTransform({
    zoom: 1,
    x: -(x - rectX),
    y: -(y - rectY),
  });
}
function onItemSelected(event) {
  state.selected = true;

  const node = props.editor.getNode.value(event.item.id);
  props.editor.addSelectedNodes([node]);

  onSelectItem(event);
  blurInput();
}
</script>
<style scoped>
input {
  transition: width 300ms ease;
}
</style>
<style>
.search-select-node>div {
  @apply ring-4;
}

.vue-flow.add-transition .vue-flow__transformationpane {
  transition: transform 250ms ease;
}
</style>
