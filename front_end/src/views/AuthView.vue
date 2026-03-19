<template>
  <div class="auth-container">
    <div class="auth-card">
      <!-- 左侧装饰区 -->
      <div class="image-side">
        <div class="image-overlay">
          <!-- 装饰图形 -->
          <div class="deco-circles">
            <div class="circle c1"></div>
            <div class="circle c2"></div>
            <div class="circle c3"></div>
          </div>

          <div class="brand-area">
            <div class="brand-logo">
              <el-icon :size="40" color="#fff"><Ship /></el-icon>
            </div>
            <h2 class="brand-title">智钓蓝海</h2>
            <p class="brand-desc">边缘-云协同的全生态智能垂钓系统</p>
          </div>

          <div class="feature-list">
            <div class="feature-item">
              <el-icon :size="18"><MapLocation /></el-icon>
              <span>智能水域监测</span>
            </div>
            <div class="feature-item">
              <el-icon :size="18"><Bell /></el-icon>
              <span>实时垂钓提醒</span>
            </div>
            <div class="feature-item">
              <el-icon :size="18"><DataAnalysis /></el-icon>
              <span>环境数据分析</span>
            </div>
          </div>

          <!-- 底部波浪 -->
          <svg class="wave-bottom" viewBox="0 0 400 60" preserveAspectRatio="none">
            <path d="M0 30 Q100 0 200 30 T400 30 V60 H0Z" fill="rgba(255,255,255,0.08)" />
            <path d="M0 40 Q100 15 200 40 T400 40 V60 H0Z" fill="rgba(255,255,255,0.05)" />
          </svg>
        </div>
      </div>

      <!-- 右侧表单区 -->
      <div class="form-side">
        <div class="form-container">
          <transition name="form-slide" mode="out-in">
            <LoginForm
              v-if="activeTab === 'login'"
              key="login"
              ref="loginFormRef"
              :loading="authStore.isLoading"
              @submit="handleLogin"
              @switch="activeTab = 'register'"
            />
            <RegisterForm
              v-else
              key="register"
              :loading="authStore.isLoading"
              @submit="handleRegister"
              @switch="activeTab = 'login'"
            />
          </transition>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { ElMessage } from 'element-plus'
import { Ship, MapLocation, Bell, DataAnalysis } from '@element-plus/icons-vue'

// 子组件
import LoginForm from './auth/LoginForm.vue'
import RegisterForm from './auth/RegisterForm.vue'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

const activeTab = ref('login')
const loginFormRef = ref<InstanceType<typeof LoginForm>>()

async function handleLogin(username: string, password: string) {
  try {
    await authStore.login(username, password)
    ElMessage.success('登录成功')
    const redirect = (route.query.redirect as string) || '/'
    router.push(redirect)
  } catch (e: any) {
    ElMessage.error(e.response?.data?.error || '登录失败')
  }
}

async function handleRegister(data: { username: string; password: string; phone: string; email: string }) {
  try {
    await authStore.register(data.username, data.password, data.phone, data.email)
    ElMessage.success('注册成功，请登录')
    activeTab.value = 'login'
    if (loginFormRef.value) {
      loginFormRef.value.form.username = data.username
    }
  } catch (e: any) {
    ElMessage.error(e.response?.data?.error || '注册失败')
  }
}
</script>

<style scoped>
.auth-container {
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 40px 20px 60px;
  min-height: calc(100vh - 160px);
}

.auth-card {
  display: flex;
  width: 900px;
  max-width: 100%;
  border-radius: 16px;
  overflow: hidden;
  box-shadow:
    0 4px 6px -1px rgba(0, 0, 0, 0.05),
    0 10px 30px -5px rgba(0, 0, 0, 0.08);
  background: #fff;
}

