<template>
    <vue-flow :id="props.id" style="height: 100vh;" :class="{ disabled: isDisabled }" @drop="onDropInEditor"
        @dragend="clearHighlightedElements" @dragover.prevent="onDragoverEditor" :default-edge-options="{
            type: 'custom',
            updatable: true,
            selectable: true,
            markerEnd: settings.arrow ? MarkerType.ArrowClosed : '',
        }">
        <Background />
        <MiniMap v-if="minimap" :node-class-name="minimapNodeClassName" class="hidden md:block" />
        <div
            class="flex flex-col gap-2 p-2 justify-center absolute z-10 inset-y-20 left-4 w-40 rounded-xl bg-gray-200 dark:bg-gray-800 overflow-auto">
            <div v-for="block in blocks" :key="block.id" :title="getBlockTitle(block)" draggable="true"
                class="transform select-none cursor-move relative p-4 rounded-lg bg-input transition group bg-white"
                @dragstart="$event.dataTransfer.setData('block', JSON.stringify(block))">
                <div
                    class="flex items-center absolute right-2 invisible group-hover:visible top-2 text-gray-600 dark:text-gray-300">
                    <a :href="`https://docs.automa.site/blocks/${block.id}.html`" title="帮助" target="_blank" rel="noopener">
                        <v-remixicon name="riInformationLine" size="18" />
                    </a>
                </div>
                <v-remixicon :path="getIconPath(block.icon)" :name="block.icon" size="24" class="mb-2" />
                <p class="leading-tight text-overflow capitalize">
                    {{ block.name }}
                </p>
            </div>
        </div>
        <div v-if="editorControls" class="flex items-center absolute w-full p-4 left-0 bottom-0 z-10 md:pr-60">
            <slot name="controls-prepend" />
            <editor-search-blocks :editor="editor" />
            <div class="flex-grow pointer-events-none" />
            <slot name="controls-append" />

            <div class="rounded-lg bg-white dark:bg-gray-800 inline-block">
                <Button v-tooltip="'重置'" class="control-button p-2 rounded-lg" @click="editor.fitView()">
                    <v-remixicon name="riFullscreenLine" />
                </Button>
                <Button v-tooltip="'放大'" class="p-2 rounded-lg relative z-10" @click="editor.zoomOut()">
                    <v-remixicon name="riSubtractLine" />
                </Button>
                <hr class="h-6 border-r inline-block" />
                <Button v-tooltip.top="'缩小'" class="p-2 rounded-lg" @click="editor.zoomIn()">
                    <v-remixicon name="riAddLine" />
                </Button>
            </div>
        </div>
        <template v-for="(node, name) in nodeTypes" :key="name" #[name]="nodeProps">
            <component :is="node" v-bind="{ ...nodeProps, editor: editor, }" @delete="deleteBlock"
                @settings="initEditBlockSettings" @edit="editBlock(nodeProps, $event)"
                @update="updateBlockData(nodeProps.id, $event)" />
        </template>
        <template #edge-custom="edgeProps">
            <editor-custom-edge v-bind="edgeProps" />
        </template>
    </vue-flow>
</template>
<script setup>
import { onMounted, onBeforeUnmount, watch, computed, reactive, markRaw } from 'vue';
import {
    VueFlow,
    useVueFlow,
    MarkerType,
    getConnectedEdges,
} from '@vue-flow/core';
import { Background, MiniMap } from '@vue-flow/additional-components';
import cloneDeep from 'lodash.clonedeep';
import { customAlphabet } from 'nanoid';
import { debounce, parseJSON, throttle } from '@/utils/helper';
import { useAppStore } from '@/stores/app';
import { getBlocks } from '@/utils/getSharedData';
import { useEditorBlock } from '@/composable/editorBlock';
import { categories } from '@/utils/shared';
import EditorCustomEdge from './editor/EditorCustomEdge.vue';
import EditorSearchBlocks from './editor/EditorSearchBlocks.vue';
import Button from 'primevue/button'

const nanoid = customAlphabet('1234567890abcdefghijklmnopqrstuvwxyz', 7);

const state = reactive({
    dataChanged: false,
})

const props = defineProps({
    id: {
        type: String,
        default: 'editor',
    },
    data: {
        type: Object,
        default: () => ({
            x: 0,
            y: 0,
            zoom: 0,
            nodes: [],
            edges: [],
        }),
    },
    options: {
        type: Object,
        default: () => ({}),
    },
    editorControls: {
        type: Boolean,
        default: true,
    },
    minimap: {
        type: Boolean,
        default: true,
    },
    disabled: Boolean,
});
const emit = defineEmits([
    'edit',
    'init',
    'update:node',
    'delete:node',
    'update:settings',
]);

