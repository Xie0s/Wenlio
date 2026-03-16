import { createApp } from 'vue'
import { createPinia } from 'pinia'
import piniaPluginPersistedstate from 'pinia-plugin-persistedstate'
import './style.css'
import App from './App.vue'
import router from './router'
import { initTheme } from '@/stores/theme'
import { registerServiceWorker } from '@/lib/sw'
import { setupBuildVersionProbe } from '@/lib/version'
import 'vue-sonner/style.css'

const pinia = createPinia()
pinia.use(piniaPluginPersistedstate)

initTheme(pinia)
registerServiceWorker()
setupBuildVersionProbe()

createApp(App)
  .use(pinia)
  .use(router)
  .mount('#app')
