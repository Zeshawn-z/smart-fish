<template>
  <div class="resource-manager">
    <!-- 顶部操作栏 -->
    <div class="manager-header">
      <div class="header-left">
        <h3 class="manager-title">{{ config.title }}</h3>
        <el-tag type="info" size="small" effect="plain" round>共 {{ totalCount }} 条</el-tag>
      </div>
      <div class="header-right">
        <el-input
          v-if="config.searchable !== false"
          v-model="searchText"
          :placeholder="`搜索${config.title.replace('管理', '')}...`"
          :prefix-icon="Search"
          clearable
          class="search-input"
          @clear="onSearch"
          @keyup.enter="onSearch"
        />
        <el-button v-if="config.creatable !== false" type="primary" :icon="Plus" @click="openCreate">
          新增
        </el-button>
        <el-button :icon="Refresh" circle @click="loadData" />
      </div>
    </div>

    <!-- 表格 -->
    <div class="table-wrapper">
      <el-table
        :data="tableData"
        v-loading="loading"
        stripe
        border
        class="manager-table"
        :empty-text="`暂无${config.title.replace('管理', '')}数据`"
        :row-class-name="tableRowClassName"
        :header-cell-style="{ background: '#f8f9fb', fontWeight: 600, color: '#303133', fontSize: '13px' }"
      >
        <el-table-column
          v-for="col in config.columns"
          :key="col.prop"
          :prop="col.prop"
          :label="col.label"
          :width="col.width"
          :min-width="col.minWidth"
          :align="col.align || 'left'"
          :show-overflow-tooltip="col.showOverflow !== false"
        >
          <template #default="{ row }" v-if="col.render || col.tag || col.formatter">
            <el-tag v-if="col.tag" :type="col.tag.type ? col.tag.type(row) : 'info'" size="small" effect="light" round>
              {{ col.tag.label ? col.tag.label(row) : row[col.prop] }}
            </el-tag>
            <span v-else-if="col.formatter">{{ col.formatter(row) }}</span>
            <span v-else-if="col.render" v-html="col.render(row)"></span>
          </template>
        </el-table-column>

        <!-- 操作列 -->
        <el-table-column label="操作" :width="config.actionWidth || 150" align="center" fixed="right">
          <template #default="{ row }">
            <el-button size="small" type="primary" link :icon="Edit" @click="openEdit(row)">
              编辑
            </el-button>
            <el-divider direction="vertical" />
            <el-button size="small" type="danger" link :icon="Delete" @click="handleDelete(row)">
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <!-- 分页 -->
    <div class="manager-footer" v-if="totalCount > pageSize">
      <el-pagination
        v-model:current-page="currentPage"
        :page-size="pageSize"
        :total="totalCount"
        layout="total, prev, pager, next, jumper"
        background
        @current-change="handlePageChange"
      />
    </div>

    <!-- 编辑/新增对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="isCreate ? `新增${config.title.replace('管理', '')}` : `编辑${config.title.replace('管理', '')}`"
      width="560px"
      :close-on-click-modal="false"
      destroy-on-close
    >
      <el-form
        v-if="formFields.length > 0"
        ref="formRef"
        :model="formData"
        label-position="top"
        class="edit-form"
      >
        <el-form-item
          v-for="field in formFields"
          :key="field.prop"
          :label="field.label"
          :required="field.required"
          :prop="field.prop"
        >
          <!-- input -->
          <el-input
            v-if="field.type === 'input'"
            v-model="formData[field.prop]"
            :placeholder="field.placeholder || `请输入${field.label}`"
            :readonly="field.readonly"
          />
          <!-- textarea -->
          <el-input
            v-else-if="field.type === 'textarea'"
            v-model="formData[field.prop]"
            type="textarea"
            :rows="3"
            :placeholder="field.placeholder || `请输入${field.label}`"
          />
          <!-- number -->
          <el-input-number
            v-else-if="field.type === 'number'"
            v-model="formData[field.prop]"
            style="width: 100%"
          />
          <!-- select -->
          <el-select
            v-else-if="field.type === 'select'"
            v-model="formData[field.prop]"
            :placeholder="field.placeholder || `请选择${field.label}`"
            style="width: 100%"
          >
            <el-option
              v-for="opt in field.options || []"
              :key="opt.value"
              :label="opt.label"
              :value="opt.value"
            />
          </el-select>
          <!-- switch -->
          <el-switch v-else-if="field.type === 'switch'" v-model="formData[field.prop]" />
          <!-- datetime -->
          <el-date-picker
            v-else-if="field.type === 'datetime'"
            v-model="formData[field.prop]"
            type="datetime"
            :placeholder="field.placeholder || `选择${field.label}`"
            style="width: 100%"
          />
        </el-form-item>
      </el-form>

      <div v-else class="no-form-hint">
        <el-empty description="该资源暂不支持表单编辑" :image-size="60" />
      </div>

      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button
          type="primary"
          :loading="isSaving"
          :disabled="formFields.length === 0"
          @click="handleSave"
        >
          {{ isCreate ? '创建' : '保存' }}
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Search, Edit, Delete, Refresh } from '@element-plus/icons-vue'
import type { ResourceConfig, FormFieldConfig } from '@/types'

const props = defineProps<{ config: ResourceConfig }>()

