import { createApp } from 'vue'
import App from '@/App.vue'
import { router } from '@/router'
import { createPinia } from 'pinia'
import PrimeVue from 'primevue/config'
import ToastService from 'primevue/toastservice'
import Tooltip from 'primevue/tooltip'
import ConfirmationService from 'primevue/confirmationservice';
import 'primeicons/primeicons.css'
import vRemixicon, { icons } from '@/components/vRemixicon';
import './style.css'
import 'primevue/resources/themes/lara-light-teal/theme.css'


const pinia = createPinia()
const app = createApp(App)
app.use(router)
app.use(pinia)
app.use(PrimeVue, {
    ptOptions: {
        mergeSections: true,
        mergeProps: false,
    },
    ripple: true
});
app.use(ToastService)
app.use(ConfirmationService)
app.directive('tooltip', Tooltip);
app.use(vRemixicon, icons)

app.mount('#app')
