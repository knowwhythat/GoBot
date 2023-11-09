<template>
    <div class="flex flex-col">
        <Toolbar class="p-2">
            <template #start>
                <Button @click="router.back()" v-tooltip="'返回'" class="mr-2">
                    <template #icon>
                        <v-remixicon name="riArrowGoBackLine" />
                    </template>
                </Button>
                <Button v-tooltip="'运行'" class="mr-2">
                    <template #icon>
                        <v-remixicon name="riPlayLine" />
                    </template>
                </Button>
            </template>

            <template #center>
                <span class="flex ">
                    <p class="pt-3 pr-2">工作流名称</p>
                    <Button v-tooltip="'编辑'">
                        <template #icon>
                            <v-remixicon name="riEditBoxLine" />
                        </template>
                    </Button>
                </span>
            </template>

            <template #end>
                <Button v-tooltip="'发布'" class="mr-2">
                    <template #icon>
                        <v-remixicon name="riSendPlaneLine" />
                    </template>
                </Button>
                <Button v-tooltip="'保存'" @click="save">
                    <template #icon>
                        <span>
                            <span v-if="state.dataChanged" class="flex h-3 w-3 absolute top-0 left-0 -ml-1 -mt-1">
                                <span
                                    class="animate-ping absolute inline-flex h-full w-full rounded-full bg-red-600 opacity-75"></span>
                                <span class="relative inline-flex rounded-full h-3 w-3 bg-red-600"></span>
                            </span>
                            <v-remixicon name="riSaveLine" />
                        </span>
                    </template>
                </Button>
            </template>
        </Toolbar>
        <WorkflowEditor v-if="loaded" :id="props.id" :data="workflow" class="workflow-editor" @init="onEditorInit"
            @edit="editBlock" @update:flow="state.dataChanged = true" @update:node="state.dataChanged = true"
            @delete:node="state.dataChanged = true" />
    </div>
</template>

<script setup>
import { reactive, ref, shallowRef, onBeforeMount } from 'vue';
import { useRouter, onBeforeRouteLeave } from 'vue-router'
import Toolbar from 'primevue/toolbar';
import Button from 'primevue/button';
import WorkflowEditor from '@/views/design/flow/WorkflowEditor.vue'
import { useToast } from "primevue/usetoast"
import { GetMainFlow, SaveMainFlow } from "@back/go/main/App"
import { trim } from 'lodash-es'
const toast = useToast()

const props = defineProps({
    id: {
        type: String,
        default: '',
    },
});

const workflow = ref(null)
const loaded = ref(false)

onBeforeMount(() => {
    console.log(props.id)
    GetMainFlow(props.id).then((result) => {
        if (trim(result).length > 0) {
            workflow.value = JSON.parse(result)
        }
        loaded.value = true
    }).catch((error) => {
        loaded.value = true
        console.error(error)
        toast.add({ severity: 'error', summary: '错误', detail: error, life: 3000 });
    })
})
const state = reactive({
    dataChanged: false,
})

onBeforeRouteLeave(onBeforeLeave);

function onBeforeLeave() {
    if (state.dataChanged) {
        const confirm = window.confirm("有修改尚未保存，是否退出");
        if (!confirm) return false;
    }
}

const router = useRouter()
const editor = shallowRef(null);

function onEditorInit(instance) {
    editor.value = instance;
}

function save() {
    const flow = editor.value.toObject();
    flow.edges = flow.edges.map((edge) => {
        delete edge.sourceNode;
        delete edge.targetNode;

        return edge;
    });

    const triggerBlock = flow.nodes.find((node) => node.label === 'Start');
    if (!triggerBlock) {
        toast.add({ severity: 'warn', summary: '警告', detail: '当前流程图中不存在开始节点', life: 3000 });
        return;
    }
    console.log(JSON.stringify(flow))
    SaveMainFlow(props.id, JSON.stringify(flow)).then((result) => {
        state.dataChanged = false
    }).catch((error) => {
        toast.add({ severity: 'error', summary: '未知异常', detail: error, life: 3000 });
    })
}
function editBlock(data) {
    console.log(data)
    if (data.id = '子流程') {
        router.push('design/sequence')
    }
}
</script>
<style scoped>
.workflow-editor {
    height: calc(100vh - 70px);
}
</style>