const fallbackBlocks = {
    BlockBasic: ['BlockExportData'],
    BlockBasicWithFallback: ['BlockWebhook'],
};

const blockComponents = import.meta.glob('@/components/block/*.vue', { eager: true });
const nodeTypes = {}
console.log(blockComponents)
for (let each in blockComponents) {
    const name = blockComponents[each].default.__name
    nodeTypes[`node-${name}`] = blockComponents[each].default
}

const getPosition = (position) => (Array.isArray(position) ? position : [0, 0]);
const setMinValue = (num, min) => (num < min ? min : num);

const appStore = useAppStore();
const isMac = appStore.settings.operateSystem.indexOf('Mac') !== -1;

const editor = useVueFlow({
    id: props.id,
    edgeUpdaterRadius: 20,
    deleteKeyCode: 'Delete',
    elevateEdgesOnSelect: true,
    defaultZoom: props.data?.zoom ?? 1,
    minZoom: setMinValue(+appStore.settings.editor.minZoom || 0.5, 0.1),
    maxZoom: setMinValue(
        +appStore.settings.editor.maxZoom || 1.2,
        +appStore.settings.editor.minZoom + 0.1
    ),
    multiSelectionKeyCode: isMac ? 'Meta' : 'Control',
    defaultPosition: getPosition(props.data?.position),
    ...props.options,
});
editor.onConnect((params) => {
    params.class = `source-${params.sourceHandle} target-${params.targetHandle}`;
    params.updatable = true;
    editor.addEdges([params]);
});
editor.onEdgeUpdate(({ edge, connection }) => {
    const isBothOutput =
        connection.sourceHandle.includes('output') &&
        connection.targetHandle.includes('output');
    if (isBothOutput) return;

    Object.assign(edge, connection);
});

const blocks = getBlocks();
function getBlockTitle({ description, id, name }) {
    const blockPath = `workflow.blocks.${id}`;
    if (!blockPath) return blocksDetail[id].name;

    const descPath = `${blockPath}.${description ? 'description' : 'name'}`;
    let blockDescription = descPath ? descPath : name;

    if (description) {
        blockDescription = `${blockPath}.name\n${blockDescription}`;
    }

    return blockDescription;
}
function getIconPath(path) {
    if (path && path.startsWith('path')) {
        const { 1: iconPath } = path.split(':');
        return iconPath;
    }

    return '';
}
function toggleHighlightElement({ target, elClass, classes }) {
    const targetEl = target.closest(elClass);

    if (targetEl) {
        targetEl.classList.add(classes);
    } else {
        const elements = document.querySelectorAll(`.${classes}`);
        elements.forEach((element) => {
            element.classList.remove(classes);
        });
    }
}
function onDragoverEditor({ target }) {
    toggleHighlightElement({
        target,
        elClass: '.vue-flow__handle.source',
        classes: 'dropable-area__handle',
    });

    if (!target.closest('.vue-flow__handle')) {
        toggleHighlightElement({
            target,
            elClass: '.vue-flow__node:not(.vue-flow__node-BlockGroup)',
            classes: 'dropable-area__node',
        });
    }
}

function onDropInEditor({ dataTransfer, clientX, clientY, target }) {
    const savedBlocks = parseJSON(dataTransfer.getData('savedBlocks'), null);
    if (savedBlocks) {
        if (savedBlocks.settings.asBlock) {
            const position = editor.project({
                x: clientX,
                y: clientY - 18,
            });
            editor.addNodes([
                {
                    position,
                    id: nanoid(),
                    data: savedBlocks,
                    type: 'BlockPackage',
                    label: 'block-package',
                },
            ]);
        } else {
            const { nodes, edges } = savedBlocks.data;
            /* eslint-disable-next-line */
            const newElements = copyElements(nodes, edges, { clientX, clientY });

            editor.addNodes(newElements.nodes);
            editor.addEdges(newElements.edges);
        }

        state.dataChanged = true;
        return;
    }

    const block = parseJSON(dataTransfer.getData('block'), null);
    if (!block || block.fromBlockBasic) return;

    clearHighlightedElements();

    const isTriggerExists =
        block.id === 'trigger' &&
        editor.getNodes.value.some((node) => node.label === 'trigger');
    if (isTriggerExists) return;

    const position = editor.project({ x: clientX - 60, y: clientY - 18 });
    const nodeId = nanoid();
    const newNode = {
        position,
        label: block.name,
        data: block.data,
        type: block.component,
        id: nodeId,
    };
    editor.addNodes([newNode]);

    state.dataChanged = true;
}

