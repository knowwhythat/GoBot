<template>
    <div style="height: 100%;">
        <DataView :value="projects" :layout="layout">
            <template #header>
                <Toolbar class="bg-gray-100 p-2">
                    <template #start>
                        <div class="flex gap-1">
                            <span class="p-input-icon-left">
                                <i class="pi pi-search" />
                                <InputText v-model="searchText" style="margin-left: 40px" placeholder="Search" />
                            </span>
                            <Button label="导入">
                                <template #icon>
                                    <v-remixicon name="riDownloadLine" />
                                </template>
                            </Button>
                            <Button label="新建" @click="newWorkflow">
                                <template #icon>
                                    <v-remixicon name="riAddLine" />
                                </template>
                            </Button>
                        </div>
                    </template>
                    <template #end>
                        <SelectButton v-model="layout" :options="options" optionValue="value" optionLabel="value"
                            dataKey="value" aria-labelledby="custom">
                            <template #option="slotProps">
                                <i :class="slotProps.option.icon" size="32"></i>
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
                            <SplitButton label="设置" icon="pi pi-play" @click="run(slotProps.data.id)"
                                :model="getItems(slotProps.data.id)">
                                <template #icon>
                                    <v-remixicon size="24" name="riSettings3Line" />
                                </template>
                            </SplitButton>
                            <Button icon="pi pi-file-edit" @click="edit(slotProps.data.id)" v-tooltip="'编辑'" />
                            <Button icon="pi pi-caret-right" v-tooltip="'运行'">
                                <template #icon>
                                    <v-remixicon size="24" name="riPlayLine" />
                                </template>
                            </Button>
                        </div>
                    </div>
                </div>
            </template>
        </DataView>
        <Dialog v-model:visible="newFlow.visible" modal header="Header">
            <template #header>
                <div class="inline-flex align-items-center justify-content-center gap-2">
                    <v-remixicon size="20" name="riAddLine" />
                    <span class="font-bold white-space-nowrap">新建流程</span>
                </div>
            </template>
            <div class="card flex justify-content-center">
                <form @submit="onSubmit" class="flex flex-col gap-2">
                    <span>
                        <label for="value">流程名称</label>
                        <InputText id="value" v-model="newFlow.value" type="text"
                            :class="{ 'p-invalid': newFlow.errorMessage }" aria-describedby="text-error" />
                    </span>
                    <small class="p-error" id="text-error">{{ newFlow.errorMessage || '&nbsp;' }}</small>
                </form>
            </div>

            <template #footer>
                <Button label="取消" icon="pi pi-check" @click="newFlow.visible = false" />
                <Button label="确定" icon="pi pi-check" @click="newFlow.visible = false" />
            </template>
        </Dialog>
    </div>
</template>

<script setup>
import { ref, onMounted, reactive, watch } from "vue"
import { useRouter } from 'vue-router'
import DataView from 'primevue/dataview'
import Toolbar from 'primevue/toolbar';
import InputText from "primevue/inputtext";
import SelectButton from "primevue/selectbutton";
import SplitButton from 'primevue/splitbutton'
import Button from 'primevue/button'
import Dialog from 'primevue/dialog'
import { useToast } from "primevue/usetoast"
import { throttle } from 'lodash-es'
import { NewWorkflow, ListProject, DeleteProject } from "@back/go/main/App"
const toast = useToast()
const router = useRouter()

const projects = ref();
const layout = ref('grid');
const options = ref([
    { icon: 'pi pi-table', value: 'grid' },
    { icon: 'pi pi-list', value: 'list' },
]);

onMounted(() => {
    listProject()
});

const searchText = ref("")
watch(searchText, throttle((newVal, oldVal) => {
    listProject(newVal)
}, 1000))

function listProject(name) {
    ListProject(name).then((result) => {
        projects.value = result.list
    }).catch((error) => {
        toast.add({ severity: 'error', summary: '加载失败', detail: error, life: 3000 });
        console.error(error)
    })
}

function getItems(id) {
    return [
        {
            label: '编辑',
            icon: 'pi pi-file-edit',
            command: () => {

            }
        },
        {
            label: '删除',
            icon: 'pi pi-times',
            command: () => {
                DeleteProject(id).then((result) => {
                    listProject()
                }).catch((error) => {
                    console.error(error)
                })
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
}


const run = () => {
    toast.add({ severity: 'success', summary: 'Success', detail: 'Data Saved', life: 3000 });
};

const edit = () => {
    router.push("/design")
}

const newFlow = reactive({
    visible: false,
    flowName: '',
    errorMessage: ""

})

function newWorkflow() {
    newFlow.visible = true
    // NewWorkflow("test").then((result) => {
    //     toast.add({ severity: 'warn', summary: 'Delete', detail: result, life: 3000 });
    // }).catch((error) => {
    //     console.error(error)
    // });
}
</script>
<style scoped>
:deep(.p-dataview-content) {
    height: calc(100vh - 298px);
    overflow: auto;
}

:deep(.p-grid) {
    @apply sm:grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4;
}
</style>