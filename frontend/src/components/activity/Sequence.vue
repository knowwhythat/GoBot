<template>
    <div class="flex flex-col min-w-[300px] bg-gray-50 border-1 border-gray-200 rounded-xl hover:ring-2 m-4">
        <div v-if="props.showLabel" class="bg-primary rounded-t-xl p-2 truncate">{{ props.label }}</div>
        <draggable :model-value="props.children" item-key="itemId"
            class="flex flex-col items-center p-x-4 mb-8 overflow-auto nowheel scroll text-sm" @mousedown.stop
            @dragover.prevent @drop="handleDrop" @update:modelValue="$emit('update', { children: $event })">
            <template #item="{ element, index }">
                <div class="activity-node p-2 rounded-lg bg-input flex items-center group" :group-item-id="element.itemId"
                    style="cursor: grab" :index="index" :activity-data="element" @dragstart="onDragStart(element, $event)"
                    @dragend="onDragEnd(element, $event)">
                    <div class="leading-tight flex-1 overflow-hidden">
                        <p class="text-overflow">
                            {{ element.label }}
                        </p>
                    </div>
                    <div class="invisible group-hover:visible">
                        <v-remixicon name="riPencilLine" size="18" class="cursor-pointer inline-block mr-2"
                            @click="editBlock(element)" />
                        <v-remixicon name="riSettings3Line" size="18" class="cursor-pointer inline-block mr-2"
                            @click="editItemSettings(element)" />
                        <v-remixicon name="riDeleteBin7Line" size="18" class="cursor-pointer inline-block"
                            @click="deleteItem(index, element.itemId)" />
                    </div>
                </div>
            </template>
            <template #footer>
                <div class="w-full p-2 rounded-lg text-gray-600 dark:text-gray-200 border text-center border-dashed">
                    拖放组件到此处
                </div>
            </template>
        </draggable>
    </div>
</template>
<script setup>
import { ref, onMounted, computed, shallowReactive, watch } from 'vue';
import draggable from 'vuedraggable';
import { nanoid } from 'nanoid';


const props = defineProps({
    id: {
        type: String,
        default: nanoid(8),
    },
    label: {
        type: String,
        default: '子模块',
    },
    showLabel: {
        type: Boolean,
        default: false
    },
    children: {
        type: Array,
        default: () => ([]),
    },
});

const emit = defineEmits([
    'update',
    'delete',
]);

function onDragStart(item, event) {
    event.dataTransfer.setData('activity', JSON.stringify({ ...item, from: "sequence", containerId: props.id }));
}

function onDragEnd(item, event) {
    const droppedBlock = JSON.parse(event.dataTransfer.getData('activity') || null);
    console.log(event)
    setTimeout(() => {
        const blockEl = document.querySelector(`[group-item-id="${item.itemId}"]`);

        if (blockEl) {
            const blockIndex = props.children.findIndex(
                (activity) => activity.itemId === item.itemId
            );

            if (blockIndex !== -1) {
                const copyActivities = [...props.children];
                copyActivities.splice(blockIndex, 1);
                emit('update', { children: copyActivities });
            }
        }
    }, 200);
}

function handleDrop(event) {
    event.preventDefault();
    event.stopPropagation();

    const droppedBlock = JSON.parse(event.dataTransfer.getData('activity') || null);
    if (!droppedBlock || droppedBlock.itemId) return;

    const copyActivities = [...props.children,];
    const { key, label, data } = droppedBlock;

    const ancestorElement = event.target.closest('.activity-node');
    if (ancestorElement) {
        let index = ancestorElement.getAttribute('index');
        copyActivities.splice(index, 0, shallowReactive({ key, label, data, itemId: nanoid(10) }))
    } else {
        copyActivities.push(shallowReactive({ key, label, data, itemId: nanoid(10) }))
    }

    emit('update', { children: copyActivities });
}
</script>