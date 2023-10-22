import { createApp } from 'vue'
import App from '@/App.vue'
import PrimeVue from 'primevue/config'
import 'primevue/resources/themes/lara-light-teal/theme.css'
import ToastService from 'primevue/toastservice'
import Tooltip from 'primevue/tooltip'
import './style.css'

const app = createApp(App)
app.use(PrimeVue, {
    ptOptions: {
        mergeSections: true,
        mergeProps: false,
    },
    ripple: true
});
app.use(ToastService)
app.directive('tooltip', Tooltip);

app.mount('#app')
