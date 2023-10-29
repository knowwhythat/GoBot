<template>
    <div class="flex flex-col">
        <Toolbar class="bg-primary p-2">
            <template #start>
                <Button icon="pi pi-plus" class="mr-2" />
                <Button icon="pi pi-print" class="mr-2" />
                <Button icon="pi pi-upload" />
            </template>

            <template #center>
                <span class="p-input-icon-left">
                    <i class="pi pi-search" />
                    <InputText placeholder="Search" />
                </span>
            </template>

            <template #end>
                <SplitButton label="Save" icon="pi pi-check" :model="items"></SplitButton>
            </template>
        </Toolbar>
        <WorkflowEditor class="workflow-editor" />
    </div>
</template>

<script setup>
import Toolbar from 'primevue/toolbar';
import Button from 'primevue/button';
import WorkflowEditor from '@/views/design/flow/WorkflowEditor.vue'
</script>
<style scoped>
.workflow-editor {
    height: calc(100vh - 40px);
}
</style>