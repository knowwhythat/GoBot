import { createApp } from 'vue'
import App from '@/App.vue'
import { router } from '@/router'
import PrimeVue from 'primevue/config'
import 'primevue/resources/themes/lara-light-teal/theme.css'
// import 'primevue/resources/themes/lara-dark-teal/theme.css'
import 'primeflex/primeflex.css'
import ToastService from 'primevue/toastservice'
import Tooltip from 'primevue/tooltip'
import 'primeicons/primeicons.css'
import './style.css'

const app = createApp(App)
app.use(router)
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
