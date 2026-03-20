<template>
  <el-dialog v-model="visible" title="发布新帖" width="640px" :close-on-click-modal="false" class="create-post-dialog">
    <el-form :model="form" label-position="top">
      <el-form-item label="标题" required>
        <el-input
          v-model="form.title"
          placeholder="给帖子起个吸引人的标题吧"
          maxlength="100"
          show-word-limit
        />
      </el-form-item>

      <el-form-item label="标签">
        <div class="tag-select">
          <button
            v-for="t in TAGS"
            :key="t"
            type="button"
            class="tag-option"
            :class="{ selected: form.tag === t }"
            @click="form.tag = form.tag === t ? '' : t"
          >
            {{ TAG_ICONS[t] }} {{ t }}
          </button>
        </div>
      </el-form-item>

      <el-form-item label="内容" required>
        <el-input
          v-model="form.body"
          type="textarea"
          :rows="6"
          placeholder="分享你的垂钓故事..."
          maxlength="5000"
          show-word-limit
          resize="vertical"
        />
      </el-form-item>

      <el-form-item label="配图（可选）">
        <el-upload
          :auto-upload="false"
          list-type="picture-card"
          :file-list="fileList"
          :on-change="handleFileChange"
          :on-remove="handleFileRemove"
          :limit="9"
          accept="image/*"
          :on-exceed="handleExceed"
          class="image-uploader"
        >
          <div class="upload-trigger">
            <el-icon :size="24"><Plus /></el-icon>
            <span>添加图片</span>
          </div>
        </el-upload>
        <div class="upload-tip">支持 jpg/png/gif，最多 9 张</div>
      </el-form-item>
    </el-form>

    <template #footer>
      <el-button @click="visible = false">取消</el-button>
      <el-button type="primary" :loading="isCreating" round @click="handleCreate">
        {{ isCreating ? '发布中...' : '发布帖子' }}
      </el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { Plus } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import type { UploadFile, UploadFiles } from 'element-plus'
import { CommunityService } from '@/services/CommunityService'

const TAGS = ['钓鱼日记', '经验分享', '装备测评', '钓点推荐', '问答求助']
const TAG_ICONS: Record<string, string> = {
  '钓鱼日记': '📖',
  '经验分享': '💡',
  '装备测评': '🔧',
  '钓点推荐': '📍',
  '问答求助': '❓'
}

const visible = defineModel<boolean>({ default: false })
const emit = defineEmits<{ (e: 'created'): void }>()

const isCreating = ref(false)
const form = ref({ title: '', body: '', tag: '' })
const fileList = ref<UploadFile[]>([])

watch(visible, (v) => {
  if (v) {
    form.value = { title: '', body: '', tag: '' }
    fileList.value = []
  }
})

function handleFileChange(_file: UploadFile, newFileList: UploadFiles) {
  fileList.value = newFileList as UploadFile[]
}

function handleFileRemove(_file: UploadFile, newFileList: UploadFiles) {
  fileList.value = newFileList as UploadFile[]
}

function handleExceed() {
  ElMessage.warning('最多上传 9 张图片')
}

async function handleCreate() {
  if (!form.value.title.trim() || !form.value.body.trim()) {
    ElMessage.warning('请填写标题和内容')
    return
  }
  isCreating.value = true
  try {
    // 1. 创建帖子
    const post = await CommunityService.createPost({
      title: form.value.title,
      body: form.value.body,
      tag: form.value.tag || undefined
    })

    // 2. 上传图片（如有）
    const postId = post.id
    if (fileList.value.length > 0 && postId) {
      const uploadPromises = fileList.value
        .filter(f => f.raw)
        .map(f => CommunityService.uploadPostImage(postId, f.raw!))
      await Promise.allSettled(uploadPromises)
    }

    ElMessage.success('发布成功！')
    visible.value = false
    emit('created')
  } catch {
    ElMessage.error('发布失败，请稍后再试')
  } finally {
    isCreating.value = false
  }
}
</script>

<style scoped>
/* 标签选择 */
.tag-select {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.tag-option {
  padding: 5px 14px;
  border-radius: 16px;
  border: 1px solid #dcdfe6;
  background: #fafafa;
  font-size: 13px;
  color: #606266;
  cursor: pointer;
  transition: all 0.2s;
}

.tag-option:hover {
  border-color: #409eff;
  color: #409eff;
}

.tag-option.selected {
  background: #ecf5ff;
  border-color: #409eff;
  color: #409eff;
  font-weight: 600;
}

/* 上传区域 */
.image-uploader :deep(.el-upload--picture-card) {
  width: 100px;
  height: 100px;
  border-radius: 10px;
  border: 1px dashed #d9d9d9;
  transition: border-color 0.2s;
}

.image-uploader :deep(.el-upload--picture-card:hover) {
  border-color: #409eff;
}

.image-uploader :deep(.el-upload-list__item) {
  width: 100px;
  height: 100px;
  border-radius: 10px;
}

.upload-trigger {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 4px;
  color: #909399;
  font-size: 12px;
  height: 100%;
}

.upload-tip {
  font-size: 12px;
  color: #a0a4ad;
  margin-top: 6px;
}
</style>
