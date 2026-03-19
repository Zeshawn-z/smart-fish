<template>
  <el-form :model="form" label-width="80px" style="max-width: 500px;">
    <el-form-item label="用户名">
      <el-input :model-value="username" disabled />
    </el-form-item>
    <el-form-item label="角色">
      <el-tag>{{ role }}</el-tag>
    </el-form-item>
    <el-form-item label="手机号">
      <el-input v-model="form.phone" placeholder="请输入手机号" />
    </el-form-item>
    <el-form-item label="邮箱">
      <el-input v-model="form.email" placeholder="请输入邮箱" />
    </el-form-item>
    <el-form-item>
      <el-button type="primary" @click="handleSave" :loading="saving">保存修改</el-button>
    </el-form-item>
  </el-form>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { ElMessage } from 'element-plus'

const authStore = useAuthStore()

const username = authStore.user?.username || ''
const role = authStore.user?.role || ''
const form = reactive({
  phone: authStore.user?.phone || '',
  email: authStore.user?.email || ''
})
const saving = ref(false)

async function handleSave() {
  saving.value = true
  try {
    await authStore.updateProfile({ phone: form.phone, email: form.email })
    ElMessage.success('个人信息已更新')
  } catch {
    ElMessage.error('更新失败')
  } finally {
    saving.value = false
  }
}
</script>
