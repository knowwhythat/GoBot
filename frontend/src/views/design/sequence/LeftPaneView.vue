<template>
  <TabView>
    <TabPanel header="组件库">
      <Tree
        :value="nodes"
        :filter="true"
        selectionMode="single"
        filterMode="lenient"
        filterPlaceholder="输入组件名称进行搜索"
        class="w-full h-full"
        :pt="{
          label: (options) => ({
            class: ['w-full'],
          }),
        }"
      >
        <template #default="slotProps">
          <div
            v-if="isLeaf(slotProps)"
            draggable="true"
            @dragstart="toolDragStart($event, slotProps.node)"
            @dragend="toolDragEnd($event)"
            class="w-full flex cursor-move"
          >
            <p class="pr-2">
              <v-remixicon v-bind="getIconPath(slotProps.node.icon_path)" />
            </p>
            <div class="transform select-none text-ellipsis max-w-xs">
              {{ slotProps.node.label }}
            </div>
          </div>
          <div v-else class="flex">
            <p class="pr-2">
              <v-remixicon v-bind="getIconPath(slotProps.node.icon_path)" />
            </p>
            <b class="cursor-default">{{ slotProps.node.label }}</b>
          </div>
        </template>
      </Tree>
    </TabPanel>
    <TabPanel header="项目树"></TabPanel>
  </TabView>
</template>
<script setup>
import TabView from "primevue/tabview";
import TabPanel from "primevue/tabpanel";
import Tree from "primevue/tree";
import { getIconPath } from "@/utils/helper";
import { ref, onMounted, inject } from "vue";
import { ParseAllPlugin } from "@back/go/main/App";

const nodes = ref(null);

onMounted(() => {
  ParseAllPlugin()
    .then((result) => {
      nodes.value = result;
    })
    .catch((err) => {
      console.error(err);
      toast.add({
        severity: "error",
        summary: "加载组件树异常",
        detail: err,
        life: 3000,
      });
    });
});

function isLeaf(data) {
  return !(data.node.children && data.node.children.length > 0);
}

const { dragBlockId } = inject("dragBlockId");
const { dropBlockId } = inject("dropBlockId");

function toolDragStart(event, element) {
  dragBlockId.value = "";
  event.dataTransfer.setData("activity", JSON.stringify(element));
}

function toolDragEnd(event) {
  dropBlockId.value = "";
}
</script>
<style scoped>
:deep(.p-tree-wrapper) {
  height: calc(100vh - 168px);
}
:deep(.p-tree) {
  padding: 0.25rem;
}
:deep(.p-treenode-icon) {
  display: none;
}
:deep(.p-treenode-label) {
  margin-left: -0.75rem;
}
</style>