function clearHighlightedElements() {
    const elements = document.querySelectorAll(
        '.dropable-area__node, .dropable-area__handle'
    );
    elements.forEach((element) => {
        element.classList.remove('dropable-area__node');
        element.classList.remove('dropable-area__handle');
    });
}

const settings = appStore.settings.editor;
const isDisabled = computed(() => props.options.disabled ?? props.disabled);

const blockSettingsState = reactive({
    show: false,
    data: {},
});

function initEditBlockSettings({ blockId, details, data, itemId }) {
    blockSettingsState.data = {
        itemId,
        blockId,
        id: details.id,
        data: cloneDeep(data),
    };
    blockSettingsState.show = true;
}
function clearBlockSettings() {
    Object.assign(blockSettingsState, {
        data: null,
        show: false,
    });
}
function minimapNodeClassName({ label }) {
    const block = useEditorBlock(label);
    return block.category.color;
}
function updateBlockData(nodeId, data = {}) {
    if (isDisabled.value) return;

    const node = editor.getNode.value(nodeId);
    node.data = { ...node.data, ...data };

    emit('update:node', node);
}
function updateBlockSettingsData(newSettings) {
    if (isDisabled.value) return;

    const nodeId = blockSettingsState.data.blockId;
    const node = editor.getNode.value(nodeId);

    if (blockSettingsState.data.itemId) {
        const index = node.data.blocks.findIndex(
            (item) => item.itemId === blockSettingsState.data.itemId
        );
        if (index === -1) return;

        node.data.blocks[index].data = {
            ...node.data.blocks[index].data,
            ...newSettings,
        };
    } else {
        node.data = { ...node.data, ...newSettings };
    }

    emit('update:settings', {
        settings: newSettings,
        itemId: blockSettingsState.data.itemId,
        blockId: blockSettingsState.data.blockId,
    });
}
function editBlock({ id, label, data }, additionalData = {}) {
    if (isDisabled.value) return;

    emit('edit', {
        id: label,
        blockId: id,
        data: cloneDeep(data),
        ...additionalData,
    });
}
function deleteBlock(nodeId) {
    if (isDisabled.value) return;

    editor.removeNodes([nodeId]);
    emit('delete:node', nodeId);
}
function onMousedown(event) {
    if (isDisabled.value && event.shiftKey) {
        event.stopPropagation();
        event.preventDefault();
    }
}
function applyFlowData() {
    if (settings.snapToGrid) {
        editor.snapToGrid.value = true;
        editor.snapGrid.value = Object.values(settings.snapGrid);
    }

    editor.setNodes(
        props.data?.nodes?.map((node) => ({ ...node, events: {} })) || []
    );
    editor.setEdges(props.data?.edges || []);
    editor.setTransform({
        x: props.data?.x || 0,
        y: props.data?.y || 0,
        zoom: props.data?.zoom || 1,
    });
}

watch(
    () => props.disabled,
    (value) => {
        const keys = [
            'nodesDraggable',
            'edgesUpdatable',
            'nodesConnectable',
            'elementsSelectable',
        ];

        keys.forEach((key) => {
            editor[key].value = !value;
        });
    },
    { immediate: true }
);
watch(editor.getSelectedNodes, (nodes, _, cleanup) => {
    const connectedEdges = getConnectedEdges(nodes, editor.getEdges.value);

    connectedEdges.forEach((edge) => {
        edge.class = 'connected-edges';
    });

    cleanup(() => {
        connectedEdges.forEach((edge) => {
            edge.class = undefined;
        });
    });
});

onMounted(() => {
    applyFlowData();
    window.addEventListener('mousedown', onMousedown, true);
    emit('init', editor);
});
onBeforeUnmount(() => {
    window.removeEventListener('mousedown', onMousedown, true);
});
</script>
<style>
@import '@vue-flow/core/dist/style.css';
@import '@vue-flow/core/dist/theme-default.css';

.control-button {
    @apply p-2 rounded-lg bg-white dark:bg-gray-800 transition-colors;
}
</style>
  