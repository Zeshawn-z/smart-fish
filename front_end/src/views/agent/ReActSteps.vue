<template>
  <div class="react-steps">
    <div
      v-for="step in steps"
      :key="step.id"
      class="react-step"
      :class="[`react-step--${step.type}`]"
    >
      <!-- 折叠行 -->
      <div class="react-step__row" @click="toggleCollapse(step.id)">
        <!-- Chevron 箭头 SVG（右=折叠，下=展开，旋转过渡） -->
        <svg
          class="react-step__chevron"
          :class="{ 'is-expanded': !isCollapsed(step.id) }"
          viewBox="0 0 16 16"
          width="14"
          height="14"
          fill="none"
        >
          <path d="M6 4l4 4-4 4" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"/>
        </svg>

        <!-- 步骤标签 -->
        <span class="react-step__tag" :class="`react-step__tag--${step.type}`">
          {{ stepTag(step.type) }}
        </span>

        <!-- 工具名（仅 tool_call） -->
        <code v-if="step.type === 'tool_call'" class="react-step__tool-name">
          {{ step.toolCall?.name || 'unknown' }}
        </code>

        <!-- 状态指示 -->
        <span v-if="step.type === 'tool_call' && step.toolCall" class="react-step__status">
          <span v-if="step.toolCall.status === 'calling'" class="status-calling">
            <span class="dot-dot"></span><span class="dot-dot"></span><span class="dot-dot"></span>
          </span>
          <span v-else-if="step.toolCall.status === 'success'" class="status-success">
            <el-icon :size="12"><CircleCheck /></el-icon>
          </span>
          <span v-else-if="step.toolCall.status === 'error'" class="status-error">✗</span>
        </span>

        <!-- 流式呼吸圆点 -->
        <span
          v-if="step.streaming && !(step.type === 'tool_call' && step.toolCall?.status === 'calling')"
          class="breathing-dots"
        >
          <span></span><span></span><span></span>
        </span>
      </div>

      <!-- 折叠内容（CSS grid row 过渡实现丝滑动画） -->
      <div class="react-step__collapse" :data-open="!isCollapsed(step.id)">
        <div class="react-step__collapse-inner">
          <!-- 工具参数 -->
          <div v-if="step.toolCall?.args && Object.keys(step.toolCall.args).length > 0" class="react-step__params">
            <code v-for="(val, key) in step.toolCall.args" :key="key" class="react-step__param-tag">
              {{ key }}: {{ JSON.stringify(val) }}
            </code>
          </div>

          <!-- 内容文本 -->
          <div class="react-step__content">
            <span v-html="renderText(step.content)"></span>
          </div>

          <!-- 加载动画 -->
          <div v-if="step.type === 'tool_call' && step.toolCall?.status === 'calling' && !step.streaming" class="react-step__loading">
            <span class="breathing-dots">
              <span></span><span></span><span></span>
            </span>
            <span>等待工具返回</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { CircleCheck } from '@element-plus/icons-vue'
import type { ReActStep } from '@/types/agent-chat'

const props = defineProps<{
  steps: ReActStep[]
  finalStarted?: boolean
}>()

const collapsedIds = ref<Set<string>>(new Set())
const knownStepIds = ref<Set<string>>(new Set())
const manuallyExpandedIds = ref<Set<string>>(new Set())
const previousStepIds = ref<string[]>([])