/* 左侧装饰区 */
.image-side {
  position: relative;
  width: 42%;
  background: linear-gradient(160deg, #1e3a5f 0%, #2c6fbb 50%, #3d8bda 100%);
  overflow: hidden;
  flex-shrink: 0;
}

.image-overlay {
  position: relative;
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  padding: 40px 30px;
  z-index: 1;
}

/* 装饰圆 */
.deco-circles {
  position: absolute;
  inset: 0;
  pointer-events: none;
}

.circle {
  position: absolute;
  border-radius: 50%;
  border: 1px solid rgba(255, 255, 255, 0.1);
}

.circle.c1 {
  width: 300px;
  height: 300px;
  top: -80px;
  right: -100px;
  background: radial-gradient(circle, rgba(255, 255, 255, 0.06) 0%, transparent 70%);
}

.circle.c2 {
  width: 200px;
  height: 200px;
  bottom: -60px;
  left: -60px;
  background: radial-gradient(circle, rgba(255, 255, 255, 0.04) 0%, transparent 70%);
}

.circle.c3 {
  width: 120px;
  height: 120px;
  top: 40%;
  left: 10%;
  border: 1px solid rgba(255, 255, 255, 0.08);
}

/* 品牌 */
.brand-area {
  text-align: center;
  z-index: 2;
}

.brand-logo {
  width: 72px;
  height: 72px;
  margin: 0 auto 20px;
  background: rgba(255, 255, 255, 0.15);
  border-radius: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  backdrop-filter: blur(8px);
  border: 1px solid rgba(255, 255, 255, 0.2);
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
}

.brand-title {
  font-size: 28px;
  font-weight: 700;
  color: #fff;
  margin: 0 0 8px;
  letter-spacing: 2px;
}

.brand-desc {
  font-size: 13px;
  color: rgba(255, 255, 255, 0.75);
  margin: 0;
  line-height: 1.6;
}

/* 特性列表 */
.feature-list {
  margin-top: 40px;
  z-index: 2;
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.feature-item {
  display: flex;
  align-items: center;
  gap: 10px;
  color: rgba(255, 255, 255, 0.85);
  font-size: 14px;
  padding: 8px 16px;
  border-radius: 8px;
  background: rgba(255, 255, 255, 0.08);
  backdrop-filter: blur(4px);
  transition: background 0.2s;
}

.feature-item:hover {
  background: rgba(255, 255, 255, 0.14);
}

/* 底部波浪 */
.wave-bottom {
  position: absolute;
  bottom: 0;
  left: 0;
  width: 100%;
  height: 60px;
}

/* 右侧表单区 */
.form-side {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 36px 44px;
  background: #fff;
  min-height: 480px;
}

.form-container {
  width: 100%;
  max-width: 380px;
}

/* 切换动画 */
.form-slide-enter-active,
.form-slide-leave-active {
  transition: all 0.3s ease;
}

.form-slide-enter-from {
  opacity: 0;
  transform: translateX(30px);
}

.form-slide-leave-to {
  opacity: 0;
  transform: translateX(-30px);
}

/* 移动端响应式 */
@media (max-width: 768px) {
  .auth-container {
    padding: 20px 16px 40px;
    min-height: auto;
  }

  .auth-card {
    flex-direction: column;
    min-height: auto;
  }

  .auth-card.register-mode {
    min-height: auto;
  }

  .image-side {
    width: 100%;
    min-height: 200px;
    max-height: 220px;
  }

  .image-overlay {
    padding: 24px 20px;
  }

  .brand-logo {
    width: 56px;
    height: 56px;
    border-radius: 16px;
    margin-bottom: 12px;
  }

  .brand-logo .el-icon {
    font-size: 28px !important;
  }

  .brand-title {
    font-size: 22px;
  }

  .brand-desc {
    font-size: 12px;
  }

  .feature-list {
    display: none;
  }

  .form-side {
    padding: 28px 24px;
  }
}

@media (max-width: 480px) {
  .image-side {
    max-height: 180px;
  }

  .form-side {
    padding: 24px 16px;
  }
}
</style>
