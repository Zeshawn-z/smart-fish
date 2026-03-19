<template>
  <el-form ref="formRef" :model="form" :rules="rules" label-width="100px" style="max-width: 500px;">
    <el-form-item label="当前密码" prop="oldPassword">
      <el-input v-model="form.oldPassword" type="password" show-password />
    </el-form-item>
    <el-form-item label="新密码" prop="newPassword">
      <el-input v-model="form.newPassword" type="password" show-password />
    </el-form-item>
    <el-form-item label="确认新密码" prop="confirmPassword">
      <el-input v-model="form.confirmPassword" type="password" show-password />
    </el-form-item>
    <el-form-item>
      <el-button type="primary" @click="handleSubmit" :loading="saving">修改密码</el-button>
    </el-form-item>
  </el-form>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'

const authStore = useAuthStore()
const formRef = ref<FormInstance>()
const saving = ref(false)

const form = reactive({
  oldPassword: '',
  newPassword: '',
  confirmPassword: ''
})

const rules: FormRules = {
  oldPassword: [{ required: true, message: '请输入当前密码', trigger: 'blur' }],
  newPassword: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, message: '密码至少 6 位', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请确认新密码', trigger: 'blur' },
    {
      validator: (_rule: any, value: string, callback: (err?: Error) => void) => {
        if (value !== form.newPassword) callback(new Error('两次密码不一致'))
        else callback()
      },
      trigger: 'blur'
    }
  ]
}

async function handleSubmit() {
  await formRef.value?.validate()
  saving.value = true
  try {
    await authStore.updatePassword(form.oldPassword, form.newPassword)
    ElMessage.success('密码已修改')
    form.oldPassword = ''
    form.newPassword = ''
    form.confirmPassword = ''
  } catch (e: any) {
    ElMessage.error(e.response?.data?.error || '修改失败')
  } finally {
    saving.value = false
  }
}
</script>
