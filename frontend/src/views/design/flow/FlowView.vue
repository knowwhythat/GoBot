<template>
    <VueFlow fit-view-on-init class="flow" v-model="elements">
        <Background />
        <Panel :position="PanelPosition.BottomRight">
            <div>
                <label for="ishidden">
                    hidden
                    <input id="ishidden" v-model="isHidden" type="checkbox" />
                </label>
            </div>
        </Panel>
        <Controls />
    </VueFlow>
</template>
<script setup>
import '@vue-flow/core/dist/style.css'
import '@vue-flow/core/dist/theme-default.css'
import { ref, onMounted } from "vue"
import { Background, Panel, PanelPosition, Controls } from '@vue-flow/additional-components'
import { VueFlow, useVueFlow } from '@vue-flow/core'

const elements = ref([
    // Nodes
    // An input node, specified by using `type: 'input'`
    { id: '1', type: 'input', label: 'Node 1', position: { x: 250, y: 5 } },

    // Default nodes, you can omit `type: 'default'`
    { id: '2', label: 'Node 2', position: { x: 100, y: 100 }, },
    { id: '3', label: 'Node 3', position: { x: 400, y: 100 } },

    // An output node, specified by using `type: 'output'`
    { id: '4', type: 'output', label: 'Node 4', position: { x: 400, y: 200 } },

    // Edges
    // Most basic edge, only consists of an id, source-id and target-id
    { id: 'e1-3', source: '1', target: '3' },

    // An animated edge
    { id: 'e1-2', source: '1', target: '2', animated: true },
])
let { onPaneReady, onNodeDragStop, onConnect, addEdges, setTransform, toObject, nodes, edges } = useVueFlow()

const isHidden = ref(false)

</script>
<style scoped>
.flow {
    height: 100vh;
    width: 100vw;
}
</style>