<template>
    <Splitter style="height: 100vh" class="mb-5" stateStorage="local">
        <SplitterPanel class="flex align-items-center justify-content-center" :size="25">
            <div class="card flex justify-content-center">
                <Tree :value="nodes" :filter="true" filterMode="lenient" class="w-full md:w-30rem">
                    <template #default="slotProps">
                        <span v-if="isLeaf(slotProps)" draggable="true" class="transform select-none cursor-move"
                            @dragstart="$event.dataTransfer.setData('activity', JSON.stringify(slotProps.node))">{{
                                slotProps.node.label }}</span>
                        <b v-else>{{ slotProps.node.label }}</b>
                    </template>
                </Tree>
            </div>
        </SplitterPanel>
        <SplitterPanel :size="75">
            <div class="flex justify-around">
                <Sequence :id="activity.id" :label="activity.label" :showLabel="activity.showLabel"
                    :children="activity.children" @update="update" />
            </div>

        </SplitterPanel>
    </Splitter>
</template>
<script setup>
import Splitter from 'primevue/splitter';
import SplitterPanel from 'primevue/splitterpanel';
import Tree from 'primevue/tree';
import { ref, onMounted, reactive } from 'vue';
import Sequence from '@/components/activity/Sequence.vue'
import { nanoid } from 'nanoid';

const nodes = ref(null);

const activity = reactive({
    id: nanoid(8),
    label: '主流程',
    showLabel: true,
    children: []
})

function update({ children }) {
    activity.children = children
}

onMounted(() => {
    nodes.value = [{
        key: '0',
        label: 'Documents',
        data: 'Documents Folder',
        icon: 'pi pi-fw pi-inbox',
        children: [
            {
                key: '0-0',
                label: 'work',
                data: 'Work Folder11111111111111111111111111111111111111111111111111111111111111111111111111111111111',
                icon: 'pi pi-fw pi-cog',
                children: [
                    { key: '0-0-0', label: 'Work1111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111', icon: 'pi pi-fw pi-file', data: 'Expenses Document' },
                    { key: '0-0-1', label: 'Resume.doc', icon: 'pi pi-fw pi-file', data: 'Resume Document' }
                ]
            },
            {
                key: '0-1',
                label: 'Home',
                data: 'Home Folder',
                icon: 'pi pi-fw pi-home',
                children: [{ key: '0-1-0', label: 'Invoices.txt', icon: 'pi pi-fw pi-file', data: 'Invoices for this month' }]
            }
        ]
    }]
});
function isLeaf(data) {
    return !(data.node.children && data.node.children.length > 0)
}



</script>