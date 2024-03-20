import {createRouter, createWebHistory} from "vue-router";
import Login from "@/views/Login.vue";
import Main from "@/views/main/MainView.vue";
import Flow from "@/views/main/FlowView.vue";
import Trigger from "@/views/main/TriggerView.vue";
import Record from "@/views/main/RecordView.vue";
import Setting from "@/views/main/SettingView.vue";
import Monitor from "@/views/main/MonitorView.vue";
import DesignView from "@/views/design/DesignView.vue";
import WorkflowView from "@/views/design/flow/WorkflowView.vue";
import SequenceView from "@/views/design/sequence/MainView.vue";

export const router = createRouter({
    history: createWebHistory(),
    routes: [
        {
            path: '/',
            component: Login
        },
        {
            path: "/main",
            component: Main,
            children: [
                {path: "", component: Flow},
                {path: "trigger", component: Trigger},
                {path: "record", component: Record},
                {path: "monitor", component: Monitor},
                {path: "setting", component: Setting},
            ],
        },
        {
            path: "/design",
            component: DesignView,
            children: [
                {
                    path: "",
                    component: WorkflowView,
                    props: (route) => ({id: route.query.id}),
                },
                {
                    path: "sequence",
                    component: SequenceView,
                    props: (route) => ({
                        id: route.query.id,
                        subflowId: route.query.subflowId,
                        label: route.query.label,
                    }),
                },
            ],
        },
    ],
});
