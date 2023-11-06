import { createRouter, createWebHistory } from 'vue-router'
import Main from '@/views/main/MainView.vue'
import Flow from '@/views/main/FlowView.vue'
import Trigger from '@/views/main/TriggerView.vue'
import Record from '@/views/main/RecordView.vue'
import Log from '@/views/main/LogView.vue'
import DesignView from '@/views/design/DesignView.vue'
import WorkflowView from '@/views/design/flow/WorkflowView.vue'
import SequenceView from '@/views/design/sequence/SequenceView.vue'

export const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      component: Main,
      children: [
        { path: '', component: Flow },
        { path: 'trigger', component: Trigger },
        { path: 'record', component: Record },
        { path: 'log', component: Log },
      ],
    },
    {
      path: '/design',
      component: DesignView,
      children :[
        { path: '', component: WorkflowView, props: route => ({ id: route.query.id })},
        { path: 'sequence', component: SequenceView },
      ]
    }
  ],
})
