// main.js or main.ts

import {createApp} from 'vue';
import App from './App.vue';
import axios from './services/axios.js';
import router from './router';


// createApp(App).use(router).mount('#app')
const app = createApp(App)
app.config.globalProperties.$axios = axios;
app.use(router)
app.mount('#app')