watch(
  () => props.steps.map(step => `${step.id}:${step.streaming ? '1' : '0'}`).join('|'),
  () => {
    const currentStepIds = props.steps.map(step => step.id)
    const currentIdSet = new Set(currentStepIds)

    if (currentStepIds.length === 0) {
      collapsedIds.value = new Set()
      knownStepIds.value = new Set()
      manuallyExpandedIds.value = new Set()
      previousStepIds.value = []
      return
    }

    const prevStepIds = previousStepIds.value
    const newlyAddedIds = currentStepIds.filter(id => !prevStepIds.includes(id))

    // 清理已不存在的步骤状态，避免跨会话残留。
    collapsedIds.value = new Set([...collapsedIds.value].filter(id => currentIdSet.has(id)))
    knownStepIds.value = new Set([...knownStepIds.value].filter(id => currentIdSet.has(id)))
    manuallyExpandedIds.value = new Set([...manuallyExpandedIds.value].filter(id => currentIdSet.has(id)))

    for (const stepId of newlyAddedIds) {
      const stepIndex = currentStepIds.indexOf(stepId)
      const step = props.steps[stepIndex]

      // 新步骤初始默认为折叠，再自动展开当前流式步骤。
      knownStepIds.value.add(stepId)
      collapsedIds.value.add(stepId)
      if (step?.streaming) {
        collapsedIds.value.delete(stepId)
      }

      // 在“下一个步骤出现时”折叠前一个步骤。
      const previousStepId = stepIndex > 0 ? currentStepIds[stepIndex - 1] : null
      if (previousStepId && !manuallyExpandedIds.value.has(previousStepId)) {
        collapsedIds.value.add(previousStepId)
      }
    }

    for (const stepId of currentStepIds) {
      if (!knownStepIds.value.has(stepId)) {
        knownStepIds.value.add(stepId)
        collapsedIds.value.add(stepId)
      }
    }

    previousStepIds.value = currentStepIds
    // 重新赋值触发依赖更新。
    collapsedIds.value = new Set(collapsedIds.value)
  },
  { immediate: true },
)

watch(
  () => !!props.finalStarted,
  (started, prevStarted) => {
    if (!started || prevStarted) return

    const lastStepId = props.steps[props.steps.length - 1]?.id
    if (!lastStepId) return

    collapsedIds.value.add(lastStepId)
    collapsedIds.value = new Set(collapsedIds.value)
  },
  { immediate: true },
)

function toggleCollapse(id: string) {
  if (collapsedIds.value.has(id)) {
    collapsedIds.value.delete(id)
    manuallyExpandedIds.value.add(id)
  } else {
    collapsedIds.value.add(id)
    manuallyExpandedIds.value.delete(id)
  }
  collapsedIds.value = new Set(collapsedIds.value)
}

function isCollapsed(id: string): boolean {
  return collapsedIds.value.has(id)
}

function stepTag(type: string): string {
  const map: Record<string, string> = {
    thought: '思考',
    tool_call: '调用工具',
    tool_result: '观察结果',
    action: '执行',
    final_answer: '回答',
  }
  return map[type] || type
}

