import { ref, computed, nextTick, type Ref } from 'vue'
import type { AgentMessage, ChatSession, StreamEvent, ReActStep } from '@/types/agent-chat'
import { createMockStreamReply } from './useMockAgentStream'

/**
 * 智能体对话 composable
 *
 * 管理会话列表、消息收发、流式事件处理。
 * 视图层只需消费返回的 ref 即可。
 */
export function useAgentChat() {
  // ===== 会话 =====
  const sessions = ref<ChatSession[]>([])
  const currentSessionId = ref<string | null>(null)
  const sidebarCollapsed = ref(false)

  const currentSession = computed(() =>
    sessions.value.find(s => s.id === currentSessionId.value) ?? null,
  )
  const currentMessages = computed<AgentMessage[]>(() =>
    currentSession.value?.messages ?? [],
  )

  function createSession(title = '新对话'): ChatSession {
    const session: ChatSession = {
      id: `session_${Date.now()}`,
      title,
      messages: [],
      createdAt: new Date(),
    }
    sessions.value.unshift(session)
    currentSessionId.value = session.id
    return session
  }

  function startNewChat() {
    createSession()
  }

  function switchSession(id: string) {
    currentSessionId.value = id
    nextTick(scrollToBottom)
  }

  function deleteSession(id: string) {
    const idx = sessions.value.findIndex(s => s.id === id)
    if (idx === -1) return
    sessions.value.splice(idx, 1)
    if (currentSessionId.value === id) {
      currentSessionId.value = sessions.value[0]?.id ?? null
    }
  }

  // ===== 消息与流式 =====
  const inputText = ref('')
  const isGenerating = ref(false)
  const scrollRef = ref<HTMLElement>()
  const inputRef: Ref<any> = ref(null)
  let cancelStream: (() => void) | null = null

  /** 追踪每条助手消息是否已收到 final_start */
  const finalStartedMap = ref<Record<string, boolean>>({})

  /** 只有收到 final_start 后才展示最终回答 */
  function showFinalAnswer(msg: AgentMessage): boolean {
    return !!finalStartedMap.value[msg.id]
  }

  function handleSend() {
    const text = inputText.value.trim()
    if (!text || isGenerating.value) return
    inputText.value = ''
    sendMessage(text)
  }

  function sendQuickPrompt(text: string) {
    if (isGenerating.value) return
    sendMessage(text)
  }

  function sendMessage(text: string) {
    // 确保活跃会话
    if (!currentSession.value) {
      const session = createSession(text.slice(0, 20))
      session.title = text.length > 20 ? text.slice(0, 20) + '...' : text
    } else if (currentSession.value.messages.length === 0) {
      currentSession.value.title = text.length > 20 ? text.slice(0, 20) + '...' : text
    }

    const session = currentSession.value!

    // 用户消息
    const userMsg: AgentMessage = {
      id: `msg_${Date.now()}_user`,
      role: 'user',
      content: text,
      steps: [],
      createdAt: new Date(),
    }
    session.messages.push(userMsg)

    // 空助手消息
    const assistantMsg: AgentMessage = {
      id: `msg_${Date.now()}_assistant`,
      role: 'assistant',
      content: '',
      steps: [],
      streaming: true,
      createdAt: new Date(),
    }
    session.messages.push(assistantMsg)

    finalStartedMap.value[assistantMsg.id] = false

    isGenerating.value = true
    scrollToBottom()

    const assistantMsgId = assistantMsg.id
    const { cancel } = createMockStreamReply(text, (event: StreamEvent) => {
      handleStreamEvent(event, assistantMsgId)
    })
    cancelStream = cancel
  }

  function findMessageById(messageId: string): AgentMessage | null {
    for (const session of sessions.value) {
      const msg = session.messages.find(m => m.id === messageId)
      if (msg) return msg
    }
    return null
  }

  function updateStep(
    msg: AgentMessage,
    stepId: string | undefined,
    updater: (step: ReActStep) => ReActStep,
  ) {
    if (!stepId) return
    const idx = msg.steps.findIndex(s => s.id === stepId)
    if (idx === -1) return
    msg.steps.splice(idx, 1, updater(msg.steps[idx]))
  }

  function handleStreamEvent(event: StreamEvent, messageId: string) {
    const msg = findMessageById(messageId)
    if (!msg) return

    switch (event.type) {
      case 'step_start': {
        const step: ReActStep = {
          id: event.stepId!,
          type: event.stepType ?? 'thought',
          content: event.content ?? '',
          streaming: true,
        }
        msg.steps.push(step)
        break
      }

      case 'step_delta': {
        updateStep(msg, event.stepId, step => ({
          ...step,
          content: step.content + (event.delta ?? ''),
        }))
        break
      }

      case 'step_end': {
        updateStep(msg, event.stepId, step => ({
          ...step,
          streaming: false,
        }))
        break
      }

      case 'tool_start': {
        if (event.toolCall) {
          updateStep(msg, event.stepId, step => ({
            ...step,
            toolCall: { ...event.toolCall!, status: 'calling' },
          }))
        }
        break
      }

      case 'tool_result': {
        if (event.toolCall) {
          updateStep(msg, event.stepId, step => ({
            ...step,
            toolCall: { ...event.toolCall!, status: event.toolCall!.status },
          }))
        }
        break
      }

      case 'final_start': {
        finalStartedMap.value[msg.id] = true
        break
      }

      case 'final_delta': {
        msg.content += event.delta ?? ''
        break
      }

      case 'final_end':
      case 'done': {
        msg.streaming = false
        isGenerating.value = false
        cancelStream = null
        break
      }
    }
  }

  function scrollToBottom() {
    nextTick(() => {
      const el = scrollRef.value
      if (el) el.scrollTop = el.scrollHeight
    })
  }

  // 初始化默认会话
  function init() {
    if (sessions.value.length === 0) {
      createSession()
    }
  }

  return {
    // 会话
    sessions,
    currentSessionId,
    currentSession,
    currentMessages,
    sidebarCollapsed,
    // 输入
    inputText,
    inputRef,
    scrollRef,
    isGenerating,
    // 方法
    showFinalAnswer,
    handleSend,
    sendQuickPrompt,
    startNewChat,
    switchSession,
    deleteSession,
    scrollToBottom,
    init,
  }
}
