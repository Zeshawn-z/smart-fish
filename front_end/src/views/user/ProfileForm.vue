<template>
  <div class="profile-form">
    <!-- 只读信息 -->
    <div class="form-section">
      <div class="section-title">基本信息</div>
      <div class="info-grid">
        <div class="info-item">
          <span class="info-label">用户名</span>
          <span class="info-value">{{ username }}</span>
        </div>
        <div class="info-item">
          <span class="info-label">角色</span>
          <span class="info-value">
            <el-tag size="small" :type="roleTagType" effect="plain" round>{{ roleText }}</el-tag>
          </span>
        </div>
        <div class="info-item">
          <span class="info-label">注册时间</span>
          <span class="info-value">{{ registerDate }}</span>
        </div>
      </div>
    </div>

    <!-- 可编辑信息 -->
    <div class="form-section">
      <div class="section-title">联系方式</div>
      <el-form :model="form" label-position="top" class="edit-form">
        <div class="form-grid">
          <el-form-item label="手机号">
            <el-input v-model="form.phone" placeholder="请输入手机号" clearable />
          </el-form-item>
          <el-form-item label="邮箱">
            <el-input v-model="form.email" placeholder="请输入邮箱" clearable />
          </el-form-item>
        </div>
        <div class="form-actions">
          <el-button type="primary" @click="handleSave" :loading="saving" :disabled="!hasChanged">
            保存修改
          </el-button>
          <el-button @click="handleReset" :disabled="!hasChanged">重置</el-button>
        </div>
      </el-form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref, computed } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { ElMessage } from 'element-plus'

const authStore = useAuthStore()

const username = authStore.user?.username || ''
const role = authStore.user?.role || 'user'

const roleText = computed(() => {
  const map: Record<string, string> = { admin: '管理员', staff: '工作人员', user: '普通用户' }
  return map[role] || '普通用户'
})

const roleTagType = computed(() => {
  if (role === 'admin') return 'danger'
  if (role === 'staff') return 'warning'
  return 'info'
})

const registerDate = computed(() => {
  const t = authStore.user?.register_time
  if (!t) return '--'
  return t.substring(0, 10)
})

const originalPhone = authStore.user?.phone || ''
const originalEmail = authStore.user?.email || ''

const form = reactive({
  phone: originalPhone,
  email: originalEmail
})

const saving = ref(false)

const hasChanged = computed(() => {
  return form.phone !== originalPhone || form.email !== originalEmail
})

function handleReset() {
  form.phone = originalPhone
  form.email = originalEmail
}

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

<style scoped>
.profile-form {
  display: flex;
  flex-direction: column;
  gap: 28px;
}

.form-section {
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.section-title {
  font-size: 14px;
  font-weight: 600;
  color: #1d2129;
  padding-bottom: 10px;
  border-bottom: 1px solid #f0f0f0;
}

/* 信息展示 */
.info-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 0;
}

.info-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
  padding: 12px 0;
}

.info-label {
  font-size: 12px;
  color: #a0a4ad;
}

.info-value {
  font-size: 14px;
  color: #1d2129;
  font-weight: 500;
}

/* 表单编辑 */
.edit-form :deep(.el-form-item__label) {
  font-size: 13px;
  color: #86909c;
  padding-bottom: 4px;
}

.form-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 0 20px;
}

.form-actions {
  display: flex;
  gap: 10px;
  padding-top: 4px;
}

@media (max-width: 768px) {
  .info-grid {
    grid-template-columns: 1fr 1fr;
  }

  .form-grid {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 480px) {
  .info-grid {
    grid-template-columns: 1fr;
  }
}
</style>