const emit = defineEmits<{
  (e: 'create'): void
  (e: 'edit', row: any): void
}>()

// ===== 数据状态 =====
const tableData = ref<any[]>([])
const loading = ref(false)
const searchText = ref('')
const currentPage = ref(1)
const pageSize = 12
const totalCount = ref(0)

// ===== 对话框状态 =====
const dialogVisible = ref(false)
const isCreate = ref(false)
const isSaving = ref(false)
const formData = ref<Record<string, any>>({})
const editingRowId = ref<number | null>(null)

const formFields = computed<FormFieldConfig[]>(() => props.config.formFields || [])

// ===== 数据加载（后端分页） =====
async function loadData() {
  loading.value = true
  try {
    const params: Record<string, any> = {}
    if (searchText.value) params.search = searchText.value
    params.page = currentPage.value
    params.page_size = pageSize
    const result = await props.config.loadFn(params)
    tableData.value = result.data
    totalCount.value = result.total
  } catch (err) {
    console.error(`加载${props.config.title}失败:`, err)
    ElMessage.error(`加载${props.config.title}失败`)
  } finally {
    loading.value = false
  }
}

function onSearch() {
  currentPage.value = 1
  loadData()
}

function handlePageChange(page: number) {
  currentPage.value = page
  loadData()
}

function tableRowClassName({ rowIndex }: { rowIndex: number }) {
  return rowIndex % 2 === 0 ? '' : 'stripe-row'
}

// ===== CRUD 操作 =====
function openCreate() {
  if (formFields.value.length === 0 && !props.config.createFn) {
    emit('create')
    return
  }
  isCreate.value = true
  editingRowId.value = null
  formData.value = {}
  // 初始化默认值
  for (const field of formFields.value) {
    formData.value[field.prop] = field.type === 'number' ? 0 : field.type === 'switch' ? false : ''
  }
  dialogVisible.value = true
}

function openEdit(row: any) {
  if (formFields.value.length === 0 && !props.config.updateFn) {
    emit('edit', row)
    return
  }
  isCreate.value = false
  editingRowId.value = row.id ?? row.device_id
  formData.value = {}
  for (const field of formFields.value) {
    formData.value[field.prop] = row[field.prop] ?? ''
  }
  dialogVisible.value = true
}

async function handleSave() {
  // 简单校验必填
  for (const field of formFields.value) {
    if (field.required && !formData.value[field.prop] && formData.value[field.prop] !== 0) {
      ElMessage.warning(`请填写${field.label}`)
      return
    }
  }

  isSaving.value = true
  try {
    if (isCreate.value) {
      if (props.config.createFn) {
        await props.config.createFn(formData.value)
        ElMessage.success('创建成功')
      } else {
        ElMessage.info('创建功能暂未实现')
      }
    } else {
      if (props.config.updateFn && editingRowId.value != null) {
        await props.config.updateFn(editingRowId.value, formData.value)
        ElMessage.success('更新成功')
      } else {
        ElMessage.info('编辑功能暂未实现')
      }
    }
    dialogVisible.value = false
    loadData()
  } catch (err: any) {
    ElMessage.error(err?.message || (isCreate.value ? '创建失败' : '更新失败'))
  } finally {
    isSaving.value = false
  }
}

async function handleDelete(row: any) {
  try {
    await ElMessageBox.confirm('确认删除？此操作不可恢复。', '警告', {
      confirmButtonText: '确认',
      cancelButtonText: '取消',
      type: 'warning'
    })
    const id = row.id ?? row.device_id
    await props.config.deleteFn(id)
    ElMessage.success('删除成功')
    loadData()
  } catch (err: any) {
    if (err !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

// ===== 暴露刷新方法给父组件 =====
defineExpose({ loadData })

// ===== 监听 config 变化 =====
watch(
  () => props.config.resource,
  () => {
    searchText.value = ''
    currentPage.value = 1
    loadData()
  }
)

onMounted(() => loadData())
</script>

<style scoped>
.resource-manager {
  height: 100%;
  display: flex;
  flex-direction: column;
  background: #fff;
}

/* 头部操作栏 */
.manager-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 20px 24px 16px;
  flex-shrink: 0;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.manager-title {
  font-size: 18px;
  font-weight: 700;
  color: #1d2129;
  margin: 0;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 10px;
}

.search-input {
  width: 240px;
}

/* 表格区域 */
.table-wrapper {
  flex: 1;
  overflow: hidden;
  padding: 0 24px;
}

.manager-table {
  width: 100%;
}

.manager-table :deep(.el-table__header-wrapper th) {
  font-size: 13px;
}

.manager-table :deep(.el-table__body-wrapper) {
  font-size: 13.5px;
}

.manager-table :deep(.el-table__row) {
  transition: background 0.15s;
}

.manager-table :deep(.el-table__row:hover > td) {
  background: #f0f5ff !important;
}

/* 分页 */
.manager-footer {
  padding: 16px 24px;
  display: flex;
  justify-content: flex-end;
  flex-shrink: 0;
  border-top: 1px solid #f0f0f0;
}

/* 编辑表单 */
.edit-form {
  max-height: 60vh;
  overflow-y: auto;
}

.edit-form :deep(.el-form-item__label) {
  font-weight: 600;
  color: #303133;
}

.no-form-hint {
  padding: 24px 0;
}
</style>