function renderText(text: string): string {
  return text
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/\*\*(.*?)\*\*/g, '<strong>$1</strong>')
    .replace(/`([^`]+)`/g, '<code>$1</code>')
    .replace(/\n/g, '<br>')
}
</script>

<style scoped>
.react-steps {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

/* ===== 步骤行 ===== */
.react-step__row {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 8px;
  border-radius: 6px;
  cursor: pointer;
  user-select: none;
  transition: background 0.15s;
}

.react-step__row:hover {
  background: rgba(0, 0, 0, 0.02);
}

/* ===== Chevron 箭头 ===== */
.react-step__chevron {
  color: #b0b3b8;
  transition: transform 0.25s cubic-bezier(0.4, 0, 0.2, 1);
  flex-shrink: 0;
}

.react-step__chevron.is-expanded {
  transform: rotate(90deg);
}

/* ===== 步骤标签 ===== */
.react-step__tag {
  font-size: 12px;
  font-weight: 500;
  border-radius: 4px;
  padding: 1px 7px;
  flex-shrink: 0;
}

.react-step__tag--thought {
  background: rgba(124, 77, 255, 0.08);
  color: #7c4dff;
}

.react-step__tag--tool_call {
  background: rgba(255, 109, 0, 0.08);
  color: #e65100;
}

.react-step__tag--tool_result {
  background: rgba(0, 188, 212, 0.08);
  color: #00838f;
}

.react-step__tag--action {
  background: rgba(76, 175, 80, 0.08);
  color: #2e7d32;
}

.react-step__tag--final_answer {
  background: rgba(33, 150, 243, 0.08);
  color: #1565c0;
}

/* ===== 工具名 ===== */
.react-step__tool-name {
  font-family: ui-monospace, Consolas, monospace;
  font-size: 11px;
  color: #909399;
}

/* ===== 状态指示 ===== */
.react-step__status {
  display: flex;
  align-items: center;
  margin-left: 2px;
}

.status-calling {
  display: flex;
  gap: 2px;
}

.status-calling .dot-dot {
  width: 4px;
  height: 4px;
  border-radius: 50%;
  background: #ff9800;
  animation: dot-bounce 1.2s ease-in-out infinite;
}

.status-calling .dot-dot:nth-child(2) { animation-delay: 0.15s; }
.status-calling .dot-dot:nth-child(3) { animation-delay: 0.3s; }

.status-success {
  display: inline-flex;
  align-items: center;
  color: #43a047;
}

.status-error {
  color: #e53935;
  font-size: 12px;
  font-weight: 600;
}

/* ===== 呼吸圆点 ===== */
.breathing-dots {
  display: inline-flex;
  align-items: center;
  gap: 3px;
  margin-left: 4px;
}

.breathing-dots span {
  width: 4px;
  height: 4px;
  border-radius: 50%;
  background: #409eff;
  animation: breathe 1.4s ease-in-out infinite;
}

.breathing-dots span:nth-child(2) { animation-delay: 0.2s; }
.breathing-dots span:nth-child(3) { animation-delay: 0.4s; }

/* ===== 折叠区域（CSS Grid 过渡 —— 丝滑动画） ===== */
.react-step__collapse {
  display: grid;
  grid-template-rows: 0fr;
  transition: grid-template-rows 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.react-step__collapse[data-open="true"] {
  grid-template-rows: 1fr;
}

.react-step__collapse-inner {
  overflow: hidden;
  min-height: 0;
}

/* ===== 折叠内容样式 ===== */
.react-step__params {
  display: flex;
  flex-wrap: wrap;
  gap: 5px;
  padding: 6px 0 4px 32px;
}

.react-step__param-tag {
  font-family: ui-monospace, Consolas, monospace;
  font-size: 11px;
  padding: 2px 7px;
  border-radius: 4px;
  background: #f5f5f7;
  color: #909399;
}

.react-step__content {
  padding: 4px 0 8px 32px;
  font-size: 13px;
  line-height: 1.75;
  color: #a8a8b0;
  white-space: pre-wrap;
  word-break: break-word;
}

.react-step__content :deep(strong) {
  color: #b8b8c0;
  font-weight: 600;
}

.react-step__content :deep(code) {
  font-family: ui-monospace, Consolas, monospace;
  font-size: 11px;
  padding: 1px 5px;
  background: #f5f5f7;
  border-radius: 3px;
  color: #909399;
}

/* ===== 加载提示 ===== */
.react-step__loading {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 4px 0 8px 32px;
  font-size: 12px;
  color: #b8b8c0;
}

/* ===== Animations ===== */
@keyframes dot-bounce {
  0%, 80%, 100% { transform: scale(0.4); opacity: 0.4; }
  40% { transform: scale(1); opacity: 1; }
}

@keyframes breathe {
  0%, 100% { opacity: 0.25; transform: scale(0.8); }
  50% { opacity: 1; transform: scale(1); }
}

/* ===== Dark mode ===== */
@media (prefers-color-scheme: dark) {
  .react-step__row:hover {
    background: rgba(255, 255, 255, 0.03);
  }

  .react-step__chevron {
    color: #5a5b65;
  }

  .react-step__tag--thought {
    background: rgba(124, 77, 255, 0.12);
    color: #b388ff;
  }

  .react-step__tag--tool_call {
    background: rgba(255, 109, 0, 0.12);
    color: #ffab40;
  }

  .react-step__tag--tool_result {
    background: rgba(0, 188, 212, 0.12);
    color: #4dd0e1;
  }

  .react-step__tool-name {
    color: #70727a;
  }

  .react-step__param-tag {
    background: #2a2b34;
    color: #70727a;
  }

  .react-step__content {
    color: #70727a;
  }

  .react-step__content :deep(strong) {
    color: #8a8b95;
  }

  .react-step__content :deep(code) {
    background: #2a2b34;
    color: #70727a;
  }

  .react-step__loading {
    color: #70727a;
  }

  .breathing-dots span {
    background: #409eff;
  }
}
</style>
