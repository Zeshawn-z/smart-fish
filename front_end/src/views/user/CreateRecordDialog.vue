<template>
  <el-dialog v-model="visible" title="新增垂钓记录" width="500px">
    <el-form :model="form" label-position="top">
      <el-form-item label="开始时间" required>
        <el-date-picker v-model="form.start_time" type="datetime" placeholder="选择开始时间" style="width: 100%" />
      </el-form-item>
      <el-form-item label="结束时间" required>
        <el-date-picker v-model="form.end_time" type="datetime" placeholder="选择结束时间" style="width: 100%" />
      </el-form-item>
      <el-form-item label="纬度">
        <el-input-number v-model="form.latitude" :precision="6" :step="0.01" style="width: 100%" />
      </el-form-item>
      <el-form-item label="经度">
        <el-input-number v-model="form.longitude" :precision="6" :step="0.01" style="width: 100%" />
      </el-form-item>
      <el-form-item label="关联设备 ID（可选）">
        <el-input v-model="form.device_id" placeholder="设备编号" />
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="visible = false">取消</el-button>
      <el-button type="primary" :loading="isCreating" @click="handleCreate">确定</el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { FishingRecordService } from '@/services/FishingRecordService'

const visible = defineModel<boolean>({ default: false })
const emit = defineEmits<{ (e: 'created'): void }>()

const isCreating = ref(false)
const form = ref({
  start_time: null as Date | null,
  end_time: null as Date | null,
  latitude: 0,
  longitude: 0,
  device_id: ''
})

watch(visible, (v) => {
  if (v) form.value = { start_time: null, end_time: null, latitude: 0, longitude: 0, device_id: '' }
})

async function handleCreate() {
  if (!form.value.start_time || !form.value.end_time) {
    ElMessage.warning('请填写开始和结束时间')
    return
  }
  isCreating.value = true
  try {
    await FishingRecordService.createRecord({
      start_time: (form.value.start_time as Date).toISOString(),
      end_time: (form.value.end_time as Date).toISOString(),
      latitude: form.value.latitude,
      longitude: form.value.longitude,
      device_id: form.value.device_id || undefined
    })
    ElMessage.success('记录创建成功')
    visible.value = false
    emit('created')
  } catch {
    ElMessage.error('创建失败')
  } finally {
    isCreating.value = false
  }
}
</script>
