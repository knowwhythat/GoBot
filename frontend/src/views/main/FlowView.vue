<template>
    <div style="height: 100%;">
        <DataView :value="products" :layout="layout">
            <template #header>
                <Toolbar class="bg-gray-100 p-2">
                    <template #start>
                        <span class="p-input-icon-left">
                            <i class="pi pi-search" />
                            <InputText v-model="value1" style="margin-left: 40px" placeholder="Search" />
                        </span>
                    </template>
                    <template #end>
                        <Button label="导入" class="bg-primary text-white p-2 mr-8">
                            <template #icon>
                                <v-remixicon name="riDownloadLine" />
                            </template>
                        </Button>
                        <Button label="新建" class="bg-primary text-white p-2 mr-8">
                            <template #icon>
                                <v-remixicon name="riAddLine" />
                            </template>
                        </Button>
                        <SelectButton v-model="layout" :options="options" optionValue="value" optionLabel="value"
                            dataKey="value" aria-labelledby="custom">
                            <template #option="slotProps">
                                <i :class="slotProps.option.icon"></i>
                            </template>
                        </SelectButton>
                    </template>
                </Toolbar>
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
                <div class="border-gray-200 border-current border rounded-lg">
                    <div class="p-4 border-1 surface-border border-round">
                        <div class="flex flex-wrap align-items-center justify-content-between gap-2">
                            <div class="flex align-items-center gap-2">
                                <v-remixicon size="18" name="riPriceTag3Line" />
                                <span class="font-semibold">{{ slotProps.data.category }}</span>
                            </div>
                        </div>
                        <div class="flex flex-column align-items-center gap-3 py-5">
                            <div class="text-xl font-bold truncate ">{{ slotProps.data.name }}</div>
                        </div>
                        <div class="flex flex-row-reverse gap-1">
                            <SplitButton class="bg-primary text-white p-2" label="设置" icon="pi pi-play" @click="run"
                                :model="items">
                                <template #icon>
                                    <v-remixicon size="24" name="riSettings3Line" />
                                </template>
                            </SplitButton>
                            <Button class="bg-primary text-white p-2" icon="pi pi-file-edit" @click="edit"
                                v-tooltip="'编辑'" />
                            <Button class="bg-primary text-white p-2" icon="pi pi-caret-right" v-tooltip="'运行'">
                                <template #icon>
                                    <v-remixicon size=" 24" name="riPlayLine" />
                                </template>
                            </Button>
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
import Toolbar from 'primevue/toolbar';
import InputText from "primevue/inputtext";
import SelectButton from "primevue/selectbutton";
import SplitButton from 'primevue/splitbutton'
import Button from 'primevue/button'
import { useToast } from "primevue/usetoast"
const toast = useToast()
const router = useRouter()

onMounted(() => {
    products.value = [
        { 'category': 'test', 'name': '测试中文名称很长的按实际发生备份脚本的时间副书记快捷方式国际快递发几个快递发几个防水等级覅是降低房价多少' },
        { 'category': 'test', 'name': 'test' },
        { 'category': 'test', 'name': 'test' },
        { 'category': 'test', 'name': 'test' },
        { 'category': 'test', 'name': 'test' },
    ]
});

const products = ref();
const layout = ref('grid');
const options = ref([
    { icon: 'pi pi-table', value: 'grid' },
    { icon: 'pi pi-list', value: 'list' },
]);

const items = [
    {
        label: '编辑',
        icon: 'pi pi-file-edit',
        command: () => {
            router.push("/design")
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
            toast.add({ severity: 'warn', summary: 'Delete', detail: 'Data Deleted', life: 3000 });
        }
    },
    { label: 'Upload', icon: 'pi pi-upload', to: '/fileupload' }
];

const run = () => {
    toast.add({ severity: 'success', summary: 'Success', detail: 'Data Saved', life: 3000 });
};

const edit = () => {
    router.push("/design")
}
</script>
<style scoped>
.card:deep(.p-dataview-content) {
    height: calc(100vh - 268px);
    overflow: auto;
}

:deep(.p-grid) {
    @apply sm:grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4;
}
</style>