<template>
  <div class="login-form-wrapper">
    <h2 class="form-title">登录</h2>
    <p class="form-subtitle">欢迎回来，请登录您的账号</p>

    <el-form
      ref="formRef"
      :model="form"
      :rules="rules"
      @submit.prevent="handleLogin"
    >
      <el-form-item prop="username">
        <el-input v-model="form.username" placeholder="用户名" :prefix-icon="User" size="large" />
      </el-form-item>
      <el-form-item prop="password">
        <el-input v-model="form.password" type="password" placeholder="密码" :prefix-icon="Lock" size="large" show-password />
      </el-form-item>
      <el-form-item>
        <el-button type="primary" size="large" :loading="loading" native-type="submit" class="submit-btn">
          {{ loading ? '登录中...' : '登录' }}
        </el-button>
      </el-form-item>
    </el-form>

    <div class="toggle-view">
      还没有账号？<el-button type="primary" link @click="$emit('switch')">立即注册</el-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import type { FormInstance, FormRules } from 'element-plus'
import { User, Lock } from '@element-plus/icons-vue'

defineProps<{
  loading: boolean
}>()

const emit = defineEmits<{
  (e: 'submit', username: string, password: string): void
  (e: 'switch'): void
}>()

const formRef = ref<FormInstance>()
const form = reactive({ username: '', password: '' })

const rules: FormRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 20, message: '长度在 3 到 20 个字符', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, max: 20, message: '长度在 6 到 20 个字符', trigger: 'blur' }
  ]
}

async function handleLogin() {
  await formRef.value?.validate()
  emit('submit', form.username, form.password)
}

defineExpose({ form })
</script>

<style scoped>
.login-form-wrapper {
  max-width: 380px;
  width: 100%;
  margin: 0 auto;
}

.form-title {
  font-size: 26px;
  font-weight: 700;
  color: #1d2129;
  margin: 0 0 6px;
}

.form-subtitle {
  font-size: 14px;
  color: #86909c;
  margin: 0 0 28px;
}

.submit-btn {
  width: 100%;
  height: 44px;
  font-size: 15px;
  font-weight: 600;
  border-radius: 8px;
  margin-top: 4px;
}

.toggle-view {
  text-align: center;
  margin-top: 24px;
  font-size: 14px;
  color: #86909c;
}
</style>
