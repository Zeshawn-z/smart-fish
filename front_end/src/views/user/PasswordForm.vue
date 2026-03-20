<template>
  <div class="password-form">
    <div class="form-section">
      <div class="section-title">修改登录密码</div>
      <p class="section-desc">修改后需要使用新密码重新登录</p>

      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-position="top"
        class="pwd-form"
      >
        <el-form-item label="当前密码" prop="oldPassword">
          <el-input v-model="form.oldPassword" type="password" show-password placeholder="请输入当前密码" />
        </el-form-item>
        <el-form-item label="新密码" prop="newPassword">
          <el-input v-model="form.newPassword" type="password" show-password placeholder="至少 6 位" />
        </el-form-item>
        <el-form-item label="确认新密码" prop="confirmPassword">
          <el-input v-model="form.confirmPassword" type="password" show-password placeholder="再次输入新密码" />
        </el-form-item>
        <div class="form-actions">
          <el-button type="primary" @click="handleSubmit" :loading="saving">修改密码</el-button>
        </div>
      </el-form>
    </div>
  </div>
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

<style scoped>
.password-form {
  display: flex;
  flex-direction: column;
}

.form-section {
  display: flex;
  flex-direction: column;
}

.section-title {
  font-size: 14px;
  font-weight: 600;
  color: #1d2129;
  padding-bottom: 6px;
}

.section-desc {
  font-size: 12.5px;
  color: #a0a4ad;
  margin: 0 0 18px;
}

.pwd-form {
  max-width: 380px;
}

.pwd-form :deep(.el-form-item__label) {
  font-size: 13px;
  color: #86909c;
  padding-bottom: 4px;
}

.form-actions {
  padding-top: 4px;
}
</style>
