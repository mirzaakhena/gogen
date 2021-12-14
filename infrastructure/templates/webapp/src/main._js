import {createApp} from 'vue'
import App from './App._vue'
import router from './router/index._js'

import 'bootstrap/scss/bootstrap.scss'
import 'sweetalert2/src/sweetalert2.scss'

createApp(App).//
    use(router).//
    mount('#app')