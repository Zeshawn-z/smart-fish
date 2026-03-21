<template>
  <el-dialog v-model="visible" title="编辑帖子" width="640px" :close-on-click-modal="false" class="edit-post-dialog">
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
    </el-form>

    <template #footer>
      <el-button @click="visible = false">取消</el-button>
      <el-button type="primary" :loading="isSaving" round @click="handleSave">
        {{ isSaving ? '保存中...' : '保存修改' }}
      </el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { PostService } from '@/services/CommunityService'
import type { Post } from '@/types'

const TAGS = ['钓鱼日记', '经验分享', '装备测评', '钓点推荐', '问答求助']
const TAG_ICONS: Record<string, string> = {
  '钓鱼日记': '📖',
  '经验分享': '💡',
  '装备测评': '🔧',
  '钓点推荐': '📍',
  '问答求助': '❓'
}

const props = defineProps<{
  post: Post | null
}>()

const visible = defineModel<boolean>({ default: false })
const emit = defineEmits<{ (e: 'updated'): void }>()

const isSaving = ref(false)
const form = ref({ title: '', body: '', tag: '' })

watch(visible, (v) => {
  if (v && props.post) {
    form.value = {
      title: props.post.title || '',
      body: props.post.body || '',
      tag: props.post.tag || ''
    }
  }
})

async function handleSave() {
  if (!form.value.title.trim() || !form.value.body.trim()) {
    ElMessage.warning('请填写标题和内容')
    return
  }
  if (!props.post) return

  isSaving.value = true
  try {
    await PostService.update(props.post.id, {
      title: form.value.title,
      body: form.value.body,
      tag: form.value.tag || undefined
    } as Partial<Post>)

    ElMessage.success('修改成功！')
    visible.value = false
    emit('updated')
  } catch {
    ElMessage.error('修改失败，请稍后再试')
  } finally {
    isSaving.value = false
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
</style>
