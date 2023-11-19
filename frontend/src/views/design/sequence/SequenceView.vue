<template>
  <div style="height: 100vh" class="overflow-hidden">
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
        <span class="flex">
          <p class="pt-3 pr-2 text-xl font-serif font-semibold">{{ label }}</p>
        </span>
      </template>

      <template #end>
        <Button v-tooltip="'保存'" @click="save">
          <template #icon>
            <span>
              <span
                v-if="dataChanged"
                class="flex h-3 w-3 absolute top-0 left-0 -ml-1 -mt-1"
              >
                <span
                  class="animate-ping absolute inline-flex h-full w-full rounded-full bg-red-600 opacity-75"
                ></span>
                <span
                  class="relative inline-flex rounded-full h-3 w-3 bg-red-600"
                ></span>
              </span>
              <v-remixicon name="riSaveLine" />
            </span>
          </template>
        </Button>
      </template>
    </Toolbar>
    <Splitter id="sequence-designer" class="mb-5" stateStorage="local">
      <SplitterPanel :size="25">
        <Tree
          :value="nodes"
          :filter="true"
          selectionMode="single"
          filterMode="lenient"
          filterPlaceholder="输入组件名称进行搜索"
          class="w-full h-full"
        >
          <template #default="slotProps">
            <div class="flex">
              <p class="pr-2">
                <v-remixicon name="riPlayLine" />
              </p>
              <div
                v-if="isLeaf(slotProps)"
                draggable="true"
                class="transform select-none cursor-move text-ellipsis max-w-xs"
                @dragstart="
                  $event.dataTransfer.setData(
                    'activity',
                    JSON.stringify(slotProps.node)
                  )
                "
              >
                {{ slotProps.node.label }}
              </div>
              <b v-else>{{ slotProps.node.label }}</b>
            </div>
          </template>
        </Tree>
      </SplitterPanel>
      <SplitterPanel class="h-full" :size="75">
        <div class="flex justify-around">
          <Sequence
            :id="activity.id"
            :label="activity.label"
            :showLabel="activity.showLabel"
            :children="activity.children"
            @update="update"
          />
        </div>
      </SplitterPanel>
    </Splitter>
  </div>
</template>
<script setup>
import Button from "primevue/button";
import Splitter from "primevue/splitter";
import SplitterPanel from "primevue/splitterpanel";
import { useRouter, onBeforeRouteLeave } from "vue-router";
import Toolbar from "primevue/toolbar";
import Tree from "primevue/tree";
import { ref, onMounted, reactive } from "vue";
import Sequence from "@/components/activity/Sequence.vue";
import { nanoid } from "nanoid";
const props = defineProps({
  id: {
    type: String,
    default: "",
  },
  subflowId: {
    type: String,
    default: "",
  },
  label: {
    type: String,
    default: "",
  },
});
const nodes = ref(null);
const router = useRouter();
const dataChanged = ref(false);

const activity = reactive({
  id: nanoid(8),
  label: "主流程",
  showLabel: true,
  children: [],
});

function update({ children }) {
  activity.children = children;
}

onMounted(() => {
  nodes.value = [
    {
      key: "0",
      label: "Documents",
      data: {},
      icon_path: "pi pi-fw pi-inbox",
      children: [
        {
          key: "0-0",
          label: "work",
          data: {},
          icon_path: "pi pi-fw pi-cog",
          children: [
            {
              key: "0-0-0",
              label: "Work11111111111111111",
              icon_path: "pi pi-fw pi-file",
              component: "BasicActivity",
              data: {},
            },
            {
              key: "0-0-1",
              label: "Resume.doc",
              icon_path: "pi pi-fw pi-file",
              component: "BasicActivity",
              data: {},
            },
          ],
        },
        {
          key: "0-1",
          label: "Home",
          data: {},
          icon_path: "pi pi-fw pi-home",
          children: [
            {
              key: "0-1-0",
              label: "Invoices.txt",
              icon_path: "pi pi-fw pi-file",
              component: "BasicActivity",
              data: {},
            },
          ],
        },
      ],
    },
  ];
});
function isLeaf(data) {
  return !(data.node.children && data.node.children.length > 0);
}

function save() {}
</script>
<style scoped>
#sequence-designer {
  height: calc(100vh - 68px);
}
:deep(.p-splitter-panel) {
  overflow: auto;
}
</style>
