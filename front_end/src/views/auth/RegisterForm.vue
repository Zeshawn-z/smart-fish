<template>
  <div class="register-form-wrapper">
    <h2 class="form-title">注册</h2>
    <p class="form-subtitle">创建您的智钓蓝海账号</p>

    <el-form
      ref="formRef"
      :model="form"
      :rules="rules"
      @submit.prevent="handleRegister"
    >
      <el-form-item prop="username">
        <el-input v-model="form.username" placeholder="用户名" :prefix-icon="User" size="large" />
      </el-form-item>
      <el-form-item prop="password">
        <el-input v-model="form.password" type="password" placeholder="密码（至少6位）" :prefix-icon="Lock" size="large" show-password />
      </el-form-item>
      <el-form-item prop="confirmPassword">
        <el-input v-model="form.confirmPassword" type="password" placeholder="确认密码" :prefix-icon="Lock" size="large" show-password />
      </el-form-item>

      <!-- 手机号和邮箱同行，更紧凑 -->
      <div class="inline-fields">
        <el-form-item prop="phone" class="inline-field">
          <el-input v-model="form.phone" placeholder="手机号（选填）" :prefix-icon="Phone" size="large" />
        </el-form-item>
        <el-form-item prop="email" class="inline-field">
          <el-input v-model="form.email" placeholder="邮箱（选填）" :prefix-icon="Message" size="large" />
        </el-form-item>
      </div>

      <el-form-item>
        <el-button type="primary" size="large" :loading="loading" native-type="submit" class="submit-btn">
          {{ loading ? '注册中...' : '注册' }}
        </el-button>
      </el-form-item>
    </el-form>

    <div class="toggle-view">
      已有账号？<el-button type="primary" link @click="$emit('switch')">立即登录</el-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import type { FormInstance, FormRules } from 'element-plus'
import { User, Lock, Phone, Message } from '@element-plus/icons-vue'

defineProps<{
  loading: boolean
}>()

const emit = defineEmits<{
  (e: 'submit', data: { username: string; password: string; phone: string; email: string }): void
  (e: 'switch'): void
}>()

const formRef = ref<FormInstance>()
const form = reactive({
  username: '', password: '', confirmPassword: '', phone: '', email: ''
})

const rules: FormRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 50, message: '用户名长度为 3-50 个字符', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码至少 6 位', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请确认密码', trigger: 'blur' },
    {
      validator: (_rule: any, value: string, callback: (err?: Error) => void) => {
        if (value !== form.password) callback(new Error('两次密码不一致'))
        else callback()
      },
      trigger: 'blur'
    }
  ]
}

async function handleRegister() {
  await formRef.value?.validate()
  emit('submit', {
    username: form.username,
    password: form.password,
    phone: form.phone,
    email: form.email
  })
}
</script>

<style scoped>
.register-form-wrapper {
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
  margin: 0 0 24px;
}

.inline-fields {
  display: flex;
  gap: 12px;
}

.inline-field {
  flex: 1;
  min-width: 0;
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
  margin-top: 20px;
  font-size: 14px;
  color: #86909c;
}
</style>
