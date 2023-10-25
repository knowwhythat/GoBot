
<template>
    <div class="card" style="height: 100%;">
        <DataView :value="products" :layout="layout">
            <template #header>
                <div class="flex justify-content-end">
                    <DataViewLayoutOptions v-model="layout" />
                </div>
            </template>

            <template #list="slotProps">
                <div class="col-12">
                    <div class="flex flex-column xl:flex-row xl:align-items-start p-4 gap-4">
                        <div
                            class="flex flex-column sm:flex-row justify-content-between align-items-center xl:align-items-start flex-1 gap-4">
                            <div class="flex flex-column align-items-center sm:align-items-start gap-3">
                                <div class="text-2xl font-bold text-900">{{ slotProps.data.name }}</div>
                                <div class="flex align-items-center gap-3">
                                    <span class="flex align-items-center gap-2">
                                        <i class="pi pi-tag"></i>
                                        <span class="font-semibold">{{ slotProps.data.category }}</span>
                                    </span>
                                </div>
                            </div>
                            <div class="flex sm:flex-column align-items-center sm:align-items-end gap-3 sm:gap-2">
                            </div>
                        </div>
                    </div>
                </div>
            </template>

            <template #grid="slotProps">
                <div class="col-12 sm:col-6 lg:col-12 xl:col-4 p-2">
                    <div class="p-4 border-1 surface-border surface-card border-round">
                        <div class="flex flex-wrap align-items-center justify-content-between gap-2">
                            <div class="flex align-items-center gap-2">
                                <i class="pi pi-tag"></i>
                                <span class="font-semibold">{{ slotProps.data.category }}</span>
                            </div>
                        </div>
                        <div class="flex flex-column align-items-center gap-3 py-5">
                            <div class="text-2xl font-bold">{{ slotProps.data.name }}</div>
                        </div>
                        <div class="flex align-items-right">
                            <SplitButton label="运行" icon="pi pi-play" @click="run" :model="items">
                                <template #icon>
                                    <v-remixicon size="24" name="riPlayLine" />
                                </template>
                            </SplitButton>
                        </div>
                    </div>
                </div>
            </template>
        </DataView>
    </div>
</template>

<script setup>
import { ref, onMounted } from "vue"
import { useRouter } from 'vue-router'
import DataView from 'primevue/dataview'
import DataViewLayoutOptions from 'primevue/dataviewlayoutoptions'
import SplitButton from 'primevue/splitbutton'
import { useToast } from "primevue/usetoast"
const toast = useToast()
const router = useRouter()

onMounted(() => {
    products.value = [
        { 'category': 'test', 'name': 'test' },
        { 'category': 'test', 'name': 'test' },
        { 'category': 'test', 'name': 'test' },
        { 'category': 'test', 'name': 'test' },
        { 'category': 'test', 'name': 'test' },
    ]
});

const products = ref();
const layout = ref('grid');

const items = [
    {
        label: '编辑',
        icon: 'pi pi-file-edit',
        command: () => {
            router.push("/design/WorkflowView")
        }
    },
    {
        label: '删除',
        icon: 'pi pi-times',
        command: () => {
            toast.add({ severity: 'warn', summary: 'Delete', detail: 'Data Deleted', life: 3000 });
        }
    },
    {
        label: '发布',
        icon: 'pi pi-send',
        command: () => {
            window.location.href = 'https://vuejs.org/';
        }
    },
    { label: 'Upload', icon: 'pi pi-upload', to: '/fileupload' }
];

const run = () => {
    toast.add({ severity: 'success', summary: 'Success', detail: 'Data Saved', life: 3000 });
};
</script>
