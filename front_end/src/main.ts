import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import router from './router'
import { setupElementPlus } from './plugins/element-plus'
import { useAuthStore } from './stores/auth'

const app = createApp(App)
const pinia = createPinia()

app.use(pinia)
app.use(router)
setupElementPlus(app)

// 初始化认证状态
const authStore = useAuthStore()
authStore.initialize().then(() => {
  app.mount('#app')
})
